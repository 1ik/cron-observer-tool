package gmail

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	To      []string // Recipient email addresses
	Subject string   // Email subject
	Body    string   // Email body (plain text or HTML)
}

// Sender defines the interface for sending emails
type Sender interface {
	// Send sends an email message
	// Returns an error if the email could not be sent
	Send(msg EmailMessage) error
}
