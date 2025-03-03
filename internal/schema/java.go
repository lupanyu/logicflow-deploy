package schema

type JavaProperties struct {
	AppName            string `json:"appName,omitempty"`
	Host               string `json:"host,omitempty"`
	PackageSource      string `json:"packageSource,omitempty"`
	DeployPath         string `json:"deployPath,omitempty"`
	BakPath            string `json:"bakPath,omitempty"`
	ServerName         string `json:"serverName,omitempty"`
	Port               int    `json:"port,omitempty"`
	HealthUri          string `json:"healthUri,omitempty"`
	HealthCheckTimeout int    `json:"healthCheckTimeout,omitempty"`
	URL                string `json:"url,omitempty"`
	Method             string `json:"method,omitempty"`
	Body               string `json:"body,omitempty"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	Name               string `json:"name,omitempty"`
	JobName            string `json:"jobName,omitempty"`
}
