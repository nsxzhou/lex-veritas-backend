package model

import (
	"time"

	"gorm.io/datatypes"
)

// AuditLog 审计日志
type AuditLog struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// 事件类型
	Type     AuditType `json:"type" gorm:"type:varchar(30);index;not null"`
	Severity Severity  `json:"severity" gorm:"type:varchar(10);index;not null"`

	// 事件描述
	Message string         `json:"message" gorm:"type:text;not null"`
	Details datatypes.JSON `json:"details,omitempty" gorm:"type:jsonb"`

	// 来源信息
	Source     string `json:"source,omitempty" gorm:"type:varchar(255)"`
	SourceType string `json:"sourceType,omitempty" gorm:"type:varchar(30)"`

	// 关联用户
	UserID    string `json:"userId,omitempty" gorm:"type:varchar(36);index"`
	UserIP    string `json:"userIp,omitempty" gorm:"type:varchar(45)"`
	UserAgent string `json:"userAgent,omitempty" gorm:"type:varchar(500)"`

	// 处理状态
	Status     AuditStatus `json:"status" gorm:"type:varchar(20);default:'unresolved'"`
	ResolvedBy string      `json:"resolvedBy,omitempty" gorm:"type:varchar(36)"`
	ResolvedAt *time.Time  `json:"resolvedAt,omitempty"`
	Resolution string      `json:"resolution,omitempty" gorm:"type:text"`

	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime;index"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Key         string    `json:"key" gorm:"type:varchar(100);uniqueIndex;not null"`
	Value       string    `json:"value" gorm:"type:text"`
	Type        string    `json:"type" gorm:"type:varchar(20)"`
	Category    string    `json:"category" gorm:"type:varchar(50);index"`
	Description string    `json:"description,omitempty" gorm:"type:varchar(500)"`
	UpdatedBy   string    `json:"updatedBy,omitempty" gorm:"type:varchar(36)"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
