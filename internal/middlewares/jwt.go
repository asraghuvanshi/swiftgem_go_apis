package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"swiftgem_go_apis/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Authorization header required"})
			c.Abort()
			return
		}

		// Check for "Bearer " prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid authorization header format. Use Bearer <token>"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the token and user_id in the context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("user", token)                   // Set the token for compatibility with CreatePost
		c.Set("user_id", uint(userID))         // Set user_id for convenience
		c.Next()
	}
}