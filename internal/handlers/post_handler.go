// internal/handlers/post_handler.go
package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ---------- DTOs ----------
type AuthorDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PostResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	Country     string    `json:"country,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      AuthorDTO `json:"author"`
}

func postToDTO(p *models.Post) PostResponse {
	return PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		Country:     p.Country,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Author: AuthorDTO{
			ID:   p.User.ID,
			Name: p.User.Name,
		},
	}
}

// ---------- CREATE ----------
type PostRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Country     string `form:"country"`
}

func CreatePost(c *gin.Context) {
	// ---- JWT User ----
	user, _ := c.Get("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// ---- Bind form ----
	var req PostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error()})
		return
	}

	file, _ := c.FormFile("image")
	hasImage := file != nil
	if req.Title == "" && !hasImage {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: "title or image required"})
		return
	}

	post := models.Post{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Country:     req.Country,
	}

	if hasImage {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, Response{Status: false, Message: "only jpg/png allowed"})
			return
		}
		filename := time.Now().Format("20060102150405") + ext
		path := filepath.Join("uploads", filename)

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "save file: " + err.Error()})
			return
		}
		post.ImageURL = strings.ReplaceAll(path, "\\", "/")
		log.Printf("Image uploaded: %s", post.ImageURL)
	}

	if err := services.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "db error: " + err.Error()})
		return
	}

	// Re-load with User
	var full models.Post
	if err := db.DB.Preload("User").First(&full, post.ID).Error; err != nil {
		full = post
	} else {
		full.ImageURL = post.ImageURL
	}

	c.JSON(http.StatusOK, Response{
		Status:  true,
		Message: "Post created",
		Data:    postToDTO(&full),
	})
}

// ---------- LIST (with pagination) ----------
func GetHomePosts(c *gin.Context) {
	filter := models.PostFilter{
		Country:    c.Query("country"),
		TimeFilter: c.Query("time_filter"),
		Page:       1,
		PageSize:   10,
	}
	if p, _ := strconv.Atoi(c.Query("page")); p > 0 {
		filter.Page = p
	}
	if s, _ := strconv.Atoi(c.Query("page_size")); s > 0 && s <= 100 {
		filter.PageSize = s
	}
	if filter.TimeFilter == "custom" {
		if st := c.Query("start_time"); st != "" {
			if t, e := time.Parse(time.RFC3339, st); e == nil {
				filter.StartTime = t
			}
		}
		if et := c.Query("end_time"); et != "" {
			if t, e := time.Parse(time.RFC3339, et); e == nil {
				filter.EndTime = t
			}
		}
	}

	posts, total, err := services.GetPosts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: err.Error()})
		return
	}

	dtos := make([]PostResponse, len(posts))
	for i := range posts {
		dtos[i] = postToDTO(&posts[i])
	}

	resp := map[string]interface{}{
		"posts":       dtos,
		"total":       total,
		"page":        filter.Page,
		"page_size":   filter.PageSize,
		"total_pages": (total + int64(filter.PageSize) - 1) / int64(filter.PageSize),
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Posts fetched", Data: resp})
}

// ---------- EDIT ----------
func EditPost(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	currentUserID := uint(claims["user_id"].(float64))

	id, _ := strconv.Atoi(c.Param("id"))
	post, err := services.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Status: false, Message: "Post not found"})
		return
	}

	if post.UserID != currentUserID {
		c.JSON(http.StatusForbidden, Response{Status: false, Message: "You can only edit your own posts"})
		return
	}

	var req PostRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error()})
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Description != "" {
		post.Description = req.Description
	}
	if req.Country != "" {
		post.Country = req.Country
	}

	if file, err := c.FormFile("image"); err == nil && file != nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, Response{Status: false, Message: "only jpg/png"})
			return
		}
		filename := time.Now().Format("20060102150405") + ext
		path := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "save file: " + err.Error()})
			return
		}
		post.ImageURL = strings.ReplaceAll(path, "\\", "/")
	}

	if err := services.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "update failed: " + err.Error()})
		return
	}

	var full models.Post
	db.DB.Preload("User").First(&full, post.ID)
	full.ImageURL = post.ImageURL

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post updated", Data: postToDTO(&full)})
}

// ---------- DELETE ----------
func DeletePost(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	currentUserID := uint(claims["user_id"].(float64))

	id, _ := strconv.Atoi(c.Param("id"))
	post, err := services.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Status: false, Message: "Post not found"})
		return
	}

	if post.UserID != currentUserID {
		c.JSON(http.StatusForbidden, Response{Status: false, Message: "You can only delete your own posts"})
		return
	}

	if err := services.DeletePost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: "delete failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post deleted", Data: nil})
}
