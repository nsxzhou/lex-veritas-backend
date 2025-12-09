package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/auth"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/cache"
	"github.com/lexveritas/lex-veritas-backend/internal/repository"
)

var (
	ErrUserNotFound        = errors.New("用户不存在")
	ErrUserDisabled        = errors.New("用户已被禁用或封禁")
	ErrInvalidCredentials  = errors.New("邮箱或密码错误")
	ErrAccountLocked       = errors.New("由于尝试次数过多，账户已被锁定")
	ErrRefreshTokenInvalid = errors.New("刷新令牌无效或已过期")
	ErrEmailAlreadyExists  = errors.New("邮箱已被注册")
)

// refreshTokenData 刷新令牌数据 (存储在 Redis)
type refreshTokenData struct {
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

// AuthService 认证服务接口
type AuthService interface {
	// 注册与登录
	Register(ctx context.Context, req *dto.RegisterRequest) (*model.User, error)
	LoginByEmail(ctx context.Context, email, password string) (*dto.TokenPair, *model.User, error)
	LoginByPhone(ctx context.Context, phone, code string) (*dto.TokenPair, *model.User, error)

	// Token 管理
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
	LogoutAll(ctx context.Context, userID string) error

	// 用户信息
	GetCurrentUser(ctx context.Context, userID string) (*model.User, error)

	// OAuth 登录
	OAuthLogin(ctx context.Context, provider, code string) (*dto.TokenPair, *model.User, error)
	OAuthCallback(ctx context.Context, provider, code, state string) (*dto.TokenPair, *model.User, error)

	// Token 验证
	ValidateAccessToken(token string) (*auth.Claims, error)
	IsTokenBlacklisted(ctx context.Context, jti string) bool
}

// authService 认证服务实现
type authService struct {
	jwtMgr   *auth.JWTManager
	pwdMgr   *auth.PasswordManager
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService(jwtCfg *auth.JWTConfig, pwdCfg *auth.PasswordConfig) AuthService {
	return &authService{
		jwtMgr:   auth.NewJWTManager(jwtCfg),
		pwdMgr:   auth.NewPasswordManager(pwdCfg),
		userRepo: repository.NewUserRepository(),
	}
}

// Register 用户注册
func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*model.User, error) {
	// 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// 验证密码强度
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// 哈希密码
	hashedPassword, err := s.pwdMgr.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &model.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Role:         model.RoleUser,
		Status:       model.StatusActive,
		TokenQuota:   100000, // 默认额度
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// LoginByEmail 邮箱密码登录
func (s *authService) LoginByEmail(ctx context.Context, email, password string) (*dto.TokenPair, *model.User, error) {
	// 检查是否被锁定
	if s.isAccountLocked(ctx, email) {
		return nil, nil, ErrAccountLocked
	}

	// 查询用户
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.incrementLoginAttempts(ctx, email)
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	// 验证密码
	if err := s.pwdMgr.VerifyPassword(password, user.PasswordHash); err != nil {
		s.incrementLoginAttempts(ctx, email)
		return nil, nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != model.StatusActive {
		return nil, nil, ErrUserDisabled
	}

	// 清除登录尝试计数
	s.clearLoginAttempts(ctx, email)

	// 更新最后登录时间
	now := time.Now()
	_ = s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"last_login_at": now,
	})

	// 生成令牌对
	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return tokenPair, user, nil
}

// LoginByPhone 手机验证码登录 (占位实现)
func (s *authService) LoginByPhone(ctx context.Context, phone, code string) (*dto.TokenPair, *model.User, error) {
	return nil, nil, errors.New("手机验证码登录尚未实现")
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPair, error) {
	tokenHash := hashToken(refreshToken)
	key := cache.AuthRefreshTokenKey(tokenHash)

	// 从 Redis 获取 refresh token 数据
	var data refreshTokenData
	if err := cache.GetObject(ctx, key, &data); err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// 删除已使用的 refresh token
	_ = cache.Delete(ctx, key)

	// 获取用户
	user, err := s.userRepo.FindByID(ctx, data.UserID)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// 生成新的令牌对
	return s.generateTokenPair(ctx, user)
}

// Logout 登出 (使 access token 失效)
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	jti, err := s.jwtMgr.GetTokenID(accessToken)
	if err != nil {
		return nil // 忽略无效 token
	}

	// 将 token 加入黑名单
	key := cache.AuthBlacklistKey(jti)
	return cache.Set(ctx, key, "1", s.jwtMgr.GetRefreshExpire())
}

