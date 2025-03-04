package agent

import (
	"fmt"
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
				log.Printf(" [%s]连接正常关闭", utils.GetCallerInfo(1))
				return
			}
			log.Printf(" [%s]连接异常: %v", utils.GetCallerInfo(1), err)
			time.Sleep(5 * time.Second)
			a.reconnect()
			continue
		}
		// 验证消息格式
		if msg.Payload == nil {
			log.Printf(" [%s]收到无效消息: %+v", utils.GetCallerInfo(1), msg)
			continue
		}
		switch msg.Type {
		case protocol.MsgRegisterResponse:
			// 处理注册消息
			// ...
			log.Printf(" [%s]收到注册响应消息: %+v", utils.GetCallerInfo(1), msg)
		case protocol.MsgWebDeploy:
			go nodes.NewWebDeployNode(a.agentID, a.wsConn).Run(msg, msg.Payload.(schema.WebProperties))
		case protocol.MsgJavaDeploy:
			go nodes.NewJavaDeployNode(a.agentID, a.wsConn).Run(msg, msg.Payload.(schema.JavaProperties))
		case protocol.MsgHeartbeat:
			log.Printf(" [%s]收到心跳检测回应消息:%v\n", utils.GetCallerInfo(1), msg)
		default:
			log.Printf(" [%s]未知消息类型: %s", utils.GetCallerInfo(1), msg.Type)
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
				Payload:   []byte("ping"),
			}
			if err := a.wsConn.WriteJSON(heartbeat); err != nil {
				log.Printf(" [%s]心跳发送失败: %v", utils.GetCallerInfo(1), err)
				return
			} else {
				log.Printf(" [%s]心跳发送成功: %v", utils.GetCallerInfo(1), heartbeat)
			}
		case <-a.stopHeartbeat:
			return
		}
	}
}

// 在现有结构体下方添加
func (a *DeploymentAgent) Connect() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(a.serverURL, nil)
	if err != nil {
		return fmt.Errorf("连接服务器失败: %w", err)
	}
	a.wsConn = conn
	// 发送注册消息
	registerMsg := protocol.Message{
		Type:      protocol.MsgRegister,
		AgentID:   a.agentID,
		Timestamp: time.Now().UnixNano(),
	}
	log.Printf(" [%s]发送注册消息: %+v", utils.GetCallerInfo(1), registerMsg)
	if err := conn.WriteJSON(registerMsg); err != nil {
		conn.Close()
		return fmt.Errorf("注册消息发送失败: %w", err)
	}

	log.Printf(" [%s]已连接服务器 %s [AgentID: %s]", utils.GetCallerInfo(1), a.serverURL, a.agentID)
	return nil
}

// 在结构体中添加重连逻辑
func (a *DeploymentAgent) reconnect() {
	// ... 原有重连代码基础上添加日志 ...
	log.Printf(" [%s]尝试重新连接服务器...", utils.GetCallerInfo(1))
	if err := a.Connect(); err != nil {
		log.Printf(" [%s]重连失败: %v", utils.GetCallerInfo(1), err)
	}
	// 心跳上报
	go a.Heartbeat()

	// 启动主循环
	a.Run()
}
