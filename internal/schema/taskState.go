package schema

import (
	"fmt"
	"time"
)

type NodeStatus string

// flowData的node节点的状态，server端保存的总状态,当node节点每个子步骤结束后把这个状态发送给server
const (
	NodeStatePending    NodeStatus = "pending"
	NodeStateRunning    NodeStatus = "running"
	NodeStateSuccess    NodeStatus = "success"
	NodeStateFailed     NodeStatus = "failed"
	NodeStateSkipped    NodeStatus = "skipped"
	NodeStateError      NodeStatus = "error"
	NodeStateTimeout    NodeStatus = "timeout"    // 超时状态
	NodeStateRollbacked NodeStatus = "rollbacked" // 已回滚状态
)

// 节点的子状态，agent端保存的实时状态
type TaskStatus string

const (
	TaskStateRunning TaskStatus = "running"
	TaskStateSuccess TaskStatus = "success"
	TaskStateFailed  TaskStatus = "failed"
)

// 任务的执行步骤详情，实时汇报给server
type TaskStep struct {
	FlowExecutionID string     `json:"flowExecutionId"` // flowExecution的id
	AgentID         string     `json:"agentId"`         // agent的id
	NodeID          string     `json:"nodeId"`          // flowData中的nodeId
	Setup           string     `json:"setup"`
	Status          TaskStatus `json:"status"`
	Output          string     `json:"output"`
	Error           string     `json:"error"`
}

func NewTaskStep(flowExecutionId, agentId, nodeId, setup string, status TaskStatus, output string, err string) *TaskStep {
	return &TaskStep{
		FlowExecutionID: flowExecutionId,
		AgentID:         agentId,
		NodeID:          nodeId,
		Setup:           setup,
		Status:          status,
		Output:          output,
		Error:           err,
	}
}

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

func NewOutLog(level Level, message string) string {
	return fmt.Sprintf("[%s] %s %s",
		time.Now().UTC().Format(time.DateTime),
		level,
		message)
}
