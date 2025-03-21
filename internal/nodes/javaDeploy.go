package nodes

import (
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"time"
)

type JavaDeployNode struct {
	msgChan chan interface{}
	agentID string
}

func NewJavaDeployNode(agentID string, msgChan chan interface{}) *JavaDeployNode {
	return &JavaDeployNode{
		msgChan: msgChan,
		agentID: agentID,
	}
}

func (j *JavaDeployNode) Run(msg protocol.Message, task schema.JavaProperties) {
	var rollbackFn []func()
	data, _ := protocol.NewMessage(protocol.MsgTaskResult, msg.FlowExecutionID, j.agentID, msg.NodeID, schema.NodeStateSuccess)
	defer func() {
		log.Println("java部署任务执行结束...")
		// 发生错误时执行回滚操作
		if len(rollbackFn) > 0 {
			log.Println("执行回滚操作...")
			for _, fn := range rollbackFn {
				fn()
			}
			err := data.UpdatePayload(schema.NodeStateRollbacked)
			if err != nil {
				log.Printf("[%s]pdate payload err %v", utils.GetCallerInfo(), err)
			}
		}
		sendLastResult(j.msgChan, data)
	}()
	//
	// 初始化状态上报
	status := schema.NewTaskStep(msg.FlowExecutionID, j.agentID, msg.NodeID, "开始部署", schema.TaskStateRunning, "", "")

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
				return CheckAPPHealth(status, j.msgChan, task.Port, task.HealthUri, time.Duration(task.HealthCheckTimeout)*time.Second)
			},
			nil,
		},
	}

	for _, step := range steps {
		if step.rollback != nil {
			rollbackFn = append([]func(){step.rollback}, rollbackFn...)
		}
		if !handleStep(status, step.name, j.msgChan, step.action) {
			return
		}
	}
}
