package dto

import "time"

// ============================================================================
// 认证相关 DTO
// ============================================================================

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// PhoneLoginRequest 手机登录请求
type PhoneLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone,omitempty"`
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ============================================================================
// 认证响应 DTO
// ============================================================================

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // 秒
	TokenType    string `json:"tokenType"` // Bearer
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  interface{} `json:"user"`
	Token *TokenPair  `json:"token"`
}

// UserResponse 用户信息响应
type UserResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone,omitempty"`
	Name       string     `json:"name"`
	Avatar     string     `json:"avatar,omitempty"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	TokenQuota int64      `json:"tokenQuota"`
	TokenUsed  int64      `json:"tokenUsed"`
	LastLogin  *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}
