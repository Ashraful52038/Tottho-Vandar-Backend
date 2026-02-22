package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailService struct {
	from     string
	host     string
	port     string
	username string
	password string
}

func NewEmailService(host, port, username, password, from string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	body := fmt.Sprintf("Please verify your email using this link: http://localhost:3000/verify-email?token=%s", token)
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendResetPasswordEmail(to, token string) error {
	subject := "Password Reset"
	body := fmt.Sprintf("Reset your password using this link: http://localhost:3000/reset-password?token=%s", token)
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	// SMTP connection
	// auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	var auth smtp.Auth = nil

	log.Printf("📧 Sending email to: %s", to)
	log.Printf("📧 Subject: %s", subject)
	log.Printf("📧 Body: %s", body)

	if s.username != "" && s.password != "" && s.host != "localhost" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	// Email message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	// Send email
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent to %s", to)
	return nil
}
