package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidSignature = errors.New("invalid signature")
)

// Claims JWT 声明结构
type Claims struct {
	UserID   string `json:"uid"`
	Role     string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
	Issuer        string
}

// JWTManager JWT 管理器
type JWTManager struct {
	config    *JWTConfig
	secretKey []byte
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(cfg *JWTConfig) *JWTManager {
	return &JWTManager{
		config:    cfg,
		secretKey: []byte(cfg.Secret),
	}
}

// GenerateAccessToken 生成访问令牌
func (m *JWTManager) GenerateAccessToken(userID, role string) (string, error) {
	now := time.Now()
	claims := &Claims{		
		UserID:   userID,
			Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessExpire)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// GenerateRefreshToken 生成刷新令牌 (随机 UUID)
func (m *JWTManager) GenerateRefreshToken() string {
	return uuid.New().String()
}

// ParseToken 解析并验证令牌
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetTokenID 从 token 中提取 JTI (用于黑名单)
func (m *JWTManager) GetTokenID(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return "", err
	}
	if claims == nil {
		// 如果 token 过期，仍尝试获取 JTI
		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
		if token != nil {
			if c, ok := token.Claims.(*Claims); ok {
				return c.ID, nil
			}
		}
		return "", ErrInvalidToken
	}
	return claims.ID, nil
}

// GetExpiresIn 返回 access token 过期时间（秒）
func (m *JWTManager) GetExpiresIn() int64 {
	return int64(m.config.AccessExpire.Seconds())
}

// GetRefreshExpiresIn 返回 refresh token 过期时间
func (m *JWTManager) GetRefreshExpire() time.Duration {
	return m.config.RefreshExpire
}
