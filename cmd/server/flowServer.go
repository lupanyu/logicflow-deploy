package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"logicflow-deploy/internal/middleware"
	"logicflow-deploy/internal/routes"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/server"
	"logicflow-deploy/internal/utils"
	"os"
	"os/signal"
	"syscall"
)

func HandleFlowSave(c *gin.Context) {
	var flowData schema.Template
	if err := c.ShouldBindJSON(&flowData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 转换为原始JSON
	rawData, _ := json.Marshal(flowData)

	// 保存文件
	if err := server.SaveFlowModel(rawData); err != nil {
		c.JSON(500, gin.H{"error": "保存失败"})
		return
	}

	c.JSON(200, gin.H{"message": "保存成功"})
}

func main() {
	utils.InitLog()
	s := server.NewServer()
	r := gin.Default()
	r.Use(middleware.LogResponseBody())
	r.Use(middleware.CorsMiddleware())
	routes.RegisterAPIRoutes(r, s)
	s.SetHttp(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("正在保存流程执行数据...")
		// 实际保存逻辑需要根据存储类型实现
		if ms, ok := s.GetStorage().(*server.MemoryStorage); ok {
			if err := server.SaveMemStorageToFile(ms, "flow_storage.json"); err != nil {
				log.Printf("保存失败: %v", err)
			}
		}
		os.Exit(0)
	}()

	s.Start("0.0.0.0", 8080)
}
