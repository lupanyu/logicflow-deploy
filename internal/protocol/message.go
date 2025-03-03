package protocol

type MessageType int

const (
	MsgRegister         MessageType = iota + 1
	MsgRegisterResponse             // 注册响应
	MsgJavaDeploy                   // 部署java应用  server->agent
	MsgWebDeploy                    // 部署web应用   server->agent
	MsgTask                         // 部署任务控制  server->agent
	MsgStatus                       // 部署任务最终状态  agent->server 异常或者最终状态
	MsgHeartbeat                    // 心跳检测   	 server->agent  ,agent->server
	MsgTaskStep                     // 任务详细步骤以及日志信息  agent->server
	MsgTaskResult                   // 任务结果		agent->server
	MsgError                        // 错误信息		agent->server
)

// 部署消息  msgjavadeploy msgwebdeply msgtask
type Message struct {
	Type            MessageType `json:"type"`
	FlowExecutionID string      `json:"flowExecutionID"`   // flowExecution ID 实际执行flow时生成的id， 或者为空
	AgentID         string      `json:"agentID,omitempty"` // agentID, 一般是hostname
	NodeID          string      `json:"nodeID,omitempty"`  // flow中的nodeID
	Payload         interface{} `json:"payload,omitempty"`
	Timestamp       int64       `json:"timestamp"`
}
