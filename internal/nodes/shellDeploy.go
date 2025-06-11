package nodes

import (
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
)

type ShellDeployNode struct {
	msgChan chan interface{}
	agentID string
}

func NewShellDeployNode(agentID string, msgChan chan interface{}) *ShellDeployNode {
	return &ShellDeployNode{
		msgChan: msgChan,
		agentID: agentID,
	}
}

func (j *ShellDeployNode) Run(msg protocol.Message, task schema.ShellProperties) {
	data, _ := protocol.NewMessage(protocol.MsgTaskResult, msg.FlowExecutionID, j.agentID, msg.NodeID, schema.NodeStateSuccess)
	defer func() {
		log.Println("shell部署任务执行结束...")
		sendLastResult(j.msgChan, data)
	}()
	//
	// 初始化状态上报
	status := schema.NewTaskStep(msg.FlowExecutionID, j.agentID, msg.NodeID, "开始部署", schema.TaskStateRunning, "", "")
	steps := []struct {
		name   string
		action func() bool
	}{
		{
			"前置脚本",
			func() bool {
				return handleShellDeploy(status, "前置脚本", task.PreScriptContent, task.Timeout, j.msgChan)
			},
		},
		{
			"部署脚本",
			func() bool {
				return handleShellDeploy(status, "部署脚本", task.DeployScriptContent, task.Timeout, j.msgChan)
			},
		},
		{
			"后置脚本",
			func() bool {
				return handleShellDeploy(status, "后置脚本", task.PostScriptContent, task.Timeout, j.msgChan)
			},
		},
	}
	for _, step := range steps {
		if !step.action() {
			_ = data.UpdatePayload(schema.NodeStateRollbacked)
			return
		}
	}
}
