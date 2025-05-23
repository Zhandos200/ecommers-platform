package usecase

import (
	"fmt"
	"net/smtp"
	"os"
)

// Mailer уже должен быть объявлен как
// type Mailer interface { SendVerification(email, token string) error }

// SMTPMailer удовлетворяет интерфейсу Mailer
type SMTPMailer struct {
	host string
	port int
	user string
	pass string
	from string
}

func NewSMTPMailer() *SMTPMailer {
	p := 587
	return &SMTPMailer{
		host: os.Getenv("SMTP_HOST"),
		port: p,
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
		from: os.Getenv("MAIL_FROM"),
	}
}

func (m *SMTPMailer) SendVerification(email, token string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)

	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	link := fmt.Sprintf("%s/verify?token=%s", baseURL, token)

	msg := []byte(
		"Subject: Please verify your email\r\n" +
			"From: " + m.from + "\r\n" +
			"To: " + email + "\r\n\r\n" +
			"Click to verify your account:\r\n" + link + "\r\n",
	)
	return smtp.SendMail(addr, auth, m.from, []string{email}, msg)
}
