package controller

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/services"
)

type DeployController struct {
	deployService *services.DeployService
}

func NewDeployController() *DeployController {
	return &DeployController{
		deployService: services.NewDeployService(),
	}
}

func (c *DeployController) StartDeploy(ctx *gin.Context) {
	// 启动部署流程
}

// 具体控制器方法
func (c *DeployController) GetDeployStatus(ctx *gin.Context) {
	id := ctx.Param("id")                   // 获取路径参数
	status := c.deployService.GetStatus(id) // 调用服务层
	ctx.JSON(200, gin.H{"status": status})  // 返回响应
}
