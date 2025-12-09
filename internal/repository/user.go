package repository

import (
	"context"
	"errors"

	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/database"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// 查询
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// 创建/更新
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error

	// 额度相关
	IncrementTokenUsed(ctx context.Context, id string, amount int64) error
	GetQuotaInfo(ctx context.Context, id string) (quota, used int64, err error)

	// 列表查询 (admin)
	List(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
	ListWithFilters(ctx context.Context, page, pageSize int, filters *dto.UserListRequest) ([]model.User, int64, error)

	// 删除
	Delete(ctx context.Context, id string) error
}

// userRepository 用户数据访问实现
type userRepository struct{}

// NewUserRepository 创建用户数据访问实例
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// FindByID 按 ID 查询用户
func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := database.DB().WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail 按邮箱查询用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := database.DB().WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone 按手机号查询用户
func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	if err := database.DB().WithContext(ctx).First(&user, "phone = ?", phone).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := database.DB().WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return database.DB().WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return database.DB().WithContext(ctx).Save(user).Error
}

// UpdateFields 更新指定字段
func (r *userRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	result := database.DB().WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// IncrementTokenUsed 增加已使用 Token 数量
func (r *userRepository) IncrementTokenUsed(ctx context.Context, id string, amount int64) error {
	result := database.DB().WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("token_used", gorm.Expr("token_used + ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// GetQuotaInfo 获取额度信息
func (r *userRepository) GetQuotaInfo(ctx context.Context, id string) (quota, used int64, err error) {
	var user model.User
	if err = database.DB().WithContext(ctx).
		Select("token_quota", "token_used").
		First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, ErrUserNotFound
		}
		return 0, 0, err
	}
	return user.TokenQuota, user.TokenUsed, nil
}

// List 分页查询用户列表
func (r *userRepository) List(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.DB().WithContext(ctx).Model(&model.User{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListWithFilters 带筛选条件的分页查询
func (r *userRepository) ListWithFilters(ctx context.Context, page, pageSize int, filters *dto.UserListRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.DB().WithContext(ctx).Model(&model.User{})

	// 应用筛选条件
	if filters.Email != "" {
		db = db.Where("email = ?", filters.Email)
	}
	if filters.Role != "" {
		db = db.Where("role = ?", filters.Role)
	}
	if filters.Status != "" {
		db = db.Where("status = ?", filters.Status)
	}
	if filters.Keyword != "" {
		// 模糊搜索邮箱或姓名
		db = db.Where("email LIKE ? OR name LIKE ?", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Delete 软删除用户
func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := database.DB().WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
