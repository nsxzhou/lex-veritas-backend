package model

import "time"

// TokenUsage Token 使用记录 - 按 API 细分统计
type TokenUsage struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// 归属
	TenantID  string `json:"tenantId" gorm:"type:varchar(36);index;not null"`
	UserID    string `json:"userId" gorm:"type:varchar(36);index;not null"`
	SessionID string `json:"sessionId,omitempty" gorm:"type:varchar(36);index"`

	// API 细分
	APIEndpoint string `json:"apiEndpoint" gorm:"type:varchar(100);index"`
	APIMethod   string `json:"apiMethod" gorm:"type:varchar(10)"`

	// Token 统计
	TokensIn    int `json:"tokensIn"`
	TokensOut   int `json:"tokensOut"`
	TotalTokens int `json:"totalTokens"`

	// 模型信息
	Model     string `json:"model" gorm:"type:varchar(50)"`
	ModelType string `json:"modelType" gorm:"type:varchar(20)"`

	// 费用估算
	EstimatedCost float64 `json:"estimatedCost,omitempty" gorm:"type:decimal(10,6)"`

	// 请求信息
	RequestID string `json:"requestId,omitempty" gorm:"type:varchar(36)"`
	Latency   int    `json:"latency,omitempty"`
	Success   bool   `json:"success" gorm:"default:true"`
	ErrorMsg  string `json:"errorMsg,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;index"`
}

// TokenUsageDaily Token 使用日统计 (用于 Dashboard)
type TokenUsageDaily struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// 维度
	TenantID    string    `json:"tenantId" gorm:"type:varchar(36);index;not null"`
	UserID      string    `json:"userId,omitempty" gorm:"type:varchar(36);index"`
	APIEndpoint string    `json:"apiEndpoint,omitempty" gorm:"type:varchar(100)"`
	Model       string    `json:"model,omitempty" gorm:"type:varchar(50)"`
	Date        time.Time `json:"date" gorm:"type:date;index;not null"`

	// 聚合统计
	RequestCount int64   `json:"requestCount"`
	TokensIn     int64   `json:"tokensIn"`
	TokensOut    int64   `json:"tokensOut"`
	TotalTokens  int64   `json:"totalTokens"`
	TotalCost    float64 `json:"totalCost,omitempty" gorm:"type:decimal(12,4)"`
}
