package main

import (
	"log"

	"hospital-api/src/models"
	"hospital-api/src/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	models.InitDB()

	r := gin.Default()

	// 注册各种业务路由
	routers.SetupRouter(r)

	// 全局中间件
	// r.Use(Cors())

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 路由挂载位置
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	log.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
