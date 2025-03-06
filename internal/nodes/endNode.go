package nodes

import (
	"encoding/json"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"time"
)

type EndNodeExecutor struct{}

func (e *EndNodeExecutor) Execute(flowExecutionID, nodeID string, ch chan schema.TaskStep, result chan protocol.Message) {
	stat := protocol.Message{
		FlowExecutionID: flowExecutionID,
		Type:            protocol.MsgTaskResult,
		AgentID:         "",
		NodeID:          nodeID,
		Timestamp:       time.Now().Unix(),
		Payload:         schema.NodeStateSuccess,
	}

	result <- stat
}

func (e *EndNodeExecutor) NodeType() string {
	return "end"
}

func (e *EndNodeExecutor) AgentId() string {
	return "end"
}

func NewEndNodeExecutor(node schema.Node) *EndNodeExecutor {
	var data EndNodeExecutor
	// 从node.Properties中获取属性
	err := json.Unmarshal(node.Properties, &data)
	if err != nil {
		return nil
	}
	return &data
}
