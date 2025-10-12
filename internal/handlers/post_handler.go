// internal/handlers/post_handler.go
package handlers

import (
	"net/http"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/services"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id") // From JWT middleware
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}
	post.UserID = userID

	err := services.CreatePost(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Post created", Data: post})
}

func GetHomePosts(c *gin.Context) {
	userID := c.GetUint("user_id")
	posts, err := services.GetHomePosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Status: true, Message: "Posts retrieved", Data: posts})
}
