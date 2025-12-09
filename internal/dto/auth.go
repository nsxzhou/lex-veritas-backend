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

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required,oneof=register reset_password"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6"`
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

// ============================================================================
// 用户自服务 DTO
// ============================================================================

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// UpdateProfileRequest 更新用户资料
type UpdateProfileRequest struct {
	Name   *string `json:"name,omitempty"`   // 允许空字符串以支持清空操作
	Phone  *string `json:"phone,omitempty"`  // 允许空字符串以支持清空操作
	Avatar *string `json:"avatar,omitempty"` // 允许空字符串以支持清空操作
}
