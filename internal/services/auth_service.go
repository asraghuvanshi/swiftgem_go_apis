package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"swiftgem_go_apis/internal/config"
	"swiftgem_go_apis/internal/models"
	"swiftgem_go_apis/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func SendOTP(email, otp string) error {
	from := config.AppConfig.MailFrom
	to := []string{email}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Your OTP Code\r\n"+
		"\r\n"+
		"Your OTP is: %s\r\n", email, otp))

	auth := smtp.PlainAuth("", config.AppConfig.MailUsername, config.AppConfig.MailPassword, config.AppConfig.MailHost)
	addr := fmt.Sprintf("%s:%s", config.AppConfig.MailHost, config.AppConfig.MailPort)
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		log.Printf("Failed to send OTP to %s: %v", email, err)
		return fmt.Errorf("failed to send OTP: %w", err)
	}
	log.Printf("OTP sent to %s", email)
	return nil
}

func Signup(user *models.User) error {
	_, err := repositories.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPass)
	user.IsVerified = false
	// OTP and OTPExpiry are not set here; will be set in SendOTPService

	err = repositories.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func SendOTPService(email string) error {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		return errors.New("user already verified")
	}

	otp, err := GenerateOTP()
	if err != nil {
		return err
	}
	user.OTP = otp
	user.OTPExpiry = time.Now().Add(10 * time.Minute)

	err = repositories.UpdateUser(user)
	if err != nil {
		return err
	}

	return SendOTP(user.Email, otp)
}

func ResendOTP(email string) error {
	// Same logic as SendOTPService; kept separate for clarity
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		return errors.New("user already verified")
	}

	otp, err := GenerateOTP()
	if err != nil {
		return err
	}
	user.OTP = otp
	user.OTPExpiry = time.Now().Add(10 * time.Minute)

	err = repositories.UpdateUser(user)
	if err != nil {
		return err
	}

	return SendOTP(user.Email, otp)
}

func VerifyOTP(email, otp string) error {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.OTP != otp || time.Now().After(user.OTPExpiry) {
		return errors.New("invalid or expired OTP")
	}

	user.IsVerified = true
	user.OTP = ""
	user.OTPExpiry = time.Time{}
	return repositories.UpdateUser(user)
}

func Login(email, password string) (string, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsVerified {
		return "", errors.New("email not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * time.Duration(config.AppConfig.JWTExpirationMin)).Unix(),
	})

	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
