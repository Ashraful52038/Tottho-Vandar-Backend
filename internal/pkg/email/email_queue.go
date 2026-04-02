package email

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// EmailMessage struct
type EmailMessage struct {
	Type      string `json:"type"`
	To        string `json:"to"`
	Username  string `json:"username"`
	Token     string `json:"token,omitempty"`
	Commenter string `json:"commenter,omitempty"`
	Content   string `json:"content,omitempty"`
	PostID    uint   `json:"postId,omitempty"`
}

// EmailQueueService struct
type EmailQueueService struct {
	emailService *EmailService
	rabbitConn   *amqp.Connection
	rabbitChan   *amqp.Channel
	queueName    string
}

// NewEmailQueueService - constructor
func NewEmailQueueService(emailService *EmailService, rabbitURL string) (*EmailQueueService, error) {
	// RabbitMQ connection
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare queue
	queue, err := ch.QueueDeclare(
		"email_notifications", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	return &EmailQueueService{
		emailService: emailService,
		rabbitConn:   conn,
		rabbitChan:   ch,
		queueName:    queue.Name,
	}, nil
}

// PublishEmail - publish message to queue
func (eq *EmailQueueService) PublishEmail(msg EmailMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = eq.rabbitChan.Publish(
		"",           // exchange
		eq.queueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return err
	}

	log.Printf("📨 Email message published to queue: %s", msg.Type)
	return nil
}

// ConsumeEmails - start consumer
func (eq *EmailQueueService) ConsumeEmails() error {
	msgs, err := eq.rabbitChan.Consume(
		eq.queueName, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	go func() {
		log.Println("📬 Email consumer started, waiting for messages...")
		for d := range msgs {
			var msg EmailMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				d.Nack(false, false) // reject, don't requeue
				continue
			}

			log.Printf("📧 Processing email: %s to %s", msg.Type, msg.To)

			// Process email based on type
			var emailErr error
			switch msg.Type {
			case "verification":
				emailErr = eq.emailService.SendVerificationEmail(msg.To, msg.Token)
			case "reset":
				emailErr = eq.emailService.SendResetPasswordEmail(msg.To, msg.Token)
			case "reply":
				emailErr = eq.sendReplyNotificationEmail(msg)
			default:
				log.Printf("Unknown email type: %s", msg.Type)
			}

			if emailErr != nil {
				d.Nack(false, true) // requeue
			} else {
				d.Ack(false) // acknowledge
			}
		}
	}()

	return nil
}

// sendReplyNotificationEmail - helper for reply notifications
func (eq *EmailQueueService) sendReplyNotificationEmail(msg EmailMessage) error {
	subject := fmt.Sprintf("%s replied to your comment", msg.Commenter)
	body := fmt.Sprintf(`
Hello %s,

%s replied to your comment:
"%s"

View the conversation: http://localhost:3000/posts/%d

---
You're receiving this because you have notifications enabled.
`, msg.Username, msg.Commenter, msg.Content, msg.PostID)

	return eq.emailService.sendEmail(msg.To, subject, body)
}

// Close - close connections
func (eq *EmailQueueService) Close() {
	if eq.rabbitChan != nil {
		eq.rabbitChan.Close()
	}
	if eq.rabbitConn != nil {
		eq.rabbitConn.Close()
	}
	log.Println("RabbitMQ connection closed")
}
