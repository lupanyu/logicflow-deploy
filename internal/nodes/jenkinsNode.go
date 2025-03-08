package nodes

import (
	"context"
	"github.com/bndr/gojenkins"
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"time"
)

// JenkinsNode 定义Jenkins流水线节点结构
type JenkinsNodeExecutor struct {
	properties schema.JenkinsProperties
	agent      schema.JenkinsConnection
	ctx        context.Context
	jenkins    *gojenkins.Jenkins
}

// Execute 触发Jenkins任务并等待完成
func (j *JenkinsNodeExecutor) Execute(flowExecutionID, nodeID string, ch chan schema.TaskStep, result chan protocol.Message) {

	stat := schema.TaskStep{
		FlowExecutionID: flowExecutionID,
		NodeID:          nodeID,
		Status:          schema.TaskStateSuccess,
		Setup:           "发送部署指令",
		AgentID:         j.properties.NodeName,
		Output:          schema.NewOutLog(schema.LevelInfo, "发送构建指令成功", "..."),
	}
	builder, err := j.build()
	if err != nil {
		stat.Status = schema.TaskStateFailed
		stat.Output = schema.NewOutLog(schema.LevelError, "发送部署指令出现异常", "")
		stat.Error = schema.NewOutLog(schema.LevelError, "发送部署指令失败", err.Error())
		ch <- schema.TaskStep{}
	}
	ch <- stat
	// 轮询构建结果
	go j.waitForBuildCompletion(builder, stat, ch, result)
}

// waitForBuildCompletion 轮询构建状态
func (j *JenkinsNodeExecutor) waitForBuildCompletion(builder *gojenkins.Build, taskStep schema.TaskStep, ch chan schema.TaskStep, result chan protocol.Message) {
	msg, _ := protocol.NewMessage(protocol.MsgTaskResult, taskStep.FlowExecutionID, taskStep.AgentID, taskStep.NodeID, schema.NodeStateSuccess)

	for {
		time.Sleep(5 * time.Second)
		stat := j.getStatus(builder)
		if stat == "" {

		} else {
			if stat == "ABORTED" {
				taskStep.Status = schema.TaskStateFailed
				taskStep.Output = schema.NewOutLog(schema.LevelInfo, "构建任务", "构建任务被终止")
				taskStep.Error = schema.NewOutLog(schema.LevelError, "构建任务", "构建任务被终止")
			}
			if stat == "SUCCESS" {
				taskStep.Status = schema.TaskStateSuccess
				taskStep.Output = schema.NewOutLog(schema.LevelInfo, "构建任务", "构建任务完成")
				taskStep.Error = ""
			}
			ch <- taskStep
			time.Sleep(100 * time.Millisecond)
			if taskStep.Status == schema.TaskStateFailed {
				msg.UpdatePayload(schema.NodeStateFailed)
			}
			result <- msg
			break
		}
	}
}

func (j *JenkinsNodeExecutor) getStatus(builder *gojenkins.Build) string {
	poll, err := builder.Poll(j.ctx)
	if err != nil {
		return ""
	}
	log.Printf("[%s] 获取状态成功:响应码 %d 状态%s", utils.GetCallerInfo(), poll, builder.GetResult())
	return builder.GetResult()
}

func NewJenkinsNodeExecutor(data schema.JenkinsProperties, agent schema.JenkinsConnection) *JenkinsNodeExecutor {
	jenkins := gojenkins.CreateJenkins(nil, agent.URL, agent.Username, agent.Password)
	ctx := context.Background()
	_, err := jenkins.Init(ctx)

	if err != nil {
		panic(err)
	}

	return &JenkinsNodeExecutor{
		properties: data,
		agent:      agent,
		ctx:        ctx,
		jenkins:    jenkins,
	}

}

func (j *JenkinsNodeExecutor) build() (*gojenkins.Build, error) {

	queueid, err := j.jenkins.BuildJob(j.ctx, j.properties.JobName, nil)
	if err != nil {
		return nil, err
	}
	result, err := j.jenkins.GetBuildFromQueueID(j.ctx, queueid)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] 构建任务触发成功,revsion:%s buildNumber:%d", utils.GetCallerInfo(), result.GetRevision(), result.GetBuildNumber())
	return result, nil
}
