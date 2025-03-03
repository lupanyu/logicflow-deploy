package nodes

import (
	"context"
	"encoding/json"
	"logicflow-deploy/internal/schema"
)

type EndNodeExecutor struct{}

func (e *EndNodeExecutor) Execute(ctx context.Context, state chan schema.TaskStep) {
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
