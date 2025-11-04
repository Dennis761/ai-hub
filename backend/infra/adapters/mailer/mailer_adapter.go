package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"ai_hub.com/app/core/ports/adminports"
)

var _ adminports.Mailer = (*SMTPMailer)(nil)

type SMTPMailer struct {
	host string
	user string
	pass string
}

func NewSMTPMailer(host, user, pass string) *SMTPMailer {
	return &SMTPMailer{
		host: host,
		user: user,
		pass: pass,
	}
}

func (m *SMTPMailer) send(_ context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:587", m.host)

	msg := []byte(fmt.Sprintf(
		"From: AI Hub <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		m.user, to, subject, body,
	))

	auth := smtp.PlainAuth("", m.user, m.pass, m.host)
	if err := smtp.SendMail(addr, auth, m.user, []string{to}, msg); err != nil {
		return fmt.Errorf("send mail error: %w", err)
	}
	return nil
}

func (m *SMTPMailer) SendVerificationCode(ctx context.Context, email, code string, expiresAt time.Time) error {
	body := fmt.Sprintf("Your verification code is: %s\nExpires: %s",
		code, expiresAt.Format(time.RFC3339))
	return m.send(ctx, email, "Registration Confirmation", body)
}

func (m *SMTPMailer) SendResetCode(ctx context.Context, email, code string, expiresAt time.Time) error {
	body := fmt.Sprintf("Your reset code is: %s\nExpires: %s",
		code, expiresAt.Format(time.RFC3339))
	return m.send(ctx, email, "Password Reset Code", body)
}
