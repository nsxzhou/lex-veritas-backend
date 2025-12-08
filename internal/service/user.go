package service

import (
	"context"
	"errors"

	"github.com/lexveritas/lex-veritas-backend/internal/dto"
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
