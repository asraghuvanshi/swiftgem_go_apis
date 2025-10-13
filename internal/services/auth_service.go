package services

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
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

	addr := fmt.Sprintf("%s:%s", config.AppConfig.MailHost, config.AppConfig.MailPort)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to mail server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.AppConfig.MailHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.AppConfig.MailHost,
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	auth := smtp.PlainAuth("", config.AppConfig.MailUsername, config.AppConfig.MailPassword, config.AppConfig.MailHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT TO failed: %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send DATA: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	client.Quit()

	log.Printf("OTP sent successfully to %s", email)
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

	log.Printf("Verify attempt - Input OTP: '%s' (len: %d)", otp, len(otp))
	log.Printf("Stored OTP: '%s' (len: %d)", user.OTP, len(user.OTP))
	log.Printf("Expiry: %v, Now: %v, Is Expired: %t", user.OTPExpiry, time.Now(), time.Now().After(user.OTPExpiry))

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
