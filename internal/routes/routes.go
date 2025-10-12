// internal/routes/routes.go
package routes

import (
	"swiftgem_go_apis/internal/handlers"
	"swiftgem_go_apis/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/send-otp", handlers.SendOTP)
			auth.POST("/resend-otp", handlers.ResendOTP)
			auth.POST("/verify-otp", handlers.VerifyOTP)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middlewares.JWTAuth())
		{
			protected.POST("/posts", handlers.CreatePost)
			protected.GET("/home/posts", handlers.GetHomePosts)
			// Add more: user/profile, chats, etc.
		}
	}
}
