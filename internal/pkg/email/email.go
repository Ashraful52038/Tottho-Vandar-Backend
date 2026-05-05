package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailService struct {
	from        string
	host        string
	port        string
	username    string
	password    string
	frontendURL string
}

func NewEmailService(host, port, username, password, from, frontendURL string) *EmailService {
	return &EmailService{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		from:        from,
		frontendURL: frontendURL,
	}
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	body := fmt.Sprintf("Please verify your email using this link: %s/verify-email?token=%s", s.frontendURL, token)
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendResetPasswordEmail(to, token string) error {
	subject := "Password Reset"
	body := fmt.Sprintf("Reset your password using this link: %s/reset-password?token=%s", s.frontendURL, token)
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendReplyNotification(to, username, commenter, content string, postID uint) error {
	subject := fmt.Sprintf("%s replied to your comment", commenter)
	body := fmt.Sprintf(`
Hello %s,

%s replied to your comment:
"%s"

View the conversation: %s/posts/%d

---
You're receiving this because you have notifications enabled.
`, username, commenter, content, s.frontendURL, postID)

	return s.sendEmail(to, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	var auth smtp.Auth = nil

	if s.username != "" && s.password != "" && s.host != "localhost" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	// Email message
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, to, subject, body))

	// Send email
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent to %s", to)
	return nil
}
