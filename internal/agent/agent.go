package agent

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"net/http"
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
				log.Printf("连接正常关闭")
				return
			}
			log.Printf("连接异常: %v", err)
			//a.reconnect()
			continue
		}
		// 验证消息格式
		if msg.Payload == nil {
			log.Printf("收到无效消息: %+v", msg)
			continue
		}
		switch msg.Type {
		case protocol.MsgRegisterResponse:
			// 处理注册消息
			// ...
			log.Printf("收到注册响应消息: %+v", msg)
		case protocol.MsgJavaDeploy:
			go a.handleJavaTask(msg)
		case protocol.MsgHeartbeat:
			log.Printf("收到心跳检测回应消息:%v\n", msg)
		default:
			log.Printf("未知消息类型: %s", msg.Type)
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

// 错误处理闭包
func (a *DeploymentAgent) handleStep(nodeId, stepName string, fn func() ([]byte, error)) bool {
	status := schema.NewTaskStep(a.agentID, nodeId, stepName, schema.TaskStateSuccess, "", "")
	out, err := fn()
	if err != nil {
		status.Status = schema.TaskStateFailed
		status.Output = string(out)
		status.Error = err.Error()
		//a.sendStatus(*status)
		return false
	}
	status.Output = string(out)
	//a.sendStatus(*status)
	return true
}
func (a *DeploymentAgent) handleRollback(rollbackFn []func()) {
	if rollbackFn != nil {
		for _, fn := range rollbackFn {
			fn()
		}
	}
}

func (a *DeploymentAgent) handleJavaTask(msg protocol.Message) {
	var task schema.JavaProperties
	err := mapstructure.Decode(msg.Payload, &task)
	if err != nil {
		log.Printf("解析任务数据失败: %v", err)
		MsgTaskResult := protocol.Message{
			Type:            protocol.MsgTaskResult,
			FlowExecutionID: msg.FlowExecutionID,
			AgentID:         a.agentID,
			NodeID:          msg.NodeID,
			Timestamp:       time.Now().UnixNano(),
			Payload:         schema.NodeStateSuccess,
		}
		var rollbackFn = new([]func())
		defer func() {
			// 执行回滚操作
			fmt.Println("在defer里...")
			MsgTaskResult.Timestamp = time.Now().UnixNano()
			if len(*rollbackFn) != 0 {
				for _, fn := range *rollbackFn {
					fmt.Println("执行回滚...")
					fn()
					MsgTaskResult.Payload = schema.NodeStateRollbacked
				}
			}
			a.sendLastResult(MsgTaskResult)
		}()
		// 执行部署步骤
		status := schema.NewTaskStep("开始部署", msg.AgentID, msg.NodeID, schema.TaskStateRunning, "", "")
		a.sendStatus(*status)

		*rollbackFn = append(*rollbackFn, func() {
			a.StartService(task.ServerName)
		})
		// 1. 停止服务
		if !a.handleStep("停止服务", msg.NodeID, func() ([]byte, error) {
			return a.StopService(task.ServerName)
		}) {
			MsgTaskResult.Payload = schema.NodeStateFailed
			a.sendLastResult(MsgTaskResult)
			return
		}

		// 2. 备份旧版本
		if !a.handleStep("备份旧版本", msg.NodeID, func() ([]byte, error) {
			return a.BakOld(task.DeployPath, task.BakPath)
		}) {
			MsgTaskResult.Payload = schema.NodeStateFailed
			a.sendLastResult(MsgTaskResult)
			return
		}

		// 恢复原始状态
		newRollbackFn := func() { a.Rollback(task.BakPath, task.DeployPath, task.ServerName) }
		*rollbackFn = append([]func(){newRollbackFn}, *rollbackFn...)
		// 3. 下载文件
		if !a.handleStep("下载最新代码包", msg.NodeID, func() ([]byte, error) {
			return a.UpdateFile(task.PackageSource, task.DeployPath)
		}) {
			MsgTaskResult.Payload = schema.NodeStateFailed
			a.sendLastResult(MsgTaskResult)
			return
		}

		// 4. 启动服务
		if !a.handleStep("启动服务", msg.NodeID, func() ([]byte, error) {
			return a.StartService(task.ServerName)
		}) {
			MsgTaskResult.Payload = schema.NodeStateFailed
			a.sendLastResult(MsgTaskResult)
			return
		}

		// 5. 健康检查
		if !a.handleStep("健康检查", msg.NodeID, func() ([]byte, error) {
			return a.checkHealth(msg.NodeID, task.Port, task.HealthUri, time.Duration(task.HealthCheckTimeout)*time.Second)
		}) {
			MsgTaskResult.Payload = schema.NodeStateFailed
			a.sendLastResult(MsgTaskResult)
			return
		}

		// 6. 清理旧版本
	}
}
func (a *DeploymentAgent) StopService(service string) ([]byte, error) {
	return utils.RunShell("systemctl stop " + service)
}

// 正确写法应该像其他方法一样有空格：
func (a *DeploymentAgent) StartService(service string) ([]byte, error) {
	return utils.RunShell("systemctl start " + service)
}

func (a *DeploymentAgent) BakOld(old, new string) ([]byte, error) {

	return utils.RunShell("cp -r " + old + " " + new)
}

func (a *DeploymentAgent) UpdateFile(downloadURL, new string) ([]byte, error) {
	return utils.RunShell("curl -o " + new + " " + downloadURL)
}

func (a *DeploymentAgent) Rollback(backup, new, service string) ([]byte, error) {
	return utils.RunShell("cp -r " + backup + " " + new)
}

func (a *DeploymentAgent) checkHealth(nodeId string, port int, uri string, timeout time.Duration) ([]byte, error) {
	status := schema.NewTaskStep(a.agentID, nodeId, "健康检查", schema.TaskStateSuccess, "", "")
	defer a.sendStatus(*status)
	client := http.Client{Timeout: 3 * time.Second}
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%d%s", port, uri))
		if err == nil && resp.StatusCode == 200 {
			status.Status = schema.TaskStateSuccess
			return nil, nil
		}
		time.Sleep(5 * time.Second)
	}
	status.Status = schema.TaskStateFailed
	status.Error = "健康检查超时"
	return nil, errors.New("健康检查超时")
}

