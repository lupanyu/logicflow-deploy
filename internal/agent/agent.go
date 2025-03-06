package agent

import (
	"log"
	"logicflow-deploy/internal/nodes"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type DeploymentAgent struct {
	serverURL     string
	agentID       string
	wsConn        *websocket.Conn
	stopHeartbeat chan struct{} // 心跳停止信号
}

// 在结构体初始化时增加参数传递
func NewDeploymentAgent(serverURL string) *DeploymentAgent {
	hostname, _ := os.Hostname() // 获取主机名

	return &DeploymentAgent{
		serverURL: serverURL,
		agentID:   hostname, // 使用主机名作为 agentID
	}
}

func (a *DeploymentAgent) Run() {
	for {
		var msg protocol.Message
		if err := a.wsConn.ReadJSON(&msg); err != nil {
			// 处理连接错误类型
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf(" [%s]连接正常关闭", utils.GetCallerInfo())
				return
			}
			log.Printf(" [%s]连接异常: %v", utils.GetCallerInfo(), err)
			time.Sleep(5 * time.Second)
			a.reconnect()
			continue
		}
		// 验证消息格式
		if msg.Payload == nil {
			log.Printf(" [%s]收到无效消息: %+v", utils.GetCallerInfo(), msg)
			continue
		}
		switch msg.Type {
		case protocol.MsgRegisterResponse:
			// 处理注册消息
			// ...
			log.Printf(" [%s]收到注册响应消息: %+v", utils.GetCallerInfo(), msg)
		case protocol.MsgWebDeploy:
			go nodes.NewWebDeployNode(a.agentID, a.wsConn).Run(msg, msg.Payload.(schema.WebProperties))
		case protocol.MsgJavaDeploy:
			java := msg.Payload.(schema.JavaProperties)
			node := nodes.NewJavaDeployNode(a.agentID, a.wsConn)
			go node.Run(msg, java)
		case protocol.MsgHeartbeat:
			log.Printf(" [%s]收到心跳检测回应消息:%v\n", utils.GetCallerInfo(), msg)
		default:
			log.Printf(" [%s]未知消息类型: %s", utils.GetCallerInfo(), msg.Type)
			a.sendErrorResponse(msg.FlowExecutionID, "unsupported message type")
		}
	}
}

// 新增错误响应方法
func (a *DeploymentAgent) sendErrorResponse(taskID string, reason string) {
	event := protocol.Message{
		Type:            protocol.MsgError,
		FlowExecutionID: taskID,
		AgentID:         a.agentID,
		Timestamp:       time.Now().UnixNano(),
		Payload: schema.ErrorDetail{
			Code:    400,
			Message: reason,
		},
	}
	a.wsConn.WriteJSON(event)
}

func (a *DeploymentAgent) handleRollback(rollbackFn []func()) {
	if rollbackFn != nil {
		for _, fn := range rollbackFn {
			fn()
		}
	}
}

// 每10s发一次心跳检测
func (a *DeploymentAgent) Heartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	log.Println("心跳检测启动...")
	// 在退出时停止心跳

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartbeat := protocol.Message{
				Type:      protocol.MsgHeartbeat,
				AgentID:   a.agentID,
				Timestamp: time.Now().UnixNano(),
				Payload:   "ping",
			}
			if err := a.wsConn.WriteJSON(heartbeat); err != nil {
				log.Printf(" [%s]心跳发送失败: %v", utils.GetCallerInfo(), err)
				return
			} else {
				log.Printf(" [%s]心跳发送成功: %v", utils.GetCallerInfo(), heartbeat)
			}
		case <-a.stopHeartbeat:
			return
		}
	}
}

// 在现有结构体下方添加
func (a *DeploymentAgent) Connect() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(a.serverURL, nil)
	if err != nil {
		log.Printf("[%s] 连接服务器失败: %v", utils.GetCallerInfo(), err)
	}
	a.wsConn = conn
	// 发送注册消息
	registerMsg := protocol.Message{
		Type:      protocol.MsgRegister,
		AgentID:   a.agentID,
		Timestamp: time.Now().UnixNano(),
	}
	log.Printf(" [%s]发送注册消息: %+v", utils.GetCallerInfo(), registerMsg)
	if err := conn.WriteJSON(registerMsg); err != nil {
		conn.Close()
		log.Printf("[%s]注册消息发送失败: %v", utils.GetCallerInfo(), err)
	}

	log.Printf(" [%s]已连接服务器 %s [AgentID: %s]", utils.GetCallerInfo(), a.serverURL, a.agentID)
	// 心跳上报
	go a.Heartbeat()

	// 启动主循环
	a.Run()
}

// 在结构体中添加重连逻辑
func (a *DeploymentAgent) reconnect() {
	// ... 原有重连代码基础上添加日志 ...
	log.Printf(" [%s]尝试重新连接服务器...", utils.GetCallerInfo())
	a.Connect()
}
