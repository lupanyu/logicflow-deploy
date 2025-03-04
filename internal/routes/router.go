package routes

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/controller"
	"logicflow-deploy/internal/server"
)

func RegisterAPIRoutes(r *gin.Engine, s *server.Server) {
	// 部署相关路由组
	deployGroup := r.Group("/api/v1/deploy")
	{
		deployCtrl := controller.NewDeployController()
		deployGroup.POST("", func(context *gin.Context) {
			deployCtrl.StartDeploy(context, s)
		}) // 开始部署
		deployGroup.GET("/:id", deployCtrl.GetDeployStatus) // 获取部署状态
		deployGroup.DELETE("/:id", deployCtrl.CancelDeploy) // 取消当前部署
	}

	// r.GET("/ws", func(c *gin.Context) { services.HandleFlowExecution(s, c) })
}
