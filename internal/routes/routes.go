package routes

import (
	"swiftgem_go_apis/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/sendOtp", handlers.SendOTP)
	r.POST("/resendOtp", handlers.ResendOTP)
	r.POST("/verifyOtp", handlers.VerifyOTP)

	// Protected routes (require JWT)
	// auth := r.Group("/")
	// auth.Use(middlewares.JWTMiddleware())
	// {
	// 	auth.GET("/profile", handlers.GetProfile)             // example protected route
	// 	auth.POST("/posts", handlers.CreatePost)              // create a post
	// 	auth.GET("/feed", handlers.GetFeed)                   // fetch feed
	// 	auth.GET("/notifications", handlers.GetNotifications) // notifications
	// }
}
