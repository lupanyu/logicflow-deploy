package nodes

import (
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
)

type StartNodeExecutor struct {
}

func (e *StartNodeExecutor) Execute(flowExecutionID, nodeID string, ch chan schema.TaskStep, result chan protocol.Message) {

	stat, _ := protocol.NewMessage(protocol.MsgTaskResult, flowExecutionID, "start", nodeID, schema.NodeStateSuccess)
	result <- stat
}

func (e *StartNodeExecutor) NodeType() string {
	return "start"
}
func (e *StartNodeExecutor) AgentId() string {
	return "start"
}

func NewStartNodeExecutor(node schema.Node) *StartNodeExecutor {

	return &StartNodeExecutor{}
}
