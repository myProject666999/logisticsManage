package main

import (
	"log"
	"logisticsManage/config"
	"logisticsManage/middleware"
	"logisticsManage/models"
	"logisticsManage/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到.env文件，将使用系统环境变量")
	}

	config.LoadConfig()

	models.InitDB()
	models.AutoMigrate()
	models.SeedData()

	r := gin.Default()

	r.Use(middleware.CORS())

	routes.SetupRoutes(r)

	log.Printf("服务器启动在端口 %s", config.AppConfig.Server.Port)
	err = r.Run(":" + config.AppConfig.Server.Port)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
