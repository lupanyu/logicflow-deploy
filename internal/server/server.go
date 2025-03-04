package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"log"
	"logicflow-deploy/internal/nodes"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Server struct {
	agents map[string]*protocol.AgentConnection // 连接的Agent
	//taskQueue  chan DeploymentTask
	httpServer *gin.Engine
	agentsLock sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		agents: make(map[string]*protocol.AgentConnection),
	}
}
func (s *Server) SetHttp(g *gin.Engine) {
	s.httpServer = g
}
func (s *Server) Start(ip string, port int) {
	s.httpServer.Run(fmt.Sprintf("http://%s:%d", ip, port))
}

// 查看agent的状态 是否可以接受任务
func (s *Server) HandleAgentStatus(agentID string) bool {
	agent, ok := s.agents[agentID]
	if !ok {
		return false
	}
	// 检查Agent状态
	if agent.Status == protocol.AgentReady {
		return true
	}
	return false
}

// 根据节点ID取得connection
func (s *Server) GetAgentConnection(agentID string) (*protocol.AgentConnection, bool) {
	agent, ok := s.agents[agentID]
	return agent, ok
}

func HandleAgentConnection(s *Server, conn *websocket.Conn) {
	// 处理注册消息
	var registerMsg protocol.Message
	if err := conn.ReadJSON(&registerMsg); err != nil {
		log.Printf("读取注册消息失败: %v", err)
		return
	}
	var agentID string
	// 验证消息格式
	if registerMsg.Type != protocol.MsgRegister {
		conn.WriteJSON(protocol.Message{
			Type:      protocol.MsgRegisterResponse,
			AgentID:   registerMsg.AgentID,
			Timestamp: time.Now().UnixNano(),
			Payload: protocol.MessageAuthResponse{Code: 401,
				Message: "Invalid message type"},
		})
		log.Printf("收到无效消息: %+v", registerMsg)
		return
	} else {
		if registerMsg.AgentID == "" {
			conn.WriteJSON(protocol.Message{
				Type:      protocol.MsgRegisterResponse,
				AgentID:   registerMsg.AgentID,
				Timestamp: time.Now().UnixNano(),
				Payload: protocol.MessageAuthResponse{Code: 401,
					Message: "Invalid agent id"},
			})
			log.Printf("收到无效消息: %+v", registerMsg)
			if err := conn.Close(); err != nil {
				log.Printf("关闭连接失败: %v", err)
			}
			return
		}
		agentID = registerMsg.AgentID
		conn.WriteJSON(protocol.Message{
			Type:      protocol.MsgRegisterResponse,
			AgentID:   registerMsg.AgentID,
			Timestamp: time.Now().UnixNano(),
			Payload: protocol.MessageAuthResponse{Code: 200,
				Message: "注册成功"},
		})
	}

	//把agent加入到server的map中
	s.agentsLock.Lock()
	s.agents[agentID] = &protocol.AgentConnection{
		Conn:       conn,
		LastActive: time.Now(),
		Status:     protocol.AgentReady,
	}
	s.agentsLock.Unlock()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("关闭连接失败: %v", err)
		}
	}()
	log.Println("等待消息...")
	for {
		var msg protocol.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("读取超时")
			} else {
				log.Printf("读取消息失败: %v", err)
			}
			log.Printf("读取消息失败: %v", err)
			// 把agent从map中移除
			s.agentsLock.Lock()
			delete(s.agents, agentID)
			s.agentsLock.Unlock()
			break
		}
		log.Println("收到的消息是:", msg)
		switch msg.Type {

		case protocol.MsgStatus: // 处理最终的任务状态
			handleStatusUpdate(s, msg)
		case protocol.MsgTaskStep: // 子任务的状态
			handleStatusUpdate(s, msg)
		case protocol.MsgHeartbeat:
			handleHealthCheck(s, msg) // 处理心跳消息
		default:
			log.Println("unhandled default case")
		}
		// 收到有效消息后重置超时计时器
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	}
}

