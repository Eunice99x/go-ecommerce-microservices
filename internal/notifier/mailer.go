package notifier

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

// SMTPMailer sends plain text emails; smtp.SendMail upgrades to TLS when the server supports it
type SMTPMailer struct {
	addr string
	from string
	auth smtp.Auth
}

func NewSMTPMailer(host, port, username, password, from string) *SMTPMailer {
	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	return &SMTPMailer{
		addr: net.JoinHostPort(host, port),
		from: from,
		auth: auth,
	}
}

// Send ignores ctx because net/smtp doesn't support it
func (m *SMTPMailer) Send(_ context.Context, to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", m.from, to, subject, body)

	if err := smtp.SendMail(m.addr, m.auth, m.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

// LogMailer only logs emails, for local dev without SMTP
type LogMailer struct{}

func (LogMailer) Send(_ context.Context, to, subject, _ string) error {
	log.Printf("notifier: [dev] email to %s: %q", to, subject)
	return nil
}
