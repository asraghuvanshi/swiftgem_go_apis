package services

import (
	"errors"
	"math/rand"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("swiftgem_go_apis")

// ---------------- Password Hashing ---------------- //

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ---------------- Signup & Login ---------------- //

func Signup(user *models.User) error {
	hashed, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return repositories.CreateUser(user)
}

func Login(email, password string) (string, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}
	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return GenerateJWTByUser(user)
}

// ---------------- JWT Generation ---------------- //

func GenerateJWTByUser(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	return token.SignedString(JwtSecret)
}

func GenerateJWT(email string) (string, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	return GenerateJWTByUser(user)
}

// ---------------- OTP Functions ---------------- //

const otpLength = 6
const otpExpiry = 5 * time.Minute

func generateRandomOTP() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	otp := make([]byte, otpLength)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}

// Send OTP to user (email or phone)
func SendOTP(email, phone string) error {
	otp := generateRandomOTP()
	expiration := time.Now().Add(otpExpiry)

	// Save OTP in database
	if err := repositories.SaveOTP(email, phone, otp, expiration); err != nil {
		return err
	}

	// Send OTP via email/SMS (implement your own)
	// Example: sendEmail(email, otp) or sendSMS(phone, otp)
	return nil
}

// Resend OTP (can reuse SendOTP)
func ResendOTP(email, phone string) error {
	return SendOTP(email, phone)
}

// Verify OTP
func VerifyOTP(email, phone, otp string) (*models.User, error) {
	storedOTP, expiration, err := repositories.GetOTP(email, phone)
	if err != nil {
		return nil, errors.New("OTP not found")
	}

	if time.Now().After(expiration) {
		return nil, errors.New("OTP expired")
	}

	if storedOTP != otp {
		return nil, errors.New("Invalid OTP")
	}

	// Fetch user by email/phone
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	_ = repositories.DeleteOTP(email, phone)

	return user, nil
}
