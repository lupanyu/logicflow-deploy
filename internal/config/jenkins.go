package config

import (
	"encoding/json"
	"io/ioutil"
	"logicflow-deploy/internal/schema"
	"os"
)

// 全局连接配置存储
var JenkinsConnections = make(map[string]schema.JenkinsConnection)

// GetJenkinsConnection 获取预配置的连接信息
func GetJenkinsConnection(name string) (schema.JenkinsConnection, bool) {
	conn, ok := JenkinsConnections[name]
	return conn, ok
}

// LoadJenkinsConnections 从配置文件加载连接配置（示例）
func LoadJenkinsConnections(configData []byte) error {
	var conf struct {
		Connections []schema.JenkinsConnection `json:"connections"`
	}

	if err := json.Unmarshal(configData, &conf); err != nil {
		return err
	}

	for _, conn := range conf.Connections {
		JenkinsConnections[conn.Name] = conn
	}
	return nil
}

func init() {
	// 从默认配置文件中加载连接配置
	f, err := os.OpenFile("configs/jenkins.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(f)
	//log.Printf("加载jenkins配置文件成功， 内容是：%s", string(data))
	err = LoadJenkinsConnections(data)
	if err != nil {
		panic(err)
	}
}
