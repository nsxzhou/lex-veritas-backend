package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ChatSession 聊天会话
type ChatSession struct {
	ID     string `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID string `json:"userId,omitempty" gorm:"type:varchar(36);index"` // 登录用户 ID (可选)

	// 匿名用户标识
	GuestSessionID string `json:"guestSessionId,omitempty" gorm:"type:varchar(36);index"`

	Title   string `json:"title" gorm:"type:varchar(200)"`
	Summary string `json:"summary,omitempty" gorm:"type:text"`

	// 统计
	MessageCount int `json:"messageCount" gorm:"default:0"`
	TokensUsed   int `json:"tokensUsed" gorm:"default:0"`

	// 时间戳
	LastMessageAt *time.Time     `json:"lastMessageAt,omitempty"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	ID        string      `json:"id" gorm:"type:varchar(36);primaryKey"`
	SessionID string      `json:"sessionId" gorm:"type:varchar(36);index;not null"`
	Role      MessageRole `json:"role" gorm:"type:varchar(20);not null"`
	Content   string      `json:"content" gorm:"type:text;not null"`

	// AI 生成统计
	TokensIn  int `json:"tokensIn,omitempty"`
	TokensOut int `json:"tokensOut,omitempty"`

	// 引用关系
	Citations []MessageCitation `json:"citations,omitempty" gorm:"foreignKey:MessageID"`

	// 元数据
	Metadata datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// MessageCitation 消息中的法律引用
type MessageCitation struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	MessageID string `json:"messageId" gorm:"type:varchar(36);index;not null"`
	ChunkID   string `json:"chunkId" gorm:"type:varchar(64);index"`

	// 引用内容快照
	Text          string `json:"text" gorm:"type:text"`
	Source        string `json:"source" gorm:"type:varchar(200)"`
	ArticleNumber string `json:"articleNumber,omitempty" gorm:"type:varchar(32)"`
	LawHierarchy  string `json:"lawHierarchy,omitempty" gorm:"type:varchar(500)"`

	// 区块链验证信息
	ChunkHash      string     `json:"chunkHash" gorm:"type:varchar(64)"`
	VerificationID string     `json:"verificationId,omitempty" gorm:"type:varchar(66)"`
	BlockNumber    int64      `json:"blockNumber,omitempty"`
	Verified       bool       `json:"verified" gorm:"default:false"`
	VerifiedAt     *time.Time `json:"verifiedAt,omitempty"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
