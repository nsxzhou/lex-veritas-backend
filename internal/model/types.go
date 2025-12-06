package model

// ============================================================================
// 枚举类型定义
// ============================================================================

// UserRole 用户角色
type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

// UserStatus 用户状态
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

// MessageRole 消息角色
type MessageRole string

const (
	RoleUserMsg   MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// DocumentType 文档类型
type DocumentType string

const (
	DocTypePDF      DocumentType = "pdf"
	DocTypeDOCX     DocumentType = "docx"
	DocTypeTXT      DocumentType = "txt"
	DocTypeURL      DocumentType = "url"
	DocTypeMarkdown DocumentType = "markdown"
)

// DocumentStatus 文档状态
type DocumentStatus string

const (
	DocStatusPending    DocumentStatus = "pending"
	DocStatusProcessing DocumentStatus = "processing"
	DocStatusIndexed    DocumentStatus = "indexed"
	DocStatusMinted     DocumentStatus = "minted"
	DocStatusError      DocumentStatus = "error"
)

// AuditType 审计类型
type AuditType string

const (
	AuditTamper AuditType = "tamper"
	AuditAccess AuditType = "access"
	AuditVerify AuditType = "verify"
	AuditUpload AuditType = "upload"
	AuditMint   AuditType = "mint"
	AuditSystem AuditType = "system"
	AuditAuth   AuditType = "auth"
)

// Severity 严重程度
type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// AuditStatus 审计状态
type AuditStatus string

const (
	AuditUnresolved    AuditStatus = "unresolved"
	AuditInvestigating AuditStatus = "investigating"
	AuditResolved      AuditStatus = "resolved"
)
