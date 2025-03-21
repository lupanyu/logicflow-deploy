package nodes

import (
	"fmt"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
)

type WebDeployNode struct {
	msgChan chan interface{}
	agentID string
}

func NewWebDeployNode(agentID string, msgChan chan interface{}) *WebDeployNode {
	return &WebDeployNode{
		msgChan: msgChan,
		agentID: agentID,
	}
}

func (w *WebDeployNode) Run(msg protocol.Message, task schema.WebProperties) {
	var rollbackFn []func()
	data, _ := protocol.NewMessage(protocol.MsgTaskResult, msg.FlowExecutionID, w.agentID, msg.NodeID, schema.NodeStateSuccess)
	defer func() {
		if len(rollbackFn) > 0 {
			fmt.Println("执行回滚操作...")
			for _, fn := range rollbackFn {
				fn()
			}
			err := data.UpdatePayload(schema.NodeStateRollbacked)
			if err != nil {
				fmt.Printf("[%s] 生成消息异常，错误是：%v", utils.GetCallerInfo(), err.Error())
			}
			sendLastResult(w.msgChan, data)
		}
	}()

	// 初始化状态上报
	status := schema.NewTaskStep(msg.FlowExecutionID, w.agentID, msg.NodeID, "开始部署", schema.TaskStateRunning, "", "")

	// 部署步骤集合（去掉了服务重启和健康检查）
	steps := []struct {
		name     string
		action   func() ([]byte, error)
		rollback func()
	}{
		{
			"备份旧版本",
			func() ([]byte, error) { return BakOld(task.DeployPath, task.BakPath) },
			nil,
		},
		{
			"下载最新代码包",
			func() ([]byte, error) { return UpdateFile(task.PackageSource, task.DeployPath) },
			nil,
		},
		{
			"解包",
			func() ([]byte, error) { return UnpackTar(task.DeployPath, task.DeployPath) },
			nil,
		},
	}

	for _, step := range steps {
		if step.rollback != nil {
			rollbackFn = append([]func(){step.rollback}, rollbackFn...)
		}
		if !handleStep(status, step.name, w.msgChan, step.action) {
			return
		}
	}
}
