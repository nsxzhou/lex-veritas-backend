package service

import (
	"context"
	"errors"

	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/auth"
	"github.com/lexveritas/lex-veritas-backend/internal/repository"
)

var (
	ErrQuotaExceeded = errors.New("token quota exceeded")
)

// UserService 用户服务接口
type UserService interface {
	// 额度检查
	CheckQuota(ctx context.Context, userID string, required int64) error

	// 额度消费
	ConsumeTokens(ctx context.Context, userID string, amount int64) error

	// 额度调整 (admin)
	AdjustQuota(ctx context.Context, userID string, newQuota int64) error

	// 获取使用统计
	GetUsage(ctx context.Context, userID string) (*dto.UsageStats, error)
	GetAllUsage(ctx context.Context, page, pageSize int) (*dto.UsageListResponse, error)

	// 用户自服务
	UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID string, oldPwd, newPwd string) error

	// 管理员功能
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	ListUsers(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error)
	UpdateStatus(ctx context.Context, userID string, status string) error
	UpdateRole(ctx context.Context, userID string, role string) error
	DeleteUser(ctx context.Context, userID string) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
	}
}

// CheckQuota 检查额度是否足够
func (s *userService) CheckQuota(ctx context.Context, userID string, required int64) error {
	quota, used, err := s.userRepo.GetQuotaInfo(ctx, userID)
	if err != nil {
		return err
	}

	if used+required > quota {
		return ErrQuotaExceeded
	}

	return nil
}

// ConsumeTokens 消费 Token 额度
func (s *userService) ConsumeTokens(ctx context.Context, userID string, amount int64) error {
	return s.userRepo.IncrementTokenUsed(ctx, userID, amount)
}

// AdjustQuota 调整用户额度 (管理员操作)
func (s *userService) AdjustQuota(ctx context.Context, userID string, newQuota int64) error {
	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"token_quota": newQuota,
	})
}

// GetUsage 获取用户用量统计
func (s *userService) GetUsage(ctx context.Context, userID string) (*dto.UsageStats, error) {
	quota, used, err := s.userRepo.GetQuotaInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	remaining := quota - used
	var usageRate float64
	if quota > 0 {
		usageRate = float64(used) / float64(quota)
	}

	return &dto.UsageStats{
		UserID:     userID,
		TokenQuota: quota,
		TokenUsed:  used,
		Remaining:  remaining,
		UsageRate:  usageRate,
	}, nil
}

// GetAllUsage 获取所有用户用量 (管理员)
func (s *userService) GetAllUsage(ctx context.Context, page, pageSize int) (*dto.UsageListResponse, error) {
	users, total, err := s.userRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	list := make([]dto.UserUsage, len(users))
	for i, u := range users {
		remaining := u.TokenQuota - u.TokenUsed
		var usageRate float64
		if u.TokenQuota > 0 {
			usageRate = float64(u.TokenUsed) / float64(u.TokenQuota)
		}

		list[i] = dto.UserUsage{
			UsageStats: dto.UsageStats{
				UserID:     u.ID,
				TokenQuota: u.TokenQuota,
				TokenUsed:  u.TokenUsed,
				Remaining:  remaining,
				UsageRate:  usageRate,
			},
			Email: u.Email,
			Name:  u.Name,
		}
	}

	return &dto.UsageListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// ============================================================================
// 用户自服务实现
// ============================================================================

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) error {
	updates := make(map[string]interface{})

	if req.Name != nil {
		name := *req.Name
		// 非空时验证长度
		if name != "" && (len(name) < 2 || len(name) > 100) {
			return errors.New("姓名长度必须在2-100个字符之间")
		}
		updates["name"] = name
	}
	if req.Phone != nil {
		phone := *req.Phone
		// 非空时验证手机号格式(11位数字)
		if phone != "" && len(phone) != 11 {
			return errors.New("手机号必须为11位")
		}
		updates["phone"] = phone
	}
	if req.Avatar != nil {
		avatar := *req.Avatar
		// 非空时简单验证URL格式
		if avatar != "" && !isValidURL(avatar) {
			return errors.New("头像地址格式不正确")
		}
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		return nil // 没有更新
	}

	return s.userRepo.UpdateFields(ctx, userID, updates)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID string, oldPwd, newPwd string) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	pwdMgr := getPasswordManager()
	if err := pwdMgr.VerifyPassword(oldPwd, user.PasswordHash); err != nil {
		return errors.New("旧密码错误")
	}

	// 哈希新密码
	newHash, err := pwdMgr.HashPassword(newPwd)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"password_hash": newHash,
	})
}

// ============================================================================
// 管理员功能实现
// ============================================================================

// GetUserByID 获取用户详情
func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	// 设置默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 调用 Repository 层查询
	users, total, err := s.userRepo.ListWithFilters(ctx, page, pageSize, req)
	if err != nil {
		return nil, err
	}

	// 转换为 UserResponse
	list := make([]dto.UserResponse, len(users))
	for i, u := range users {
		list[i] = dto.UserResponse{
			ID:         u.ID,
			Email:      u.Email,
			Phone:      u.Phone,
			Name:       u.Name,
			Avatar:     u.Avatar,
			Role:       string(u.Role),
			Status:     string(u.Status),
			TokenQuota: u.TokenQuota,
			TokenUsed:  u.TokenUsed,
			LastLogin:  u.LastLoginAt,
			CreatedAt:  u.CreatedAt,
		}
	}

	return &dto.UserListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// UpdateStatus 修改用户状态
func (s *userService) UpdateStatus(ctx context.Context, userID string, status string) error {
	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"status": status,
	})
}

// UpdateRole 修改用户角色
func (s *userService) UpdateRole(ctx context.Context, userID string, role string) error {
	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"role": role,
	})
}

// DeleteUser 删除用户 (软删除)
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	return s.userRepo.Delete(ctx, userID)
}

// ============================================================================
// 内部辅助方法
// ============================================================================

// getPasswordManager 获取密码管理器
func getPasswordManager() *auth.PasswordManager {
	return auth.NewPasswordManager(&auth.PasswordConfig{
		BcryptCost: 10,
	})
}

// isValidURL 简单验证URL格式
func isValidURL(s string) bool {
	return len(s) > 7 && (s[:7] == "http://" || s[:8] == "https://")
}
