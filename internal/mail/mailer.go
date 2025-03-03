package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// Mailer handles sending emails
type Mailer struct {
	Host     string
	Port     string
	From     string
	Username string
	Password string
}

// NewMailer creates a new Mailer instance using environment variables
func NewMailer() *Mailer {
	return &Mailer{
		Host:     getEnv("SMTP_HOST", "localhost"),
		Port:     getEnv("SMTP_PORT", "1025"),
		From:     getEnv("SMTP_FROM", "admin@unixify.example.com"),
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
	}
}

// SendEmail sends an email
func (m *Mailer) SendEmail(to []string, subject, body string) error {
	// Set up authentication information if credentials are provided
	var auth smtp.Auth
	if m.Username \!= "" && m.Password \!= "" {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}

	// Format the email message
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", m.From, to[0], subject, body))

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	
	if err := smtp.SendMail(addr, auth, m.From, to, msg); err \!= nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent to %v", to)
	return nil
}

// SendVerificationEmail sends an email verification link
func (m *Mailer) SendVerificationEmail(to, username, token string) error {
	subject := "Verify Your Unixify Account"
	
	// In a real app, this would be a link to your verification endpoint
	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)
	
	body := fmt.Sprintf("<html><body><h2>Welcome to Unixify, %s\!</h2><p>Please verify your account by clicking the link below:</p><p><a href=\"%s\">Verify Account</a></p><p>If you didn't create this account, you can safely ignore this email.</p><p>The Unixify Team</p></body></html>", username, verificationLink)
	
	return m.SendEmail([]string{to}, subject, body)
}

// SendPasswordResetEmail sends a password reset link
func (m *Mailer) SendPasswordResetEmail(to, username, token string) error {
	subject := "Reset Your Unixify Password"
	
	// In a real app, this would be a link to your password reset page
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)
	
	body := fmt.Sprintf("<html><body><h2>Password Reset Request</h2><p>Hello %s,</p><p>We received a request to reset your password. Click the link below to create a new password:</p><p><a href=\"%s\">Reset Password</a></p><p>If you didn't request this, you can safely ignore this email.</p><p>The Unixify Team</p></body></html>", username, resetLink)
	
	return m.SendEmail([]string{to}, subject, body)
}

// SendWelcomeEmail sends a welcome email to a new user
func (m *Mailer) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to Unixify\!"
	
	body := fmt.Sprintf("<html><body><h2>Welcome to Unixify, %s\!</h2><p>Thank you for registering an account with us.</p><p>You can now manage your UNIX accounts and groups with our comprehensive interface.</p><p>If you have any questions, please don't hesitate to contact us.</p><p>The Unixify Team</p></body></html>", username)
	
	return m.SendEmail([]string{to}, subject, body)
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