func handleHealthCheck(s *Server, msg protocol.Message) {
	// 处理心跳消息
	agentID := msg.AgentID
	agent, ok := s.agents[agentID]
	log.Println("收到心跳消息:", msg)
	if !ok {
		log.Printf("未找到AgentID: %s", agentID)
		return
	}
	agent.LastActive = time.Now()
	agent.Status = protocol.AgentReady
}

// 处理agent的任务状态更新
func handleStatusUpdate(s *Server, msg protocol.Message) {
	// 处理Agent状态更新
	agentID := msg.AgentID
	// 反序列化消息
	var statusMsg schema.TaskStatus
	err := mapstructure.Decode(msg.Payload, &statusMsg)
	if err != nil {
		log.Printf("反序列化状态消息失败: %v", err)
		return
	}
	// 更新Agent状态和最后活跃时间
	connection, ok := s.agents[agentID]
	if !ok {
		log.Printf("未找到AgentID: %s", agentID)
		return
	}
	connection.LastActive = time.Now()
}

// FlowProcessor 是流程处理器的数据结构
type FlowProcessor struct {
	FlowID         string
	flowData       schema.FlowData
	executors      map[string]nodes.NodeExecutor
	stateStorage   Storage               // 存储状态
	statusChan     chan schema.TaskStep  // 把节点里的每个步骤的状态发送到这里
	nodeStatusChan chan protocol.Message // 把每个节点的状态发送到这里 payload是NodeState
}

// 初始化处理器
func NewFlowProcessor(flow schema.FlowData, server *Server) (*FlowProcessor, error) {
	fp := &FlowProcessor{
		FlowID:       generateFlowID(),
		flowData:     flow,
		executors:    make(map[string]nodes.NodeExecutor),
		stateStorage: NewFileStorage(""),
	}
	for _, node := range flow.Nodes {
		switch node.Type {
		case "start":
			fp.RegisterExecutor(node.ID, nodes.NewStartNodeExecutor(node))
		case "stop":
			fp.RegisterExecutor(node.ID, nodes.NewEndNodeExecutor(node))
		case "java":
			var data schema.JavaProperties
			// 从node.Properties中获取属性
			err := node.DeserializationProperties(data)
			if err != nil {
				return nil, err
			}
			host, ok := server.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("未找到AgentID: %s", data.Host)
			}
			if !server.HandleAgentStatus(data.Host) {
				return nil, fmt.Errorf("AgentID: %s 状态异常", data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewJavaNodeExecutor(data, host))
		case "build":
			fp.RegisterExecutor(node.ID, nodes.NewBuildNodeExecutor(node))
		case "web":
			// 反序列化properties
			var data schema.WebProperties
			// 从node.Properties中获取属性
			err := node.DeserializationProperties(data)
			if err != nil {
				return nil, err
			}
			host, ok := server.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("未找到AgentID: %s", data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewWebNodeExecuter(data, host))
		default:
			log.Printf("未注册的节点类型: %s", node.Type)
			return nil, fmt.Errorf("未注册的节点类型: %s", node.Type)
		}
	}
	return fp, nil
}

// 生成唯一的FlowID
func generateFlowID() string {
	return uuid.New().String()
}

