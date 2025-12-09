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

// ============================================================================
// 管理员功能 DTO
// ============================================================================

// UserListRequest 用户列表查询请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Email    string `form:"email"`
	Role     string `form:"role" binding:"omitempty,oneof=user admin super_admin"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive banned"`
	Keyword  string `form:"keyword"` // 搜索邮箱或姓名
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	List  []UserResponse `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"pageSize"`
}

// UpdateStatusRequest 更新用户状态请求
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive banned"`
}

// UpdateRoleRequest 更新用户角色请求
type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin super_admin"`
}
