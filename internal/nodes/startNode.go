package nodes

import (
	"logicflow-deploy/internal/schema"
	"time"
)

type StartNodeExecutor struct {
}

func (e *StartNodeExecutor) Execute() schema.TaskStep {
	stat := schema.TaskStep{
		Status:  schema.TaskStateRunning,
		Setup:   "开始部署",
		AgentID: "",
		Output:  schema.NewOutLog(schema.LevelInfo, "开始应用部署"),
	}
	time.Sleep(time.Millisecond)
	stat.Status = schema.TaskStateSuccess
	stat.Output = schema.NewOutLog(schema.LevelInfo, "应用部署成功")
	return stat
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
