package agent

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"log"
	"logicflow-deploy/internal/nodes"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type DeploymentAgent struct {
	serverURL     string
	agentID       string
	wsConn        *websocket.Conn
	stopHeartbeat chan struct{} // 心跳停止信号
	mu            sync.Mutex
	MsgChan       chan interface{}
	cancelWriter  context.CancelFunc // 取消函数
	writerCtx     context.Context    // 上下文
}

func (a *DeploymentAgent) Send(msg interface{}) {
	log.Printf(" [%s]发送消息: %+v", utils.GetCallerInfo(), msg)
	a.MsgChan <- msg
}

func (a *DeploymentAgent) WaitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	err := a.wsConn.Close()
	if err != nil {
		return
	}
}

// NewDeploymentAgent 在结构体初始化时增加参数传递
func NewDeploymentAgent(serverURL string) *DeploymentAgent {
	hostname, _ := os.Hostname()

	return &DeploymentAgent{
		serverURL:     serverURL,
		agentID:       hostname,                  // 使用主机名作为 agentID
		MsgChan:       make(chan interface{}, 1), // 初始化消息通道, 容量为1
		mu:            sync.Mutex{},
		stopHeartbeat: make(chan struct{}),
	}
}

func (a *DeploymentAgent) WriteToConn() {
	for {
		select {
		case msg := <-a.MsgChan:
			a.mu.Lock()
			err := a.wsConn.WriteJSON(msg)
			a.mu.Unlock()

			if err != nil {
				log.Printf("写入失败: %v", err)
				a.reconnect()
				return // 退出当前写协程
			}

		case <-a.writerCtx.Done():
			return // 上下文取消时退出
		}
	}
}

func (a *DeploymentAgent) Run() {
	for {
		var msg protocol.Message
		if err := a.wsConn.ReadJSON(&msg); err != nil {
			// 关闭所有写协程
			if a.cancelWriter != nil {
				a.cancelWriter()
			}
			// 处理连接错误类型
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf(" [%s]连接正常关闭", utils.GetCallerInfo())

			}
			log.Printf(" [%s]连接异常: %v", utils.GetCallerInfo(), err)
			go a.reconnect()
			return // 退出当前协程
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
			var web schema.WebProperties
			if err := json.Unmarshal(msg.Payload, &web); err != nil {
				log.Printf(" [%s]解析Web部署消息失败: %v", utils.GetCallerInfo(), err)
				continue
			}
			log.Printf(" [%s]解析Web部署消息成功: %v", utils.GetCallerInfo(), web)
			go nodes.NewWebDeployNode(a.agentID, a.MsgChan).Run(msg, web)
		case protocol.MsgJavaDeploy:
			var java schema.JavaProperties
			if err := json.Unmarshal(msg.Payload, &java); err != nil {
				log.Printf(" [%s]解析Java部署消息失败: %v", utils.GetCallerInfo(), err)
				continue
			}
			log.Printf(" [%s]解析Java部署消息成功: %v", utils.GetCallerInfo(), java)
			node := nodes.NewJavaDeployNode(a.agentID, a.MsgChan)
			go node.Run(msg, java)
		case protocol.MsgShellDeploy:
			var shell schema.ShellProperties
			if err := json.Unmarshal(msg.Payload, &shell); err != nil {
				log.Printf(" [%s]解析Shell部署消息失败: %v", utils.GetCallerInfo(), err)
				continue
			}
			node := nodes.NewShellDeployNode(a.agentID, a.MsgChan)
			// 设置默认的超时时间为10分钟
			if shell.Timeout == 0 {
				shell.Timeout = 600
			}
			log.Printf(" [%s]解析Shell部署消息成功: %v", utils.GetCallerInfo(), shell)
			go node.Run(msg, shell)
		case protocol.MsgHeartbeat:
			log.Printf(" [%s]收到心跳检测回应消息:%v\n", utils.GetCallerInfo(), msg)
		default:
			log.Printf(" [%s]未知消息类型: %d", utils.GetCallerInfo(), msg.Type)
			a.sendErrorResponse(msg.FlowExecutionID, msg.NodeID, "unsupported message type")
		}
	}
}

