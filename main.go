package main

import (
	"fmt"
	"log"

	"my-go-server/config"
	"my-go-server/database"
	"my-go-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 获取配置
	cfg := config.GetConfig()
	// 初始化数据库连接
	database.ConnectDB()
	// 创建默认的gin引擎
	r := gin.Default()
	// 注册路由
	routes.UserRoutes(r)
	// 启动服务器
	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器启动在 http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
