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
	agents       map[string]*protocol.AgentConnection // 连接的Agent
	stateStorage Storage
	fpMap        map[string]*FlowProcessor // 当前在执行的flow的处理器,key 是 flowID
	httpServer   *gin.Engine
	agentsLock   sync.RWMutex
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
	_ = s.httpServer.Run(fmt.Sprintf("%s:%d", ip, port))
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
		log.Printf(" [%s]读取注册消息失败: %v", err)
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
		log.Printf(" [%s]收到无效消息: %+v", registerMsg)
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
			log.Printf(" [%s]收到无效消息: %+v", registerMsg)
			if err := conn.Close(); err != nil {
				log.Printf(" [%s]关闭连接失败: %v", err)
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
	agent := &protocol.AgentConnection{
		Conn:       conn,
		LastActive: time.Now(),
		Status:     protocol.AgentReady,
	}
	s.agents[agentID] = agent
	log.Println("新添加的agent", agentID, agent, s.agents)
	s.agentsLock.Unlock()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf(" [%s]关闭%s连接失败: %v", agentID, err)
		} else {
			log.Printf(" [%s]关闭%s连接成功: %v", agentID, err)
		}
	}()
	log.Println("等待消息...")
	for {
		var msg protocol.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf(" [%s]读取%s消息超时:%v", agentID, err)
			} else {
				log.Printf(" [%s]读取%s消息失败: %v", agentID, err)
			}
			// 把agent从map中移除
			s.agentsLock.Lock()
			delete(s.agents, agentID)
			s.agentsLock.Unlock()
			break
		}
		switch msg.Type {
		case protocol.MsgStatus: // 处理最终的任务状态
			handleStatusUpdate(s, msg)
		case protocol.MsgTaskStep:
			handleStatusUpdate(s, msg)
		case protocol.MsgTaskResult:
			handleTaskResult(s, msg)
		case protocol.MsgHeartbeat:
			handleHealthCheck(s, msg, conn) // 处理心跳消息

		default:
			log.Println("unhandled default case", msg)
		}
		// 收到有效消息后重置超时计时器
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	}
}

func handleTaskResult(s *Server, msg protocol.Message) {
	log.Printf("[%s]收到 %s 任务结果: %+v", msg.AgentID, msg.NodeID, msg.Payload)
	// 更新任务结果到存储中
	flowExecution, ok := s.stateStorage.Get(msg.FlowExecutionID)
	if !ok {
		log.Printf(" [%s]flowExecutionID : %s", msg.FlowExecutionID)
		return
	}
	nodeID := msg.NodeID

	lastState, ok := msg.Payload.(schema.NodeStatus)
	if !ok {
		log.Printf(" [%s]反序列化状态消息失败: %v", utils.GetCallerInfo(), msg.Payload)
		return
	}
	nodeResult := flowExecution.NodeResults[nodeID]
	nodeResult.Status = lastState
	nodeResult.EndTime = utils.GetNowTime()
	flowExecution.NodeResults[nodeID] = nodeResult
	s.stateStorage.Save(flowExecution)
	// 触发后面的节点
	nextNodes := NextNodes(flowExecution.FlowData, nodeID)
	if len(nextNodes) == 0 {
		log.Printf(" [%s]flow: %s 没有剩余要执行的节点,结束", utils.GetCallerInfo(), flowExecution.FlowID)
	}
	fp := s.fpMap[flowExecution.FlowID]
	for _, nextNode := range nextNodes {
		go fp.executeNode(context.Background(), nextNode, s)
	}
}

func handleHealthCheck(s *Server, msg protocol.Message, conn *websocket.Conn) {
	// 处理心跳消息
	agentID := msg.AgentID
	agent, ok := s.agents[agentID]
	log.Println("收到心跳消息:", msg)
	if !ok {
		log.Printf(" [%s]未找到AgentID: %s", agentID)
		return
	}
	agent.LastActive = time.Now()

	// 回复心跳响应
	response := protocol.Message{
		Type:      protocol.MsgHeartbeat,
		AgentID:   agentID,
		Timestamp: time.Now().UnixNano(),
		Payload:   "pong",
	}
	conn.WriteJSON(response)

}

// 处理agent的任务状态更新
func handleStatusUpdate(s *Server, msg protocol.Message) {
	// 处理Agent状态更新
	agentID := msg.AgentID
	// 反序列化消息
	var statusMsg schema.TaskStatus
	err := mapstructure.Decode(msg.Payload, &statusMsg)
	if err != nil {
		log.Printf(" [%s]反序列化状态消息失败: %v", err)
		return
	}
	// 更新Agent状态和最后活跃时间
	connection, ok := s.agents[agentID]
	if !ok {
		log.Printf(" [%s]未找到AgentID: %s", agentID)
		return
	}
	connection.LastActive = time.Now()
}