// 新增错误响应方法
func (a *DeploymentAgent) sendErrorResponse(taskID, nodeID string, reason string) {
	payload := schema.ErrorDetail{
		Code:    400,
		Message: reason,
	}

	event, err := protocol.NewMessage(protocol.MsgTaskResult, taskID, a.agentID, nodeID, payload)
	if err != nil {
		log.Printf(" [%s]创建错误响应消息失败: %v", utils.GetCallerInfo(), err)
		return
	}
	a.Send(event)
}

func (a *DeploymentAgent) handleRollback(rollbackFn []func()) {
	if rollbackFn != nil {
		for _, fn := range rollbackFn {
			fn()
		}
	}
}

// Heartbeat 每10s发一次心跳检测
func (a *DeploymentAgent) Heartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	log.Println("心跳检测启动...")
	jump := 0
	// 在退出时停止心跳
	defer ticker.Stop()
	for {
		jump++
		select {
		case <-ticker.C:
			if a.writerCtx.Err() != nil {
				return
			}
			heartbeat, err := protocol.NewMessage(protocol.MsgHeartbeat, "", a.agentID, fmt.Sprintf("%d", jump), "ping")
			if err != nil {
				log.Printf(" [%s]心跳消息创建失败: %v", utils.GetCallerInfo(), err)
				continue
			}
			a.Send(heartbeat)
		case <-a.stopHeartbeat:
			return
		}
	}
}

// Connect 在现有结构体下方添加
func (a *DeploymentAgent) Connect() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 清理旧连接
	if a.wsConn != nil {
		_ = a.wsConn.Close()
		a.cancelWriter()
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancelWriter = cancel
	a.writerCtx = ctx

	// 统一连接逻辑
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(a.serverURL, nil)
	if err != nil {
		log.Printf("[%s] 连接服务器失败: %v", utils.GetCallerInfo(), err)
		return err
	}
	a.wsConn = conn
	log.Printf(" [%s]连接服务器 %s [AgentID: %s]", utils.GetCallerInfo(), a.serverURL, a.agentID)

	// 公共初始化逻辑
	go a.WriteToConn()
	a.sendRegister()
	go a.Heartbeat()
	go a.Run()
	return nil
}

// 提取注册消息发送逻辑
func (a *DeploymentAgent) sendRegister() {
	registerMsg := protocol.Message{
		Type:      protocol.MsgRegister,
		AgentID:   a.agentID,
		Timestamp: time.Now().UnixNano(),
	}
	log.Printf(" [%s]发送注册消息: %+v", utils.GetCallerInfo(), registerMsg)
	a.Send(registerMsg)
}

// 在结构体中添加重连逻辑
func (a *DeploymentAgent) reconnect() {
	a.mu.Lock()
	defer a.mu.Unlock()
	// ... 原有重连代码基础上添加日志 ...
	log.Printf(" [%s]尝试重新连接服务器...", utils.GetCallerInfo())
	retry := 1
	// 添加重试控制
	for {

		if a.wsConn != nil {
			a.wsConn.Close()
		}
		// 指数退避策略
		waitTime := time.Duration(retry*retry) * time.Second
		if waitTime > 30*time.Second {
			waitTime = 30 * time.Second
		}
		log.Printf(" [%s]等待%.0f秒后尝试第%d次重连...", utils.GetCallerInfo(), waitTime.Seconds(), retry)
		time.Sleep(waitTime)

		// 创建新连接
		err := a.Connect()
		if err == nil {
			log.Printf(" [%s]成功重连服务器", utils.GetCallerInfo())
			return
		}
		log.Printf(" [%s]连接失败: %v", utils.GetCallerInfo(), err)
		retry++
	}
}
