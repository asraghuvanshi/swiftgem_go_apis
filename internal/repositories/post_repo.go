// internal/repositories/post_repo.go
package repositories

import (
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
	"time"
)

func CreatePost(post *models.Post) error {
	return db.DB.Create(post).Error
}

func GetPosts(filter models.PostFilter) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := db.DB.Model(&models.Post{}).Preload("User")

	// Filters
	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}
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

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	query = query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC")

	err := query.Find(&posts).Error
	return posts, total, err
}

func GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := db.DB.Preload("User").First(&post, id).Error
	return &post, err
}

func UpdatePost(post *models.Post) error {
	post.UpdatedAt = time.Now()
	return db.DB.Save(post).Error
}

func DeletePost(id uint) error {
	return db.DB.Delete(&models.Post{}, id).Error
}
