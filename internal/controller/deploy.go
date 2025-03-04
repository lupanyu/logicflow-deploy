package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/server"
)

type DeployController struct {
	deployService *server.FlowProcessor
}

func NewDeployController() *DeployController {
	return &DeployController{
		//deployService: server.NewDeployService(),
	}
}

func (c *DeployController) StartDeploy(ctx *gin.Context, s *server.Server) {
	var req schema.FlowData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}
	process, err := server.NewFlowProcessor(req, s)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	flowExecution := process.ExecuteFlow(context.Background(), s)
	ctx.JSON(200, flowExecution)
}

// GetDeployStatus 获取部署状态
// @Summary 获取部署任务状态
// @Produce json
// @Param   id path string true "部署ID"
// @Success 200 {object} services.DeployStatus
// @Router /api/v1/deployments/{id} [get]
func (c *DeployController) GetDeployStatus(ctx *gin.Context) {
	//id := ctx.Param("id")                   // 获取路径参数
	//status := c.deployService.GetStatus(id) // 调用服务层
	//ctx.JSON(200, gin.H{"status": status})  // 返回响应
}

// 新增取消接口
// @Summary 取消进行中的部署
// @Produce json
// @Param   id path string true "部署ID"
// @Success 200 {object} services.BaseResponse
// @Router /api/v1/deployments/{id} [delete]
func (c *DeployController) CancelDeploy(ctx *gin.Context) {
	//id := ctx.Param("id")
	//if err := c.deployService.CancelDeploy(id); err != nil {
	//	ctx.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}
	//ctx.JSON(200, gin.H{"message": "取消请求已接收"})
}
