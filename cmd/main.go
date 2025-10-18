package main

import (
	"os"
	"swiftgem_go_apis/internal/config"
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db.Connect()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.SetupRoutes(r)
	// Create uploads directory
	r.Static("/v1/uploads", "./uploads")
	os.MkdirAll(config.AppConfig.UploadDir, 0755)
	r.Run(":" + config.AppConfig.Port)
}
