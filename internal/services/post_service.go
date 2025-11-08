// internal/services/post_service.go
package services

import (
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/repositories"
)

func CreatePost(post *models.Post) error {
	return repositories.CreatePost(post)
}

func GetPosts(filter models.PostFilter) ([]models.Post, int64, error) {
	return repositories.GetPosts(filter) // now same type
}

func GetPostByID(id uint) (*models.Post, error) {
	return repositories.GetPostByID(id)
}

func UpdatePost(post *models.Post) error {
	return repositories.UpdatePost(post)
}

func DeletePost(id uint) error {
	return repositories.DeletePost(id)
}
