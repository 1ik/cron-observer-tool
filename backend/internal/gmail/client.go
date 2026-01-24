package gmail

import (
	"fmt"
	"net/smtp"
	"strings"

	appconfig "github.com/yourusername/cron-observer/backend/internal/config"
)

// Client implements the Sender interface for sending emails via Gmail SMTP
type Client struct {
	config *appconfig.GmailConfig
}

// NewClient creates a new Gmail client with the provided configuration
func NewClient(config *appconfig.GmailConfig) *Client {
	return &Client{
		config: config,
	}
}

// Send sends an email message via Gmail SMTP
func (c *Client) Send(msg EmailMessage) error {
	// Validate configuration
	if c.config.User == "" {
		return fmt.Errorf("gmail user is not configured")
	}
	if c.config.Password == "" {
		return fmt.Errorf("gmail app password is not configured")
	}

	// Validate message
	if len(msg.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	if msg.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if msg.Body == "" {
		return fmt.Errorf("email body is required")
	}

	// Gmail SMTP configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Create authentication
	auth := smtp.PlainAuth("", c.config.User, c.config.Password, smtpHost)

	// Build email message
	to := strings.Join(msg.To, ", ")
	message := []byte(fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", msg.Subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		msg.Body + "\r\n")

	// Send email
	err := smtp.SendMail(smtpAddr, auth, c.config.User, msg.To, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
