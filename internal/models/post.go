// internal/models/post.go
package models

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	Media     string    `json:"media"` // URL or path
	CreatedAt time.Time `json:"created_at"`
}
