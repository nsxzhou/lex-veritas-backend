package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultBcryptCost bcrypt 默认 cost 值
	DefaultBcryptCost = 12
	// MinPasswordLength 最小密码长度
	MinPasswordLength = 8
)

var (
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak    = errors.New("password must contain uppercase, lowercase, digit and special character")
	ErrPasswordHashFailed = errors.New("failed to hash password")
	ErrPasswordMismatch   = errors.New("password does not match")
)

// PasswordConfig 密码配置
type PasswordConfig struct {
	BcryptCost int
}

// PasswordManager 密码管理器
type PasswordManager struct {
	cost int
}

// NewPasswordManager 创建密码管理器
func NewPasswordManager(cfg *PasswordConfig) *PasswordManager {
	cost := cfg.BcryptCost
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = DefaultBcryptCost
	}
	return &PasswordManager{cost: cost}
}

// HashPassword 对密码进行哈希
func (m *PasswordManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), m.cost)
	if err != nil {
		return "", ErrPasswordHashFailed
	}
	return string(bytes), nil
}

// VerifyPassword 验证密码
func (m *PasswordManager) VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
}
