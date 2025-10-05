package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/services"
	"swiftgem_go_apis/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Signup(c *gin.Context) {
	var user models.User

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&user); err != nil {
		var messages []string

		if errs, ok := err.(validator.ValidationErrors); ok {
			t := reflect.TypeOf(user)
			for _, fieldErr := range errs {
				field, _ := t.FieldByName(fieldErr.StructField())
				jsonTag := field.Tag.Get("json")
				jsonField := strings.Split(jsonTag, ",")[0]
				if jsonField == "" {
					jsonField = strings.ToLower(fieldErr.Field())
				}

				var msg string
				switch fieldErr.Tag() {
				case "required":
					msg = fmt.Sprintf("%s is required", jsonField)
				case "email":
					msg = fmt.Sprintf("%s is not valid", jsonField)
				case "min":
					msg = fmt.Sprintf("%s must be at least %s characters", jsonField, fieldErr.Param())
				case "oneof":
					msg = fmt.Sprintf("%s must be one of %s", jsonField, fieldErr.Param())
				default:
					msg = fmt.Sprintf("%s is invalid", jsonField)
				}

				messages = append(messages, msg)
			}
		} else {
			messages = append(messages, err.Error())
		}

		helpers.Error(c, http.StatusBadRequest, "Validation failed", messages)
		return
	}

	if err := services.Signup(&user); err != nil {
		helpers.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response := gin.H{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"gender":      user.Gender,
		"phoneNumber": user.Phone,
		"created_at":  user.CreatedAt,
	}

	helpers.Success(c, "User created successfully", response)
}

func Login(c *gin.Context) {
	type request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		var messages []string
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				jsonField := e.Field()
				if f, _ := reflect.TypeOf(req).FieldByName(e.StructField()); f.Tag.Get("json") != "" {
					jsonField = strings.Split(f.Tag.Get("json"), ",")[0]
				}
				msg := ""
				switch e.Tag() {
				case "required":
					msg = jsonField + " is required"
				case "email":
					msg = jsonField + " is not valid"
				default:
					msg = jsonField + " is invalid"
				}
				messages = append(messages, msg)
			}
		} else {
			messages = append(messages, err.Error())
		}
		helpers.Error(c, http.StatusBadRequest, "Validation failed", messages)
		return
	}

	token, err := services.Login(req.Email, req.Password)
	if err != nil {
		helpers.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response := gin.H{
		"email": req.Email,
		"token": token,
	}

	helpers.Success(c, "Login successful", response)
}

var validate = validator.New()

func SendOTP(c *gin.Context) {
	type request struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phoneNumber" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		var messages []string
		if errs, ok := err.(validator.ValidationErrors); ok {
			t := reflect.TypeOf(req)
			for _, e := range errs {
				field, _ := t.FieldByName(e.StructField())
				jsonField := strings.Split(field.Tag.Get("json"), ",")[0]
				if jsonField == "" {
					jsonField = strings.ToLower(e.Field())
				}
				msg := ""
				switch e.Tag() {
				case "required":
					msg = jsonField + " is required"
				case "email":
					msg = jsonField + " is not valid"
				default:
					msg = jsonField + " is invalid"
				}
				messages = append(messages, msg)
			}
		} else {
			messages = append(messages, err.Error())
		}
		helpers.Error(c, http.StatusBadRequest, "Validation failed", messages)
		return
	}

	// Call service to send OTP
	if err := services.SendOTP(req.Email, req.Phone); err != nil {
		helpers.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Only return message, code, data null
	helpers.Success(c, "OTP sent successfully", nil)
}

func ResendOTP(c *gin.Context) {
	type request struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phoneNumber" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		var messages []string
		if errs, ok := err.(validator.ValidationErrors); ok {
			t := reflect.TypeOf(req)
			for _, e := range errs {
				field, _ := t.FieldByName(e.StructField())
				jsonField := strings.Split(field.Tag.Get("json"), ",")[0]
				if jsonField == "" {
					jsonField = strings.ToLower(e.Field())
				}
				msg := ""
				switch e.Tag() {
				case "required":
					msg = jsonField + " is required"
				case "email":
					msg = jsonField + " is not valid"
				default:
					msg = jsonField + " is invalid"
				}
				messages = append(messages, msg)
			}
		} else {
			messages = append(messages, err.Error())
		}
		helpers.Error(c, http.StatusBadRequest, "Validation failed", messages)
		return
	}

	// Call service to resend OTP
	if err := services.ResendOTP(req.Email, req.Phone); err != nil {
		helpers.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.Success(c, "OTP resent successfully", nil)
}

func VerifyOTP(c *gin.Context) {
	// Request body
	type request struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phoneNumber" binding:"required"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		var messages []string
		if errs, ok := err.(validator.ValidationErrors); ok {
			t := reflect.TypeOf(req)
			for _, e := range errs {
				field, _ := t.FieldByName(e.StructField())
				jsonField := strings.Split(field.Tag.Get("json"), ",")[0]
				if jsonField == "" {
					jsonField = strings.ToLower(e.Field())
				}

				msg := ""
				switch e.Tag() {
				case "required":
					msg = jsonField + " is required"
				case "email":
					msg = jsonField + " is not valid"
				case "len":
					msg = jsonField + " must be " + e.Param() + " digits"
				default:
					msg = jsonField + " is invalid"
				}
				messages = append(messages, msg)
			}
		} else {
			messages = append(messages, err.Error())
		}
		helpers.Error(c, http.StatusBadRequest, "Validation failed", messages)
		return
	}

	// Verify OTP using service
	user, err := services.VerifyOTP(req.Email, req.Phone, req.OTP)
	if err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Generate JWT token
	token, err := services.GenerateJWT(user.Email)
	if err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	type UserResponse struct {
		*models.User
		Token string `json:"token"`
	}

	response := UserResponse{
		User:  user,
		Token: token,
	}

	helpers.Success(c, "OTP verified successfully", response)
}
