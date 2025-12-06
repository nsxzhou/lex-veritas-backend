package model

import "time"

// KnowledgeVersion 知识库版本
type KnowledgeVersion struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`

	// Merkle Tree
	MerkleRoot string `json:"merkleRoot" gorm:"type:varchar(66);uniqueIndex;not null"`
	ChunkCount int    `json:"chunkCount"`

	// 版本描述
	Description string `json:"description,omitempty" gorm:"type:varchar(500)"`

	// 区块链交易信息
	TxHash      string `json:"txHash,omitempty" gorm:"type:varchar(66)"`
	BlockNumber int64  `json:"blockNumber,omitempty"`

	// 状态
	Status string `json:"status" gorm:"type:varchar(20);default:'pending'"`

	// 关联的文档
	Documents []Document `json:"-" gorm:"many2many:version_documents;"`

	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	ConfirmedAt *time.Time `json:"confirmedAt,omitempty"`
}

// ProofRecord 存证验证记录
type ProofRecord struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// 验证主体
	ChunkID    string `json:"chunkId,omitempty" gorm:"type:varchar(64);index"`
	DocumentID string `json:"documentId,omitempty" gorm:"type:varchar(36);index"`

	// 验证数据
	LeafHash     string `json:"leafHash" gorm:"type:varchar(64);not null"`
	ComputedRoot string `json:"computedRoot" gorm:"type:varchar(66)"`
	OnChainRoot  string `json:"onChainRoot" gorm:"type:varchar(66)"`

	// 验证结果
	Verified  bool `json:"verified"`
	VersionID int  `json:"versionId" gorm:"index"`

	// 链上信息
	BlockNumber int64  `json:"blockNumber,omitempty"`
	TxHash      string `json:"txHash,omitempty" gorm:"type:varchar(66)"`

	// 触发来源
	TriggerType string `json:"triggerType" gorm:"type:varchar(20)"`
	TriggerBy   string `json:"triggerBy,omitempty" gorm:"type:varchar(36)"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
