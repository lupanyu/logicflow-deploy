package schema

type WebProperties struct {
	AppName       string `json:"appName,omitempty"`
	Host          string `json:"host,omitempty"`
	PackageSource string `json:"packageSource,omitempty"`
	DeployPath    string `json:"deployPath,omitempty"`
	BakPath       string `json:"bakPath,omitempty"`
	ServerName    string `json:"serverName,omitempty"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}
