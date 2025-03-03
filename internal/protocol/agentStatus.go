package protocol

// 定义 AgentStatus 枚举类型
type agentStatus string

const (
	AgentReady     agentStatus = "AgentReady"
	TaskInProgress agentStatus = "TaskInProgress"
	TaskTimeout    agentStatus = "TaskTimeout"
	TaskCompleted  agentStatus = "TaskCompleted"
)
