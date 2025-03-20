package nodes

import (
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
)

type ShellNodeExecutor struct {
	properties schema.ShellProperties
	agent      *protocol.AgentConnection
}

func (e *ShellNodeExecutor) Execute(flowExecutionID, nodeID string, ch chan schema.TaskStep, result chan protocol.Message) {
	stat := schema.TaskStep{
		FlowExecutionID: flowExecutionID,
		NodeID:          nodeID,
		Status:          schema.TaskStateSuccess,
		Setup:           "发送部署指令",
		AgentID:         e.properties.Host,
		Output:          schema.NewOutLog(schema.LevelInfo, "发送部署指令成功", "..."),
	}

	// 执行部署命令
	err := e.deploy(flowExecutionID, nodeID)
	if err != nil {
		stat.Status = schema.TaskStateFailed
		stat.Error = schema.NewOutLog(schema.LevelError, "发送部署指令", err.Error())
		log.Printf("[%s] 向%s发送部署指令异常， 错误是: %v", utils.GetCallerInfo(), e.properties.Host, err.Error())
	} else {
		log.Printf("[%s] 向%s发送部署指令成功", utils.GetCallerInfo(), e.properties.Host)
	}
	ch <- stat
}

// 向agent发送部署命令
func (e *ShellNodeExecutor) deploy(flowExecutionID, nodeID string) error {

	data, _ := protocol.NewMessage(protocol.MsgShellDeploy, flowExecutionID, e.properties.Host, nodeID, e.properties)
	log.Printf("[%s] 向%s发送部署指令 参数是：%v", utils.GetCallerInfo(), e.properties.Host, data)

	return e.agent.WriteJSON(data)
}

func NewShellNodeExecutor(data schema.ShellProperties, agent *protocol.AgentConnection) *ShellNodeExecutor {
	return &ShellNodeExecutor{
		properties: data,
		agent:      agent,
	}
}
