package protocol

// 定义 AgentStatus 枚举类型
type AgentStatus string

const (
	AgentIdle       AgentStatus = "Idle"
	AgentInProgress AgentStatus = "TaskInProgress"
	AgentOffline    AgentStatus = "Offline"
)

func (a AgentStatus) String() string {
	return string(a)
}
