package nodes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"net/http"
	"os"
	"os/exec"
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
	defer sendStatus(conn, status)
	client := http.Client{Timeout: 3 * time.Second}
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%d%s", port, uri))
		if err == nil && resp.StatusCode == 200 {
			status.Status = schema.TaskStateSuccess
			return []byte("健康检查成功，响应码200"), nil
		}
		time.Sleep(5 * time.Second)
	}
	status.Status = schema.TaskStateFailed
	status.Error = "健康检查超时"
	return nil, errors.New("健康检查超时")
}

func sendStatus(conn *websocket.Conn, status *schema.TaskStep) {
	event, err := protocol.NewMessage(protocol.MsgTaskStep, status.FlowExecutionID, status.AgentID,
		status.NodeID, status)
	if err != nil {
		log.Printf("[%s]发送状态失败：%v", utils.GetCallerInfo(), err)
		return
	}
	log.Printf("[%s]发送部署任务的状态，FlowID：%s,nodeID：%s,status:%s", utils.GetCallerInfo(), event.FlowExecutionID, event.NodeID, status.Status)
	_ = conn.WriteJSON(event)
}

func sendLastResult(conn *websocket.Conn, data protocol.Message) {
	data.Timestamp = time.Now().UnixNano()
	log.Printf("[%s]发送部署任务的最后结果：%s", utils.GetCallerInfo(), data.Payload)
	_ = conn.WriteJSON(data)
}

// 错误处理闭包
func handleStep(step *schema.TaskStep, stepName string, conn *websocket.Conn, fn func() ([]byte, error)) bool {
	step.Setup = stepName
	// 统一处理结果发送
	defer sendStatus(conn, step)

	out, err := fn()
	step.Output = schema.NewOutLog(schema.LevelInfo, step.Setup, string(out))
	if err != nil {
		step.Status = schema.TaskStateFailed
		step.Error = err.Error()
		log.Printf("[%s]执行步骤[%s]失败：%v %v", utils.GetCallerInfo(), stepName, step.Status, err)
		return false
	}
	return true
}

// 解压tar包
func UnpackTar(tarFile, dest string) ([]byte, error) {
	return utils.RunShell("tar -zxf " + tarFile + " -C " + dest)
}

func handleShellDeploy(step *schema.TaskStep, stepName, shell string, timeout int, conn *websocket.Conn) bool {
	step.Setup = stepName
	defer sendStatus(conn, step)
	// 执行shell脚本
	output, err := executeShellScript(shell, time.Duration(timeout)*time.Second)
	step.Output = schema.NewOutLog(schema.LevelInfo, step.Setup, string(output))
	if err != nil {
		step.Status = schema.TaskStateFailed
		step.Error = err.Error()
		log.Printf("[%s]执行步骤[%s]失败：%v %v", utils.GetCallerInfo(), stepName, step.Status, err)
		return false
	}
	step.Status = schema.TaskStateSuccess // 设置状态为成功
	return true
}
func executeShellScript(scriptContent string, timeout time.Duration) ([]byte, error) {
	// 创建临时脚本文件
	tmpFile, err := os.CreateTemp("tmp/", "script-*.sh")
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入脚本内容
	if _, err := io.WriteString(tmpFile, "#!/bin/sh\n"+scriptContent); err != nil {
		return nil, fmt.Errorf("写入脚本内容失败: %v", err)
	}
	if err := tmpFile.Chmod(0700); err != nil {
		return nil, fmt.Errorf("设置执行权限失败: %v", err)
	}
	_ = tmpFile.Close()

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 执行脚本并捕获输出
	cmd := exec.CommandContext(ctx, "sh ", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	// 执行结果检查
	if ctx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("脚本执行超时: %v", timeout)
	}
	if err != nil {
		return output, fmt.Errorf("脚本执行失败: %v", err)
	}
	return output, nil
}
