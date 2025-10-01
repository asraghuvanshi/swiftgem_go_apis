package main

import (
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect() // Connect to database

	r := gin.Default()    // Initialize Gin
	routes.SetupRoutes(r) // Setup all routes

	r.Run(":8080")
}
