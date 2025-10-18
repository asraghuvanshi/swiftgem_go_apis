// internal/repositories/post_repo.go
package repositories

import (
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
)

func CreatePost(post *models.Post) error {
	return db.DB.Create(post).Error
}

func GetPosts() ([]models.Post, error) {
	var posts []models.Post
	err := db.DB.Order("created_at DESC").Find(&posts).Error
	return posts, err
}
