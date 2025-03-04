package nodes

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"net/http"
	"time"
)

func StopService(service string) ([]byte, error) {
	return utils.RunShell("systemctl stop " + service)
}

// 正确写法应该像其他方法一样有空格：
func StartService(service string) ([]byte, error) {
	return utils.RunShell("systemctl start " + service)
}

func BakOld(old, new string) ([]byte, error) {

	return utils.RunShell("cp -r " + old + " " + new)
}

func UpdateFile(downloadURL, new string) ([]byte, error) {
	return utils.RunShell("curl -o " + new + " " + downloadURL)
}

func Rollback(backup, new, service string) ([]byte, error) {
	return utils.RunShell("cp -r " + backup + " " + new)
}

func CheckAPPHealth(agentID string, conn *websocket.Conn, nodeId string, port int, uri string, timeout time.Duration) ([]byte, error) {
	status := schema.NewTaskStep(agentID, nodeId, "健康检查", schema.TaskStateSuccess, "", "")
	defer sendStatus(conn, agentID, *status)
	client := http.Client{Timeout: 3 * time.Second}
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%d%s", port, uri))
		if err == nil && resp.StatusCode == 200 {
			status.Status = schema.TaskStateSuccess
			return nil, nil
		}
		time.Sleep(5 * time.Second)
	}
	status.Status = schema.TaskStateFailed
	status.Error = "健康检查超时"
	return nil, errors.New("健康检查超时")
}

func sendStatus(conn *websocket.Conn, agentID string, status schema.TaskStep) {
	event := protocol.Message{
		Type:            protocol.MsgTaskStep,
		FlowExecutionID: status.FlowExecutionID,
		AgentID:         agentID,
		Timestamp:       time.Now().UnixNano(),
		Payload:         status,
	}
	_ = conn.WriteJSON(event)
}

func sendLastResult(conn *websocket.Conn, data protocol.Message) {
	data.Timestamp = time.Now().UnixNano()
	_ = conn.WriteJSON(data)
}

// 错误处理闭包
func handleStep(agentID, nodeId, stepName string, fn func() ([]byte, error)) bool {
	status := schema.NewTaskStep(agentID, nodeId, stepName, schema.TaskStateSuccess, "", "")
	out, err := fn()
	if err != nil {
		status.Status = schema.TaskStateFailed
		status.Output = string(out)
		status.Error = err.Error()
		//a.sendStatus(*status)
		return false
	}
	status.Output = string(out)
	//a.sendStatus(*status)
	return true
}

// 解压tar包
func UnpackTar(tarFile, dest string) ([]byte, error) {
	return utils.RunShell("tar -zxf " + tarFile + " -C " + dest)
}
