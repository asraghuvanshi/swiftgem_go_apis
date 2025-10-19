// internal/services/post_service.go
package services

import (
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
	"time"
)

type PostFilter struct {
	Country    string
	TimeFilter string // "24h", "3d", "7d", "custom"
	StartTime  time.Time
	EndTime    time.Time
	Page       int
	PageSize   int
}

// CreatePost creates a new post
func CreatePost(post *models.Post) error {
	return db.DB.Create(post).Error
}

// GetPosts fetches posts with filters and pagination
func GetPosts(filter PostFilter) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := db.DB.Model(&models.Post{}).Preload("User").Order("created_at DESC")

	// Apply country filter
	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}

	// Apply time filter
	if filter.TimeFilter != "" {
		switch filter.TimeFilter {
		case "24h":
			query = query.Where("created_at >= ?", time.Now().Add(-24*time.Hour))
		case "3d":
			query = query.Where("created_at >= ?", time.Now().Add(-3*24*time.Hour))
		case "7d":
			query = query.Where("created_at >= ?", time.Now().Add(-7*24*time.Hour))
		case "custom":
			if !filter.StartTime.IsZero() && !filter.EndTime.IsZero() {
				query = query.Where("created_at BETWEEN ? AND ?", filter.StartTime, filter.EndTime)
			}
		}
	}

	// Count total posts for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// Fetch posts
	err := query.Find(&posts).Error
	return posts, total, err
}

// GetPostByID fetches a single post by ID
func GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := db.DB.Preload("User").First(&post, id).Error
	return &post, err
}

// UpdatePost updates an existing post
func UpdatePost(post *models.Post) error {
	post.UpdatedAt = time.Now()
	return db.DB.Save(post).Error
}

// DeletePost deletes a post by ID
func DeletePost(id uint) error {
	return db.DB.Delete(&models.Post{}, id).Error
}
