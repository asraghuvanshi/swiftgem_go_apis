// internal/models/feed.go
package models

// Feed is an aggregated view, perhaps not stored, but for simplicity, assume it's a view or temp
type Feed struct {
	PostID uint `json:"post_id"`
	// Add more fields as needed
}
