package routes

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/api"
	"logicflow-deploy/internal/server"
	"logicflow-deploy/internal/services"
)

func RegisterAPIRoutes(r *gin.Engine, s *server.Server) {
	// 部署相关路由组
	deployGroup := r.Group("/api/v1/deploy")
	{
		deployGroup.POST("", func(context *gin.Context) {
			api.StartDeploy(context, s)
		}) // 开始部署
		deployGroup.GET("/:id", func(context *gin.Context) {
			api.GetStatus(context, s)
		}) // 获取部署的状态
	}

	r.GET("/ws", func(c *gin.Context) { services.HandleFlowExecution(s, c) })
}
