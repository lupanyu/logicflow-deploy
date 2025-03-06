package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"net"
	"sync"
	"time"

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
		agents:       make(map[string]*protocol.AgentConnection),
		stateStorage: NewMemoryStorage(),
		fpMap:        make(map[string]*FlowProcessor),
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
		log.Printf(" [%s]读取注册消息失败: %v", utils.GetCallerInfo(), err)
		return
	}
	data, _ := protocol.NewMessage(protocol.MsgRegisterResponse, registerMsg.FlowExecutionID, registerMsg.AgentID,
		registerMsg.NodeID, nil)

	// 验证消息格式
	if registerMsg.Type != protocol.MsgRegister {
		err := data.UpdatePayload(protocol.MessageAuthResponse{Code: 401, Message: "Invalid message type"})
		if err != nil {
			log.Printf(" [%s]创建消息失败: %v", generateFlowID(), err)
			return
		}
		_ = conn.WriteJSON(data)
		log.Printf(" [%s]收到无效消息: %+v", utils.GetCallerInfo(), registerMsg)
		return
	} else {
		if registerMsg.AgentID == "" {
			err := data.UpdatePayload(protocol.MessageAuthResponse{Code: 401, Message: "Invalid message type"})
			if err != nil {
				log.Printf(" [%s]创建消息失败: %v", generateFlowID(), err)
				return
			}
			log.Printf(" [%s]收到无效消息: %+v", generateFlowID(), registerMsg)
			if err = conn.Close(); err != nil {
				log.Printf(" [%s]关闭连接失败: %v", generateFlowID(), err)
			}
			return
		}
		err := data.UpdatePayload(protocol.MessageAuthResponse{Code: 200, Message: "注册成功"})
		if err != nil {
			log.Printf(" [%s]创建消息失败: %v", generateFlowID(), err)
			return
		}
		response, _ := protocol.NewMessage(protocol.MsgRegisterResponse, registerMsg.FlowExecutionID,
			registerMsg.AgentID, registerMsg.NodeID,
			protocol.MessageAuthResponse{Code: 200, Message: "注册成功"})
		conn.WriteJSON(response)
	}

	//把agent加入到server的map中
	s.agentsLock.Lock()
	agent := &protocol.AgentConnection{
		Conn:       conn,
		LastActive: time.Now(),
		Status:     protocol.AgentReady,
	}
	s.agents[registerMsg.AgentID] = agent
	log.Println("新添加的agent", registerMsg.AgentID, agent, s.agents)
	s.agentsLock.Unlock()

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf(" [%s]关闭%s连接失败: %v", utils.GetCallerInfo(), registerMsg.AgentID, err)
		} else {
			log.Printf(" [%s]关闭%s连接成功: %v", utils.GetCallerInfo(), registerMsg.AgentID, err)
		}
	}()
	log.Println("等待消息...")
	for {
		var msg protocol.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf(" [%s]读取%s消息超时:%v", utils.GetCallerInfo(), registerMsg.AgentID, err)
			} else {
				log.Printf(" [%s]读取%s消息失败: %v", utils.GetCallerInfo(), registerMsg.AgentID, err)
			}
			// 把agent从map中移除
			s.agentsLock.Lock()
			delete(s.agents, registerMsg.AgentID)
			s.agentsLock.Unlock()
			break
		}
		switch msg.Type {

		case protocol.MsgTaskStep:
			handleTaskStepUpdate(s, msg)
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
	fp := s.fpMap[msg.FlowExecutionID]

	fp.taskResultChan <- msg

	// 触发后面的节点
	flowExecution, ok := s.stateStorage.Get(msg.FlowExecutionID)
	if !ok {
		log.Printf(" [%s]flowExecutionID : %s", msg.FlowExecutionID)
		return
	}
	nodeID := msg.NodeID
	nextNodes := NextNodes(flowExecution.FlowData, nodeID)
	if len(nextNodes) == 0 {
		log.Printf(" [%s]flow: %s 没有剩余要执行的节点,结束", utils.GetCallerInfo(), flowExecution.FlowID)
	}
	for _, nextNode := range nextNodes {
		go fp.executeNode(nextNode, s)
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
	// 发送心跳响应
	response, _ := protocol.NewMessage(protocol.MsgHeartbeat, msg.FlowExecutionID, msg.AgentID, msg.NodeID, "pong")
	_ = conn.WriteJSON(response)

}

// 处理agent的任务状态更新
func handleTaskStepUpdate(s *Server, msg protocol.Message) {
	// 处理Agent状态更新
	var statusMsg schema.TaskStep
	err := protocol.UnMarshalPayload(msg.Payload, statusMsg)
	if err != nil {
		log.Printf(" [%s]反序列化状态消息失败: %v", utils.GetCallerInfo(), err)
		return
	}

	fp := s.fpMap[msg.FlowExecutionID]
	fp.taskStepChan <- statusMsg
}
