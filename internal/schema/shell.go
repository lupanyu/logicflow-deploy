package schema

type ShellProperties struct {
	AppName             string `json:"appName,omitempty"`
	Host                string `json:"host,omitempty"`
	PreScriptContent    string `json:"preScriptContent,omitempty"`    // 前置脚本内容
	DeployScriptContent string `json:"deployScriptContent,omitempty"` // 部署脚本内容
	PostScriptContent   string `json:"postScriptContent,omitempty"`   // 后置脚本内容
	Timeout             int    `json:"timeout,omitempty"`             // 超时时间
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	Name                string `json:"name,omitempty"`
	JobName             string `json:"jobName,omitempty"`
}
