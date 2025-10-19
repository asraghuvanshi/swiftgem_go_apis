// internal/handlers/post_handler.go
package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PostRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Country     string `form:"country"`
}

func CreatePost(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "User not authenticated", Data: nil})
		return
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Invalid token format", Data: nil})
		return
	}

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
		Country:     c.PostForm("country"), // New field from form
		CreatedAt:   time.Now(),
	}

	if hasImage {
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Only JPG/PNG images are allowed", Data: nil})
			return
		}

		imagePath := filepath.Join("uploads", time.Now().Format("20060102150405")+ext)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to save image: " + err.Error(), Data: nil})
			return
		}
		post.ImageURL = strings.ReplaceAll(imagePath, "\\", "/")
	}

	if err := services.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to create post: " + err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post created successfully", Data: post})
}

// GetHomePosts (updated with filters, pagination, no token required)
func GetHomePosts(c *gin.Context) {
	filter := services.PostFilter{
		Country:    c.Query("country"),
		TimeFilter: c.Query("time_filter"),
		Page:       1,
		PageSize:   10, // Default page size
	}

	// Parse pagination parameters
	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	// Parse custom time range
	if filter.TimeFilter == "custom" {
		startTimeStr := c.Query("start_time")
		endTimeStr := c.Query("end_time")
		if startTimeStr != "" && endTimeStr != "" {
			startTime, err1 := time.Parse(time.RFC3339, startTimeStr)
			endTime, err2 := time.Parse(time.RFC3339, endTimeStr)
			if err1 == nil && err2 == nil {
				filter.StartTime = startTime
				filter.EndTime = endTime
			}
		}
	}

	posts, total, err := services.GetPosts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to fetch posts: " + err.Error(), Data: nil})
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"posts":       posts,
		"total":       total,
		"page":        filter.Page,
		"page_size":   filter.PageSize,
		"total_pages": (total + int64(filter.PageSize) - 1) / int64(filter.PageSize),
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Posts fetched successfully", Data: response})
}

// EditPost allows the post creator or admin to edit a post
func EditPost(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "User not authenticated", Data: nil})
		return
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Invalid token format", Data: nil})
		return
	}

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
	userRole, _ := claims["role"].(string)

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Invalid post ID", Data: nil})
		return
	}

	post, err := services.GetPostByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Status: false, Message: "Post not found", Data: nil})
		return
	}

	// Check if user is the creator or an admin
	if post.UserID != uint(userID) && userRole != "admin" {
		c.JSON(http.StatusForbidden, Response{Status: false, Message: "You are not authorized to edit this post", Data: nil})
		return
	}

	var req PostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Invalid input: " + err.Error(), Data: nil})
		return
	}

	file, err := c.FormFile("image")
	hasImage := err == nil && file != nil

	// Update fields
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Description != "" {
		post.Description = req.Description
	}
	if country := c.PostForm("country"); country != "" {
		post.Country = country
	}

	if hasImage {
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Only JPG/PNG images are allowed", Data: nil})
			return
		}

		imagePath := filepath.Join("uploads", time.Now().Format("20060102150405")+ext)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to save image: " + err.Error(), Data: nil})
			return
		}
		post.ImageURL = strings.ReplaceAll(imagePath, "\\", "/")
	}

	if err := services.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to update post: " + err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post updated successfully", Data: post})
}

// DeletePost allows the post creator or admin to delete a post
func DeletePost(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "User not authenticated", Data: nil})
		return
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Invalid token format", Data: nil})
		return
	}

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
	userRole, _ := claims["role"].(string)

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: "Invalid post ID", Data: nil})
		return
	}

	post, err := services.GetPostByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Status: false, Message: "Post not found", Data: nil})
		return
	}

	// Check if user is the creator or an admin
	if post.UserID != uint(userID) && userRole != "admin" {
		c.JSON(http.StatusForbidden, Response{Status: false, Message: "You are not authorized to delete this post", Data: nil})
		return
	}

	if err := services.DeletePost(uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "Failed to delete post: " + err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post deleted successfully", Data: nil})
}
