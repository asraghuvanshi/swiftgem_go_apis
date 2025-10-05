package repositories

import (
	"errors"
	"swiftgem_go_apis/internal/db"
	"swiftgem_go_apis/internal/models"
	"time"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Save or update OTP
func SaveOTP(email, phone, code string, expiresAt time.Time) error {
	var otp models.OTP
	err := db.DB.Where("email = ? AND phone = ?", email, phone).First(&otp).Error
	if err == nil {
		// Update existing OTP
		otp.Code = code
		otp.ExpiresAt = expiresAt
		return db.DB.Save(&otp).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new OTP
		newOTP := models.OTP{
			Email:     email,
			Phone:     phone,
			Code:      code,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
		}
		return db.DB.Create(&newOTP).Error
	} else {
		return err
	}
}

// Get OTP and expiration
func GetOTP(email, phone string) (string, time.Time, error) {
	var otp models.OTP
	err := db.DB.Where("email = ? AND phone = ?", email, phone).First(&otp).Error
	if err != nil {
		return "", time.Time{}, err
	}
	return otp.Code, otp.ExpiresAt, nil
}

// Delete OTP after verification
func DeleteOTP(email, phone string) error {
	return db.DB.Where("email = ? AND phone = ?", email, phone).Delete(&models.OTP{}).Error
}
