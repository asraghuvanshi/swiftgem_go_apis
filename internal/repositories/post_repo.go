// internal/repositories/post_repo.go
package repositories

import (
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
)

func CreatePost(post *models.Post) error {
	return db.DB.Create(post).Error
}

func GetPostsForUser(userID uint) ([]models.Post, error) {
	var posts []models.Post
	// Example: get all posts, in real: filter by follows etc.
	err := db.DB.Find(&posts).Error
	return posts, err
}
