package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID string `json:"id" gorm:"type:varchar(36);primaryKey"`

	// 基础信息
	Email        string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone        string `json:"phone,omitempty" gorm:"type:varchar(20);uniqueIndex"`
	PasswordHash string `json:"-" gorm:"type:varchar(255)"`
	Name         string `json:"name" gorm:"type:varchar(100)"`
	Avatar       string `json:"avatar,omitempty" gorm:"type:varchar(500)"`

	// 角色与状态
	Role   UserRole   `json:"role" gorm:"type:varchar(20);default:'user'"`
	Status UserStatus `json:"status" gorm:"type:varchar(20);default:'active'"`

	// 额度管理
	TokenQuota int64 `json:"tokenQuota" gorm:"default:100000"` // 默认 10 万 tokens
	TokenUsed  int64 `json:"tokenUsed" gorm:"default:0"`

	// 时间戳
	LastLoginAt *time.Time     `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	OAuthAccounts []OAuthAccount `json:"oauthAccounts,omitempty" gorm:"foreignKey:UserID"`
}

// OAuthAccount OAuth 第三方登录账户
type OAuthAccount struct {
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID string `json:"userId" gorm:"type:varchar(36);index;not null"`

	// OAuth 提供商信息
	Provider       string `json:"provider" gorm:"type:varchar(30);not null"`
	ProviderUserID string `json:"providerUserId" gorm:"type:varchar(100);not null"`

	// Token 信息
	AccessToken  string     `json:"-" gorm:"type:text"`
	RefreshToken string     `json:"-" gorm:"type:text"`
	TokenExpiry  *time.Time `json:"tokenExpiry,omitempty"`

	// 用户信息快照
	Email  string `json:"email,omitempty" gorm:"type:varchar(255)"`
	Name   string `json:"name,omitempty" gorm:"type:varchar(100)"`
	Avatar string `json:"avatar,omitempty" gorm:"type:varchar(500)"`

	// 元数据
	RawData datatypes.JSON `json:"-" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
