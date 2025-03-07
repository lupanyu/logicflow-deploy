package nodes

import (
	"github.com/gorilla/websocket"
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
)

type ShellDeployNode struct {
	conn    *websocket.Conn
	agentID string
}

func NewShellDeployNode(agentID string, conn *websocket.Conn) *ShellDeployNode {
	return &ShellDeployNode{
		conn:    conn,
		agentID: agentID,
	}
}

func (j *ShellDeployNode) Run(msg protocol.Message, task schema.ShellProperties) {
	var rollbackFn []func()
	data, _ := protocol.NewMessage(protocol.MsgTaskResult, msg.FlowExecutionID, j.agentID, msg.NodeID, schema.NodeStateSuccess)
	defer func() {
		log.Println("shell部署任务执行结束...")
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
		sendLastResult(j.conn, data)
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
				return handleShellDeploy(status, "前置脚本", task.PreScriptContent, task.Timeout, j.conn)
			},
		},
		{
			"部署脚本",
			func() bool {
				return handleShellDeploy(status, "部署脚本", task.DeployScriptContent, task.Timeout, j.conn)
			},
		},
		{
			"后置脚本",
			func() bool {
				return handleShellDeploy(status, "后置脚本", task.PostScriptContent, task.Timeout, j.conn)
			},
		},
	}
	for _, step := range steps {
		if !step.action() {
			return
		}
	}
}
