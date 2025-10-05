package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, status string, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  status,
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	SendResponse(c, "success", http.StatusOK, message, data)
}

func Error(c *gin.Context, code int, message string, data interface{}) {
	SendResponse(c, "error", code, message, data)
}
