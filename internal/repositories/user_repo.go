package repositories

import (
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
)

func CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
