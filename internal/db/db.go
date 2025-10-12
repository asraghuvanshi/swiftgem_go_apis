package db

import (
	"fmt"
	"swiftgem_go_apis/internal/config"
	"swiftgem_go_apis/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.AppConfig.DBDSN), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	fmt.Println("Database connected")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Feed{},
		&models.Notification{},
		// &models.Chat{}, // Uncomment when chats module is implemented
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migration completed")
}
