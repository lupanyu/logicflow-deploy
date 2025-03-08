package schema

type JenkinsConnection struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JenkinsProperties struct {
	Name       string            `json:"name"`     // node节点名称
	NodeName   string            `json:"nodeName"` // jenkins节点的名称
	JobName    string            `json:"jobName"`
	Parameters map[string]string `json:"params,omitempty"`
}

var JenkinsNodes = make(map[string]JenkinsConnection)

// 从配置文件中读取jenkins节点的配置
func init() {
	// 从配置文件中读取jenkins节点的配置

}
