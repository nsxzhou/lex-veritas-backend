package email

import "context"

// Sender 邮件发送接口
// 抽象邮件发送，支持 SMTP 和 Resend 实现
type Sender interface {
	// Send 发送邮件
	Send(ctx context.Context, to, subject, htmlBody string) error
}
