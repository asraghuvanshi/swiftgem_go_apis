// internal/models/user.go
package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name" binding:"required"`
	Email      string    `gorm:"unique" json:"email" binding:"required,email"`
	Password   string    `json:"password" binding:"required,min=6"`
	Phone      string    `json:"phoneNumber" binding:"required"`
	Gender     string    `json:"gender" binding:"required,oneof=Male Female Other"`
	OTP        string    `json:"-"`
	OTPExpiry  time.Time `json:"-"`
	IsVerified bool      `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}