func (a *DeploymentAgent) sendStatus(status schema.TaskStep) {
	event := protocol.Message{
		Type:            protocol.MsgTaskStep,
		FlowExecutionID: status.FlowExecutionID,
		AgentID:         a.agentID,
		Timestamp:       time.Now().UnixNano(),
		Payload:         status,
	}
	a.wsConn.WriteJSON(event)
}

func (a *DeploymentAgent) sendLastResult(data protocol.Message) {
	data.Timestamp = time.Now().UnixNano()
	a.wsConn.WriteJSON(data)
}
func (a *DeploymentAgent) sendLastErrorResult(data protocol.Message) {
	data.Timestamp = time.Now().UnixNano()
	data.Payload =
		a.wsConn.WriteJSON(data)
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
				log.Printf("心跳发送失败: %v", err)
				return
			} else {
				log.Printf("心跳发送成功: %v", heartbeat)
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
	log.Printf("发送注册消息: %+v", registerMsg)
	if err := conn.WriteJSON(registerMsg); err != nil {
		conn.Close()
		return fmt.Errorf("注册消息发送失败: %w", err)
	}

	log.Printf("已连接服务器 %s [AgentID: %s]", a.serverURL, a.agentID)
	return nil
}

// 在结构体中添加重连逻辑
func (a *DeploymentAgent) reconnect() {
	// ... 原有重连代码基础上添加日志 ...
	log.Printf("尝试重新连接服务器...")
	if err := a.Connect(); err != nil {
		log.Printf("重连失败: %v", err)
	}
}