// FlowProcessor 是流程处理器的数据结构
type FlowProcessor struct {
	FlowID         string
	flowData       schema.FlowData
	executors      map[string]nodes.NodeExecutor
	statusChan     chan schema.TaskStep  // 把节点里的每个步骤的状态发送到这里
	nodeStatusChan chan protocol.Message // 把每个节点的状态发送到这里 payload是NodeState
}

//

// 初始化处理器
func NewFlowProcessor(flow schema.FlowData, s *Server) (*FlowProcessor, error) {
	fp := &FlowProcessor{
		FlowID:    generateFlowID(),
		flowData:  flow,
		executors: make(map[string]nodes.NodeExecutor),
		//stateStorage: NewFileStorage(""),
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
			err := node.DeserializationProperties(&data)
			if err != nil {
				return nil, err
			}
			log.Println("反序列化properties", data)

			host, ok := s.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("java节点%v未找到AgentID: %s", data, data.Host)
			}
			if !s.HandleAgentStatus(data.Host) {
				return nil, fmt.Errorf("AgentID: %s 状态异常", data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewJavaNodeExecutor(data, host))
		case "build":
			fp.RegisterExecutor(node.ID, nodes.NewBuildNodeExecutor(node))
		case "web":
			// 反序列化properties
			var data schema.WebProperties
			// 从node.Properties中获取属性

			err := node.DeserializationProperties(&data)
			if err != nil {
				return nil, err
			}
			log.Println("反序列化properties", data)

			host, ok := s.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("web节点%v未找到AgentID: %s", data, data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewWebNodeExecuter(data, host))
		case "end":
			fp.RegisterExecutor(node.ID, nodes.NewEndNodeExecutor(node))
		default:
			log.Printf(" [%s]未注册的节点类型: %s", utils.GetCallerInfo(), node.Type)
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
func (fp *FlowProcessor) pushStatusUpdates(mem Storage) {

	for {
		log.Println("pushStatusUpdates waiting...")
		// 如果存储中没有flowExecution 就创建一个,每次更新的数据都从存储中获取
		flowExecution, ok := mem.Get(fp.FlowID)
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
				log.Printf(" [%s]反序列化状态消息失败: %v", err)
				return
			}
			// 更新node节点的状态
			nodeState := flowExecution.NodeResults[state.NodeID]
			nodeState.Status = nodeStatus
			mem.Save(flowExecution)
			// 节点执行成功，更新状态
			flowExecution.NodeResults[state.NodeID] = nodeState
		}
		mem.Save(flowExecution)

	}
}

// 执行flow
func (fp *FlowProcessor) ExecuteFlow(ctx context.Context, server *Server) schema.FlowExecution {
	execution := schema.FlowExecution{
		FlowID:      fp.FlowID,
		StartTime:   utils.GetNowTime(),
		EndTime:     nil,
		NodeResults: make(map[string]schema.NodeState),
		FlowData:    fp.flowData,
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

	// 执行节点
	go func() {
		// ... 原有上下文代码 ...

		log.Printf("[%s] 执行节点:%v", utils.GetCallerInfo(), startNode.ID)
		fp.executeNode(ctx, startNode, server)

	}()
	//execution.EndTime = utils.GetNowTime()
	server.stateStorage.Save(execution)
	log.Println("执行完成", fp.executors)
	return execution
}

// 执行单个节点
func (fp *FlowProcessor) executeNode(ctx context.Context, node schema.Node, server *Server) {

	// 查看存储中是否有当前节点的执行日志
	flowExecution, ok := server.stateStorage.Get(fp.FlowID)
	if !ok {
		log.Printf(" [%s]未找到FlowExecution: %s,停止执行", utils.GetCallerInfo(), fp.FlowID)
		return
	}
	nodeState, ok := flowExecution.NodeResults[node.ID]
	if !ok {
		nodeState = schema.NodeState{
			ID:     node.ID,
			Status: schema.NodeStateRunning,
			Logs:   "",
			Error:  "",
		}
		nodeState.StartTime = utils.GetNowTime()
	}
	if CheckDependency(fp.flowData, node, flowExecution) {
		log.Println("依赖节点执行成功")
	} else {
		log.Println("依赖节点未执行成功，停止执行当前节点", node.ID)
		return
	}

	executor, exists := fp.executors[node.ID]
	if !exists {
		log.Printf(" [%s]未找到节点执行器: %s", utils.GetCallerInfo(), node.ID)
		return
	}
	executor.Execute(ctx, fp.statusChan)

}

// 注册节点处理器
func (fp *FlowProcessor) RegisterExecutor(nodeId string, executor nodes.NodeExecutor) {
	fp.executors[nodeId] = executor
}
