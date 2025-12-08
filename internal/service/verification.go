package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/cache"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/email"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
	"go.uber.org/zap"
)

// 验证码相关错误
var (
	ErrCodeSendTooFrequent = errors.New("请稍后再发送验证码")
	ErrCodeInvalid         = errors.New("验证码错误")
	ErrCodeExpired         = errors.New("验证码已过期")
	ErrTooManyAttempts     = errors.New("验证尝试次数过多")
)

// VerificationService 验证码服务接口
type VerificationService interface {
	// SendCode 发送验证码到指定邮箱
	SendCode(ctx context.Context, email, purpose string) error

	// VerifyCode 校验验证码
	VerifyCode(ctx context.Context, email, code, purpose string) error
}

// verificationService 验证码服务实现
type verificationService struct {
	emailSender email.Sender
	cfg         *config.VerificationConfig
}

// NewVerificationService 创建验证码服务
func NewVerificationService(sender email.Sender, cfg *config.VerificationConfig) VerificationService {
	return &verificationService{
		emailSender: sender,
		cfg:         cfg,
	}
}

// SendCode 发送验证码
func (s *verificationService) SendCode(ctx context.Context, emailAddr, purpose string) error {
	// 检查发送频率限制
	limitKey := cache.VerificationLimitKey(emailAddr)
	exists, _ := cache.Exists(ctx, limitKey)
	if exists > 0 {
		return ErrCodeSendTooFrequent
	}

	// 生成验证码
	code, err := generateCode(s.cfg.CodeLength)
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	// 存储验证码到 Redis
	codeKey := cache.VerificationCodeKey(purpose, emailAddr)
	if err := cache.Set(ctx, codeKey, code, s.cfg.CodeExpire); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 设置发送频率限制
	if err := cache.Set(ctx, limitKey, "1", s.cfg.ResendDelay); err != nil {
		logger.Warn("设置发送频率限制失败", zap.Error(err))
	}

	// 发送邮件
	subject := getSubject(purpose)
	htmlBody := email.VerificationCodeTemplate(code, int(s.cfg.CodeExpire.Minutes()))

	if err := s.emailSender.Send(ctx, emailAddr, subject, htmlBody); err != nil {
		// 发送失败时删除验证码（允许重试）
		_ = cache.Delete(ctx, codeKey, limitKey)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	logger.Info("验证码已发送",
		zap.String("email", maskEmail(emailAddr)),
		zap.String("purpose", purpose),
	)

	return nil
}

// VerifyCode 校验验证码
func (s *verificationService) VerifyCode(ctx context.Context, emailAddr, code, purpose string) error {
	// 检查验证尝试次数
	attemptsKey := cache.VerificationAttemptsKey(emailAddr)
	attempts, _ := cache.Get(ctx, attemptsKey)
	if attempts != "" {
		var count int
		fmt.Sscanf(attempts, "%d", &count)
		if count >= s.cfg.MaxAttempts {
			return ErrTooManyAttempts
		}
	}

	// 获取存储的验证码
	codeKey := cache.VerificationCodeKey(purpose, emailAddr)
	storedCode, err := cache.Get(ctx, codeKey)
	if err != nil {
		// 增加尝试次数
		s.incrementAttempts(ctx, attemptsKey)
		return ErrCodeExpired
	}

	// 校验验证码
	if storedCode != code {
		s.incrementAttempts(ctx, attemptsKey)
		return ErrCodeInvalid
	}

	// 验证成功，删除验证码和尝试次数
	_ = cache.Delete(ctx, codeKey, attemptsKey)

	logger.Info("验证码校验成功",
		zap.String("email", maskEmail(emailAddr)),
		zap.String("purpose", purpose),
	)

	return nil
}

// incrementAttempts 增加验证尝试次数
func (s *verificationService) incrementAttempts(ctx context.Context, key string) {
	count, _ := cache.Incr(ctx, key)
	if count == 1 {
		_ = cache.Expire(ctx, key, s.cfg.CodeExpire)
	}
}

// generateCode 生成随机数字验证码
func generateCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}

// getSubject 根据用途获取邮件主题
func getSubject(purpose string) string {
	switch purpose {
	case "register":
		return "LexVeritas - 注册验证码"
	case "reset_password":
		return "LexVeritas - 重置密码验证码"
	default:
		return "LexVeritas - 验证码"
	}
}

// maskEmail 邮箱脱敏
func maskEmail(email string) string {
	if len(email) < 5 {
		return "***"
	}
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}
	if atIndex < 2 {
		return "***" + email[atIndex:]
	}
	return email[:2] + "***" + email[atIndex:]
}
