package email

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"gopkg.in/gomail.v2"
)

// SMTPSender SMTP 邮件发送实现
type SMTPSender struct {
	host        string
	port        int
	user        string
	password    string
	fromAddress string
	fromName    string
}

// NewSMTPSender 创建 SMTP 发送器
func NewSMTPSender(cfg *config.EmailConfig) *SMTPSender {
	return &SMTPSender{
		host:        cfg.SMTPHost,
		port:        cfg.SMTPPort,
		user:        cfg.SMTPUser,
		password:    cfg.SMTPPassword,
		fromAddress: cfg.FromAddress,
		fromName:    cfg.FromName,
	}
}

// Send 发送邮件
func (s *SMTPSender) Send(ctx context.Context, to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.fromAddress, s.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.host, s.port, s.user, s.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Ensure SMTPSender implements Sender
var _ Sender = (*SMTPSender)(nil)
