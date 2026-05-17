package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type MailpitConfig struct {
	Host     string
	Port     string
	From     string
	Username string
	Password string
}

type MailpitSender struct {
	config MailpitConfig
}

func NewMailpitSender(cfg MailpitConfig) *MailpitSender {
	return &MailpitSender{config: cfg}
}

func (m *MailpitSender) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.config.Host, m.config.Port)
	log.Printf("📧 [MAILPIT] Sending to: %s, addr: %s", to, addr)

	headers := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		to, m.config.From, subject)

	msg := []byte(headers + body)

	// ✅ Mailpit doesn't need authentication - use nil auth
	err := smtp.SendMail(addr, nil, m.config.From, []string{to}, msg)
	if err != nil {
		log.Printf("❌ [MAILPIT] Send failed: %v", err)
		return err
	}

	log.Printf("✅ [MAILPIT] Email sent successfully to %s", to)
	return nil
}
