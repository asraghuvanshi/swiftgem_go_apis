// internal/handlers/auth_handler.go
package handlers

import (
	"net/http"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/services"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	err := services.Signup(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Signup successful, please request OTP", Data: nil})
}

func SendOTP(c *gin.Context) {
	type Req struct {
		Email string `json:"email" binding:"required,email"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	err := services.SendOTPService(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "OTP sent to email", Data: nil})
}

func ResendOTP(c *gin.Context) {
	type Req struct {
		Email string `json:"email" binding:"required,email"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	err := services.ResendOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "OTP resent to email", Data: nil})
}

func VerifyOTP(c *gin.Context) {
	type Req struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	err := services.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "OTP verified", Data: nil})
}

func Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	token, err := services.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Login successful", Data: map[string]string{"token": token}})
}
