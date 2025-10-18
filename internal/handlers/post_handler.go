package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PostRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func CreatePost(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "User not authenticated", Data: nil})
		return
	}

	// Type assertion for JWT token
	token, ok := user.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Invalid token format", Data: nil})
		return
	}

	// Extract user_id from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Invalid token claims", Data: nil})
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Invalid user ID in token", Data: nil})
		return
	}

	var req PostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Invalid input: " + err.Error(), Data: nil})
		return
	}

	file, err := c.FormFile("image")
	hasImage := err == nil && file != nil
	if req.Title == "" && !hasImage {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Title or image is required", Data: nil})
		return
	}

	post := models.Post{
		UserID:      uint(userID),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if hasImage {
		// Validate image type
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Only JPG/PNG images are allowed", Data: nil})
			return
		}

		// Save image to a directory
		imagePath := filepath.Join("uploads", time.Now().Format("20060102150405")+ext)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to save image: " + err.Error(), Data: nil})
			return
		}
		post.ImageURL = strings.ReplaceAll(imagePath, "\\", "/")
	}

	// Save post to database
	if err := repositories.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to create post: " + err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post created successfully", Data: post})
}

func GetHomePosts(c *gin.Context) {
	posts, err := repositories.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to fetch posts: " + err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Posts fetched successfully", Data: posts})
}
