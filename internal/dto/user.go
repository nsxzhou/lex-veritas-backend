package dto

// ============================================================================
// 用户额度相关 DTO
// ============================================================================

// AdjustQuotaRequest 调整额度请求
type AdjustQuotaRequest struct {
	UserID   string `json:"userId" binding:"required"`
	NewQuota int64  `json:"newQuota" binding:"required,min=0"`
}

// UsageStats 用量统计
type UsageStats struct {
	UserID     string  `json:"userId"`
	TokenQuota int64   `json:"tokenQuota"`
	TokenUsed  int64   `json:"tokenUsed"`
	Remaining  int64   `json:"remaining"`
	UsageRate  float64 `json:"usageRate"` // 0-1
}

// UserUsage 用户用量 (管理员视图)
type UserUsage struct {
	UsageStats
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UsageListResponse 用量列表响应
type UsageListResponse struct {
	List  []UserUsage `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"pageSize"`
}
