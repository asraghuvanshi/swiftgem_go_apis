// internal/config/config.go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DBDSN            string
	JWTSecret        string
	JWTExpirationMin int
	MailHost         string
	MailPort         string
	MailUsername     string
	MailPassword     string
	MailFrom         string
	UploadDir        string // Add this for image storage path
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	AppConfig = Config{
		Port:             getEnv("PORT", "8080"),
		DBDSN:            getEnv("DB_DSN", "host=localhost user=postgres password=1234 dbname=swiftgem_go_apis port=5432 sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "swiftgem_go_apis_jwt_secret_key"),
		JWTExpirationMin: getEnvAsInt("JWT_EXPIRATION_MIN", 60),
		MailHost:         getEnv("MAIL_HOST", "sandbox.smtp.mailtrap.io"),
		MailPort:         getEnv("MAIL_PORT", "2525"),
		MailUsername:     getEnv("MAIL_USERNAME", ""),
		MailPassword:     getEnv("MAIL_PASSWORD", ""),
		MailFrom:         getEnv("MAIL_FROM", "swiftgem@csupport.test"),
		UploadDir:        getEnv("Upload_Dir","uploads"),

	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return value
}
