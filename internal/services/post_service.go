// internal/services/post_service.go
package services

import (
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/repositories"
)

func CreatePost(post *models.Post) error {
	return repositories.CreatePost(post)
}

func GetHomePosts(userID uint) ([]models.Post, error) {
	return repositories.GetPostsForUser(userID)
}
