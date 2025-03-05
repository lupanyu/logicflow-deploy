package nodes

import (
	"context"
	"encoding/json"
	"logicflow-deploy/internal/schema"
	"time"
)

type StartNodeExecutor struct{}

func (e *StartNodeExecutor) Execute(ctx context.Context, state chan schema.TaskStep) {
	stat := schema.TaskStep{
		Status:  schema.TaskStateRunning,
		Setup:   "开始部署",
		AgentID: "",
		Output:  schema.NewOutLog(schema.LevelInfo, "开始应用部署"),
	}
	state <- stat
	time.Sleep(time.Millisecond)
	stat.Status = schema.TaskStateSuccess
	stat.Output = schema.NewOutLog(schema.LevelInfo, "应用部署成功")
	state <- stat
}

func (e *StartNodeExecutor) NodeType() string {
	return "start"
}
func (e *StartNodeExecutor) AgentId() string {
	return "start"
}

func NewStartNodeExecutor(node schema.Node) *StartNodeExecutor {
	var data StartNodeExecutor
	// 从node.Properties中获取属性
	err := json.Unmarshal(node.Properties, &data)
	if err != nil {
		return nil
	}
	return &data
}
