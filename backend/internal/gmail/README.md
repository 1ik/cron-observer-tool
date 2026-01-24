# Gmail Package

This package provides a Gmail email sending service with an interface-based design for easy dependency injection.

## Configuration

Add the following environment variables to your `.env` file:

```env
GMAIL_USER=your-email@gmail.com
GMAIL_APP_PASSWORD=your-app-password
```

## Usage

### 1. Get the Gmail Sender from main.go

The `gmailSender` is initialized in `main.go` and available for injection:

```go
var gmailSender gmail.Sender // Available in main.go after config loading
```

### 2. Inject into Handlers or Services

Pass the `gmailSender` to any handler or service that needs to send emails:

```go
// Example: In a handler
type MyHandler struct {
    gmailSender gmail.Sender
    // ... other fields
}

func NewMyHandler(gmailSender gmail.Sender) *MyHandler {
    return &MyHandler{
        gmailSender: gmailSender,
    }
}
```

### 3. Send Emails

Use the `Send()` method to send emails:

```go
msg := gmail.EmailMessage{
    To:      []string{"recipient@example.com"},
    Subject: "Test Email",
    Body:    "<h1>Hello</h1><p>This is a test email.</p>",
}

if err := h.gmailSender.Send(msg); err != nil {
    log.Printf("Failed to send email: %v", err)
    return
}
```

## Interface

The `Sender` interface is defined in `interface.go`:

```go
type Sender interface {
    Send(msg EmailMessage) error
}
```

This allows for easy testing by creating mock implementations of the interface.

## EmailMessage Structure

```go
type EmailMessage struct {
    To      []string // Recipient email addresses
    Subject string   // Email subject
    Body    string   // Email body (HTML supported)
}
```

