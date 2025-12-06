package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Tenant 租户/组织
type Tenant struct {
	ID   string `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
	Slug string `json:"slug" gorm:"type:varchar(50);uniqueIndex;not null"`

	// 配置
	Plan       string `json:"plan" gorm:"type:varchar(20);default:'free'"`
	TokenQuota int64  `json:"tokenQuota" gorm:"default:1000000"`
	TokenUsed  int64  `json:"tokenUsed" gorm:"default:0"`
	MaxUsers   int    `json:"maxUsers" gorm:"default:5"`

	// 自定义设置
	Settings datatypes.JSON `json:"settings,omitempty" gorm:"type:jsonb"`

	// 状态
	Status string `json:"status" gorm:"type:varchar(20);default:'active'"`

	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TenantMember 租户成员关系
type TenantMember struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID string `json:"tenantId" gorm:"type:varchar(36);index;not null"`
	UserID   string `json:"userId" gorm:"type:varchar(36);index;not null"`
	Role     string `json:"role" gorm:"type:varchar(20);default:'member'"`

	JoinedAt time.Time `json:"joinedAt" gorm:"autoCreateTime"`
}

func (TenantMember) TableName() string {
	return "tenant_members"
}