// 接收statusChan中的数据 并将数据存到Storage中
func (fp *FlowProcessor) pushStatusUpdates() {

	for {
		log.Println("pushStatusUpdates waiting...")
		// 如果存储中没有flowExecution 就创建一个,每次更新的数据都从存储中获取
		flowExecution, ok := fp.stateStorage.Get(fp.FlowID)
		if !ok {
			flowExecution = schema.FlowExecution{
				FlowID:      fp.FlowID,
				NodeResults: make(map[string]schema.NodeState),
			}
		}
		select {
		// 收到状态更新
		case taskStep := <-fp.statusChan:
			// 把每一步的日志更新到flowExecution中
			taskStepData, ok := flowExecution.NodeResults[taskStep.NodeID]
			if !ok {
				// 如果没有当前node的执行日志就创建一个，并初始化它的数据
				taskStepData = schema.NodeState{
					ID:     taskStep.NodeID,
					Status: schema.NodeStateRunning,
					Logs:   "",
					Error:  "",
				}
				taskStepData.StartTime = utils.GetNowTime()
			}
			taskStepData.Logs += string(taskStep.Status) + taskStep.Output
			taskStepData.Error = taskStep.Error
			flowExecution.NodeResults[taskStep.NodeID] = taskStepData
			// 收到节点状态更新
		case state := <-fp.nodeStatusChan:
			// 提取关键状态参数信息
			var nodeStatus schema.NodeStatus
			err := mapstructure.Decode(state.Payload, &nodeStatus)
			if err != nil {
				log.Printf("反序列化状态消息失败: %v", err)
				return
			}
			// 更新node节点的状态
			nodeState := flowExecution.NodeResults[state.NodeID]
			nodeState.Status = nodeStatus
			fp.stateStorage.Save(flowExecution)
			// 节点执行成功，更新状态
			flowExecution.NodeResults[state.NodeID] = nodeState
		}
		fp.stateStorage.Save(flowExecution)

	}
}

// 执行flow
func (fp *FlowProcessor) ExecuteFlow(ctx context.Context, server *Server) schema.FlowExecution {
	execution := schema.FlowExecution{
		FlowID:      fp.FlowID,
		StartTime:   utils.GetNowTime(),
		EndTime:     nil,
		NodeResults: make(map[string]schema.NodeState),
	}
	//go fp.pushStatusUpdates(server)
	// 查找开始节点
	var startNode schema.Node
	for _, node := range fp.flowData.Nodes {
		if node.Type == "start" {
			startNode = node
			break
		}
	}
	currentNode := []schema.Node{startNode}

	// 执行节点
	for _, node := range currentNode {

		// 更新agent状态为处理中
		server.agents[node.ID].Status = protocol.TaskInProgress

		// 执行节点
		state := fp.executeNode(ctx, node, server)
		execution.NodeResults[node.ID] = state

		// 更新agent状态
		if state.Status == "success" {
			server.agents[node.ID].Status = protocol.AgentReady
		} else {
			server.agents[node.ID].Status = protocol.TaskCompleted
			break
		}

		// 获取下一个节点
		nextNode := NextNodes(fp.flowData, node)
		if nextNode == nil {
			break
		}
		currentNode = nextNode

	}
	execution.EndTime = utils.GetNowTime()
	fp.stateStorage.Save(execution)
	return execution
}

// 执行单个节点
func (fp *FlowProcessor) executeNode(ctx context.Context, node schema.Node, server *Server) schema.NodeState {
	// 等待agent准备就绪
	for {
		agent, exists := server.agents[node.ID]
		if !exists || agent.Status != protocol.AgentReady {
			time.Sleep(1 * time.Second)
			continue
		}
		flowExecution, ok := fp.stateStorage.Get(node.ID)
		if !ok {
			log.Printf("未注册的节点: %s", node.ID)
			return schema.NodeState{Status: "failed", Error: fmt.Sprintf("未注册的节点: %s", node.ID)}
		}
		if CheckDependency(fp.flowData, node, flowExecution) {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// 更新状态
	state := schema.NodeState{
		StartTime: utils.GetNowTime(),
	}

	executor, exists := fp.executors[node.ID]
	if !exists {
		state.Status = "failed"
		state.Error = fmt.Sprintf("未注册的node节点: %s", node.ID)
		state.EndTime = utils.GetNowTime()
		return state
	}
	executor.Execute(ctx, fp.statusChan)
	// 如果是start 或者 node 节点
	if node.Type == "start" || node.Type == "node" {
		state.EndTime = utils.GetNowTime()
		state.Status = schema.NodeStateSuccess
	} else {
		state.Status = schema.NodeStateRunning
	}
	return state
}

// 注册节点处理器
func (fp *FlowProcessor) RegisterExecutor(nodeId string, executor nodes.NodeExecutor) {
	fp.executors[nodeId] = executor
}
