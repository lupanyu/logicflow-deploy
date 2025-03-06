package nodes

import (
	"encoding/json"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"os"
)

// 这是一个在server构建节点

type BuildNodeExecutor struct {
	Name    string `json:"name,omitempty"`
	Env     string `json:"env"`
	Cmd     string `json:"cmd,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
	workDir string
	out     *os.File
}

func (b *BuildNodeExecutor) Execute(flowExecutionID, nodeID string, ch chan schema.TaskStep, result chan protocol.Message) {
	ch <- schema.TaskStep{
		FlowExecutionID: flowExecutionID,
		Status:          schema.TaskStateSuccess,
		NodeID:          nodeID,
		Setup:           "开始构建",
		Output:          schema.NewOutLog(schema.LevelInfo, "开始构建"),
	}
}
func (b *BuildNodeExecutor) NodeType() string {
	return "build"
}

func (b *BuildNodeExecutor) AgentId() string {
	return "build"
}

func NewBuildNodeExecutor(node schema.Node) *BuildNodeExecutor {
	var data BuildNodeExecutor
	// 从node.Properties中获取属性
	err := json.Unmarshal(node.Properties, &data)
	if err != nil {
		return nil
	}
	return &data
}
