package nodes

import (
	"fmt"
	"github.com/gorilla/websocket"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"time"
)

type JavaDeployNode struct {
	conn    *websocket.Conn
	agentID string
}

func NewJavaDeployNode(agentID string, conn *websocket.Conn) *JavaDeployNode {
	return &JavaDeployNode{
		conn:    conn,
		agentID: agentID,
	}
}

func (j *JavaDeployNode) Run(msg protocol.Message, task schema.JavaProperties) {
	var rollbackFn []func()
	defer func() {
		if len(rollbackFn) > 0 {
			fmt.Println("执行回滚操作...")
			for _, fn := range rollbackFn {
				fn()
			}
			sendLastResult(j.conn, protocol.Message{
				Type:            protocol.MsgTaskResult,
				FlowExecutionID: msg.FlowExecutionID,
				AgentID:         j.agentID,
				NodeID:          msg.NodeID,
				Timestamp:       time.Now().UnixNano(),
				Payload:         schema.NodeStateRollbacked,
			})
		}
	}()

	// 初始化状态上报
	status := schema.NewTaskStep(j.agentID, msg.NodeID, "开始部署", schema.TaskStateRunning, "", "")
	sendStatus(j.conn, j.agentID, *status)

	// 部署步骤集合
	steps := []struct {
		name     string
		action   func() ([]byte, error)
		rollback func()
	}{
		{
			"停止服务",
			func() ([]byte, error) { return StopService(task.ServerName) },
			func() { StartService(task.ServerName) },
		},
		{
			"备份旧版本",
			func() ([]byte, error) { return BakOld(task.DeployPath, task.BakPath) },
			nil,
		},
		{
			"下载最新代码包",
			func() ([]byte, error) { return UpdateFile(task.PackageSource, task.DeployPath) },
			func() { Rollback(task.BakPath, task.DeployPath, task.ServerName) },
		},
		{
			"启动服务",
			func() ([]byte, error) { return StartService(task.ServerName) },
			nil,
		},
		{
			"健康检查",
			func() ([]byte, error) {
				return CheckAPPHealth(j.agentID, j.conn, msg.NodeID, task.Port, task.HealthUri, time.Duration(task.HealthCheckTimeout)*time.Second)
			},
			nil,
		},
	}

	for _, step := range steps {
		if !handleStep(j.agentID, step.name, msg.NodeID, step.action) {
			return
		}
		if step.rollback != nil {
			rollbackFn = append([]func(){step.rollback}, rollbackFn...)
		}
	}
}

//func (j *JavaDeployNode) handleJavaTask(msg protocol.Message) {
//	var task schema.JavaProperties
//	err := mapstructure.Decode(msg.Payload, &task)
//	if err != nil {
//		log.Printf(" [%s]解析任务数据失败: %v", err)
//		MsgTaskResult := protocol.Message{
//			Type:            protocol.MsgTaskResult,
//			FlowExecutionID: msg.FlowExecutionID,
//			AgentID:         j.agentID,
//			NodeID:          msg.NodeID,
//			Timestamp:       time.Now().UnixNano(),
//			Payload:         schema.NodeStateSuccess,
//		}
//		var rollbackFn = new([]func())
//		defer func() {
//			// 执行回滚操作
//			fmt.Println("在defer里...")
//			MsgTaskResult.Timestamp = time.Now().UnixNano()
//			if len(*rollbackFn) != 0 {
//				for _, fn := range *rollbackFn {
//					fmt.Println("执行回滚...")
//					fn()
//					MsgTaskResult.Payload = schema.NodeStateRollbacked
//				}
//			}
//			sendLastResult(j.conn, MsgTaskResult)
//		}()
//		// 执行部署步骤
//		status := schema.NewTaskStep(msg.AgentID, msg.NodeID, "开始部署", schema.TaskStateRunning, "", "")
//		sendStatus(j.conn, j.agentID, *status)
//
//		*rollbackFn = append(*rollbackFn, func() {
//			StartService(task.ServerName)
//		})
//		// 1. 停止服务
//		if !handleStep(j.agentID, "停止服务", msg.NodeID, func() ([]byte, error) {
//			return StopService(task.ServerName)
//		}) {
//			MsgTaskResult.Payload = schema.NodeStateFailed
//			sendLastResult(j.conn, MsgTaskResult)
//			return
//		}
//
//		// 2. 备份旧版本
//		if !handleStep(j.agentID, "备份旧版本", msg.NodeID, func() ([]byte, error) {
//			return BakOld(task.DeployPath, task.BakPath)
//		}) {
//			MsgTaskResult.Payload = schema.NodeStateFailed
//			sendLastResult(j.conn, MsgTaskResult)
//			return
//		}
//
//		// 恢复原始状态
//		newRollbackFn := func() { Rollback(task.BakPath, task.DeployPath, task.ServerName) }
//		*rollbackFn = append([]func(){newRollbackFn}, *rollbackFn...)
//		// 3. 下载文件
//		if !handleStep(j.agentID, "下载最新代码包", msg.NodeID, func() ([]byte, error) {
//			return UpdateFile(task.PackageSource, task.DeployPath)
//		}) {
//			MsgTaskResult.Payload = schema.NodeStateFailed
//			sendLastResult(j.conn, MsgTaskResult)
//			return
//		}
//
//		// 4. 启动服务
//		if !handleStep(j.agentID, "启动服务", msg.NodeID, func() ([]byte, error) {
//			return StartService(task.ServerName)
//		}) {
//			MsgTaskResult.Payload = schema.NodeStateFailed
//			sendLastResult(j.conn, MsgTaskResult)
//			return
//		}
//
//		// 5. 健康检查
//		if !handleStep(j.agentID, "健康检查", msg.NodeID, func() ([]byte, error) {
//			return CheckAPPHealth(j.agentID, j.conn, msg.NodeID, task.Port, task.HealthUri, time.Duration(task.HealthCheckTimeout)*time.Second)
//		}) {
//			MsgTaskResult.Payload = schema.NodeStateFailed
//			sendLastResult(j.conn, MsgTaskResult)
//			return
//		}
//
//		// 6. 清理旧版本
//	}
//}

// ... 保留其他原有方法 ...
