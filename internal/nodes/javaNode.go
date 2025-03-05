package nodes

import (
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
)

type JavaNodeExecutor struct {
	properties schema.JavaProperties
	agent      *protocol.AgentConnection
}

func (e *JavaNodeExecutor) AgentId() string {
	return e.properties.Host
}

func (e *JavaNodeExecutor) NodeType() string {
	return "java"
}

func (e *JavaNodeExecutor) Execute() schema.TaskStep {
	stat := schema.TaskStep{
		Status:  schema.TaskStateSuccess,
		Setup:   "发送部署指令",
		AgentID: e.properties.Host,
		Output:  schema.NewOutLog(schema.LevelInfo, "开始应用部署"),
	}

	// 执行部署命令
	err := e.deploy()
	if err != nil {
		stat.Status = schema.TaskStateFailed
		stat.Error = schema.NewOutLog(schema.LevelError, err.Error())
		log.Println("向%s发送部署指令异常，参数是：%s， 错误是: %v", e.properties.Host, e.properties, err.Error())
	} else {
		log.Println("向%s发送部署指令成功，参数是：%s", e.properties.Host, e.properties)
	}
	return stat
}

// 向agent发送部署命令
func (e *JavaNodeExecutor) deploy() error {
	return e.agent.Conn.WriteJSON(e.properties)
}

func NewJavaNodeExecutor(data schema.JavaProperties, agent *protocol.AgentConnection) *JavaNodeExecutor {
	return &JavaNodeExecutor{
		properties: data,
		agent:      agent,
	}
}
