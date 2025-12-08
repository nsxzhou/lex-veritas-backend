package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Document 知识库文档
type Document struct {
	ID string `json:"id" gorm:"type:varchar(36);primaryKey"`

	// 基础信息
	Name         string       `json:"name" gorm:"type:varchar(255);not null"`
	OriginalName string       `json:"originalName" gorm:"type:varchar(255)"`
	Type         DocumentType `json:"type" gorm:"type:varchar(20);not null"`
	MimeType     string       `json:"mimeType,omitempty" gorm:"type:varchar(100)"`
	Size         int64        `json:"size"`
	FilePath     string       `json:"-" gorm:"type:varchar(500)"`

	// 法律元数据
	LawName       string     `json:"lawName,omitempty" gorm:"type:varchar(200)"`
	LawType       string     `json:"lawType,omitempty" gorm:"type:varchar(50)"`
	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`
	PublishOrg    string     `json:"publishOrg,omitempty" gorm:"type:varchar(100)"`

	// 处理状态
	Status       DocumentStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ProcessError string         `json:"processError,omitempty" gorm:"type:text"`

	// 分块统计
	ChunkCount int `json:"chunkCount" gorm:"default:0"`

	// 区块链存证
	IsMinted  bool       `json:"isMinted" gorm:"default:false"`
	MintedAt  *time.Time `json:"mintedAt,omitempty"`
	VersionID int        `json:"versionId,omitempty"`

	// 上传信息
	UploadedBy string `json:"uploadedBy" gorm:"type:varchar(36);index"`

	// 扩展元数据
	Metadata datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb"`

	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// DocumentChunk 文档分块元数据 (向量存储在 Milvus)
type DocumentChunk struct {
	ID         int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ChunkID    string `json:"chunkId" gorm:"type:varchar(64);uniqueIndex;not null"`
	DocumentID string `json:"documentId" gorm:"type:varchar(36);index;not null"`

	// 内容与哈希
	Content     string `json:"content" gorm:"type:text;not null"`
	ContentHash string `json:"contentHash" gorm:"type:varchar(64);index;not null"`

	// 法律结构
	ChunkOrder    int    `json:"chunkOrder"`
	LawHierarchy  string `json:"lawHierarchy,omitempty" gorm:"type:varchar(500)"`
	ArticleNumber string `json:"articleNumber,omitempty" gorm:"type:varchar(32)"`

	// 引用关系
	References datatypes.JSON `json:"references,omitempty" gorm:"type:jsonb"`

	// Merkle 验证数据
	MerkleIndex int            `json:"merkleIndex"`
	MerkleProof datatypes.JSON `json:"merkleProof,omitempty" gorm:"type:jsonb"`

	// 版本关联
	VersionID int `json:"versionId" gorm:"index"`

	// 向量化状态
	IsEmbedded bool       `json:"isEmbedded" gorm:"default:false"`
	EmbeddedAt *time.Time `json:"embeddedAt,omitempty"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
