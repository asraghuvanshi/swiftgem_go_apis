// internal/models/post.go
package models

import "time"

type PostFilter struct {
	Country    string
	TimeFilter string // "24h", "3d", "7d", "custom"
	StartTime  time.Time
	EndTime    time.Time
	Page       int
	PageSize   int
}

type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Title       string    `gorm:"type:text" json:"title"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	ImageURL    string    `gorm:"column:image_url;type:varchar(255)" json:"image_url,omitempty"`
	Country     string    `gorm:"type:varchar(100)" json:"country,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
