package service

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

type EmailService struct {
	dialer *mail.Dialer
	from   string
}

func NewEmailService() *EmailService {
	from := os.Getenv("GMAIL_FROM")
	password := os.Getenv("GMAIL_APP_PASSWORD")

	// Default SMTP settings for Gmail
	host := "smtp.gmail.com"
	port := 587

	// Allow custom SMTP settings via environment variables via env
	if customHost := os.Getenv("SMTP_HOST"); customHost != "" {
		host = customHost
	}
	if customPort := os.Getenv("SMTP_PORT"); customPort != "" {
		if p, err := strconv.Atoi(customPort); err == nil {
			port = p
		}
	}

	dialer := mail.NewDialer(host, port, from, password)

	return &EmailService{
		dialer: dialer,
		from:   from,
	}
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(toEmail, userName, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`hi %s,
	You requested to reset your password. Use the following token to reset it:

	Reset Token: %s

	If you didn't request this, please ignore this email.

	Best regards,
	The Team`, userName, resetToken)

	return s.sendEmail(toEmail, subject, body)
}

// SendWelcomeEmail sends a welcoming email
func (s *EmailService) SendWelcomeEmail(toEmail, userName string) error {
	subtle := "Welcome to LiveChat"
	body := fmt.Sprintf(`Hello %s,

Welcome to our platform! We're excited to have you on board.

If you have any questions, feel free to reach out to us.

Best regards,
Livechat team`, userName)

	return s.sendEmail(toEmail, subtle, body)
}

// sendEmail is the internal method that handles the actual sending
func (s *EmailService) sendEmail(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return s.dialer.DialAndSend(m)
}

// SendMail sends a generic email
func (s *EmailService) SendMail(toEmail, subject, body string) error {
	return s.sendEmail(toEmail, subject, body)
}

// SendHTMLEmail sends an HTML formatted email
func (s *EmailService) SendHTMLEmail(to, subject, htmlBody string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", htmlBody)

	return s.dialer.DialAndSend(m)
}