// LogoutAll 登出所有设备
func (s *authService) LogoutAll(ctx context.Context, userID string) error {
	// TODO: 实现删除所有 refresh tokens
	return nil
}

// GetCurrentUser 获取当前用户信息
func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// ValidateAccessToken 验证访问令牌
func (s *authService) ValidateAccessToken(token string) (*auth.Claims, error) {
	return s.jwtMgr.ParseToken(token)
}

// IsTokenBlacklisted 检查令牌是否在黑名单中
func (s *authService) IsTokenBlacklisted(ctx context.Context, jti string) bool {
	key := cache.AuthBlacklistKey(jti)
	count, _ := cache.Exists(ctx, key)
	return count > 0
}

// ========================= 内部辅助方法 =========================

// generateTokenPair 生成令牌对
func (s *authService) generateTokenPair(ctx context.Context, user *model.User) (*dto.TokenPair, error) {
	accessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := s.jwtMgr.GenerateRefreshToken()

	// 存储 refresh token 到 Redis
	tokenHash := hashToken(refreshToken)
	key := cache.AuthRefreshTokenKey(tokenHash)
	data := refreshTokenData{
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
	if err := cache.SetObject(ctx, key, data, s.jwtMgr.GetRefreshExpire()); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtMgr.GetExpiresIn(),
		TokenType:    "Bearer",
	}, nil
}

// isAccountLocked 检查账户是否被锁定
func (s *authService) isAccountLocked(ctx context.Context, identifier string) bool {
	key := cache.AuthLoginAttemptsKey(identifier)
	val, err := cache.Get(ctx, key)
	if err != nil || val == "" {
		return false
	}
	// 简单判断：如果尝试次数存在且超过限制则锁定
	var count int64
	fmt.Sscanf(val, "%d", &count)
	return count >= int64(cache.AuthMaxLoginAttempts)
}

// incrementLoginAttempts 增加登录尝试计数
func (s *authService) incrementLoginAttempts(ctx context.Context, identifier string) {
	key := cache.AuthLoginAttemptsKey(identifier)
	count, _ := cache.Incr(ctx, key)
	if count == 1 {
		_ = cache.Expire(ctx, key, cache.AuthLoginAttemptsTTL)
	}
}

// clearLoginAttempts 清除登录尝试计数
func (s *authService) clearLoginAttempts(ctx context.Context, identifier string) {
	key := cache.AuthLoginAttemptsKey(identifier)
	_ = cache.Delete(ctx, key)
}

// hashToken 对 token 进行哈希
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// ========================= OAuth 登录 (占位实现) =========================

// OAuthLogin OAuth 登录重定向 (占位实现)
func (s *authService) OAuthLogin(ctx context.Context, provider, code string) (*dto.TokenPair, *model.User, error) {
	// TODO: 实现 OAuth 登录逻辑
	// 1. 根据 provider 构建 OAuth 授权 URL
	// 2. 重定向用户到第三方授权页面
	return nil, nil, errors.New("OAuth 登录尚未实现,支持的 provider: google, github, wechat")
}

// OAuthCallback OAuth 回调处理 (占位实现)
func (s *authService) OAuthCallback(ctx context.Context, provider, code, state string) (*dto.TokenPair, *model.User, error) {
	// TODO: 实现 OAuth 回调逻辑
	// 1. 验证 state 防止 CSRF
	// 2. 使用 code 换取 access_token
	// 3. 获取用户信息
	// 4. 查询或创建本地用户
	// 5. 创建或更新 OAuthAccount
	// 6. 生成 JWT token
	return nil, nil, errors.New("OAuth 回调处理尚未实现")
}
