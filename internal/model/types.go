// Package model 定义系统中使用的数据结构
package model

import "time"

// User 用户模型
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Phone     string    `json:"phone,omitempty"`
	Password  string    `json:"-"` // 不序列化到 JSON
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChatSession 聊天会话
type ChatSession struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userId" gorm:"index"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	SessionID string     `json:"sessionId" gorm:"index"`
	Role      string     `json:"role"` // user | assistant
	Content   string     `json:"content"`
	Citations []Citation `json:"citations,omitempty" gorm:"-"`
	CreatedAt time.Time  `json:"createdAt"`
}

// Citation 引用信息
type Citation struct {
	ID             string            `json:"id"`
	Text           string            `json:"text"`
	Source         string            `json:"source"`
	VerificationID string            `json:"verificationId"`
	BlockNumber    int64             `json:"blockNumber"`
	Timestamp      time.Time         `json:"timestamp"`
	Metadata       map[string]string `json:"metadata"`
}

// Document 知识库文档
type Document struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Type       string    `json:"type"` // pdf | docx | txt | url
	Size       string    `json:"size"`
	Status     string    `json:"status"` // indexed | processing | error | minted
	UploadedBy string    `json:"uploadedBy"`
	UploadDate time.Time `json:"uploadDate"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Chunk 文档分块
type Chunk struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	DocumentID  string    `json:"documentId" gorm:"index"`
	Content     string    `json:"content"`
	Hash        string    `json:"hash"`
	MerkleProof []string  `json:"merkleProof" gorm:"-"`
	Embedding   []float32 `json:"-" gorm:"-"` // 向量数据存储在 Milvus
	CreatedAt   time.Time `json:"createdAt"`
}

// AuditLog 审计日志
type AuditLog struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type"`     // tamper | access | verify | system
	Severity  string    `json:"severity"` // high | medium | low
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Status    string    `json:"status"` // unresolved | investigating | resolved
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"createdAt"`
}

// ProofRecord 存证记录
type ProofRecord struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	MerkleRoot  string    `json:"merkleRoot"`
	TxHash      string    `json:"txHash"`
	BlockNumber int64     `json:"blockNumber"`
	VersionID   int       `json:"versionId"`
	DocumentIDs []string  `json:"documentIds" gorm:"-"`
	CreatedAt   time.Time `json:"createdAt"`
}
