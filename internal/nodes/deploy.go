package nodes

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
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

func CheckAPPHealth(status *schema.TaskStep, conn *websocket.Conn, port int, uri string, timeout time.Duration) ([]byte, error) {
	status.Setup = "健康检查"
	defer sendStatus(conn, *status)
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

func sendStatus(conn *websocket.Conn, status schema.TaskStep) {

	event, err := protocol.NewMessage(protocol.MsgTaskStep, status.FlowExecutionID, status.AgentID,
		status.NodeID, status)
	if err != nil {
		log.Printf("[%s]发送状态失败：%v", utils.GetCallerInfo(), err)
		return
	}
	log.Printf("[%s]发送部署任务的状态：%v", utils.GetCallerInfo(), event)
	_ = conn.WriteJSON(event)
}

func sendLastResult(conn *websocket.Conn, data protocol.Message) {
	data.Timestamp = time.Now().UnixNano()
	log.Printf("[%s]发送部署任务的最后结果：%v", utils.GetCallerInfo(), data)
	_ = conn.WriteJSON(data)
}

// 错误处理闭包
func handleStep(step *schema.TaskStep, stepName string, conn *websocket.Conn, fn func() ([]byte, error)) bool {
	step.Setup = stepName
	// 统一处理结果发送
	defer sendStatus(conn, *step)

	out, err := fn()
	step.Output = string(out)

	if err != nil {
		step.Status = schema.TaskStateFailed
		step.Error = err.Error()
		return false
	}
	return true
}

// 解压tar包
func UnpackTar(tarFile, dest string) ([]byte, error) {
	return utils.RunShell("tar -zxf " + tarFile + " -C " + dest)
}
