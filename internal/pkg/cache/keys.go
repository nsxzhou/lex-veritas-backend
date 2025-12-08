package cache

import "time"

// ============================================================================
// Redis 键前缀常量
// 所有模块的 Redis 键前缀统一在此定义，便于维护和避免冲突
// ============================================================================

const (
	// ========== 认证模块 ==========

	// KeyAuthRefreshToken 刷新令牌存储键
	// 格式: auth:refresh:{tokenHash}
	// 值: RefreshTokenData JSON
	// TTL: 7 days
	KeyAuthRefreshToken = "auth:refresh:"

	// KeyAuthBlacklist Token 黑名单键 (logout 后使 access token 失效)
	// 格式: auth:blacklist:{jti}
	// 值: "1"
	// TTL: 等于 access token 剩余有效期
	KeyAuthBlacklist = "auth:blacklist:"

	// KeyAuthLoginAttempts 登录尝试计数键 (防暴力破解)
	// 格式: auth:attempts:{identifier}
	// 值: 尝试次数
	// TTL: 15 minutes
	KeyAuthLoginAttempts = "auth:attempts:"

	// ========== 会话模块 ==========

	// KeySessionUser 用户会话缓存键
	// 格式: session:user:{userId}
	// 值: User JSON
	KeySessionUser = "session:user:"

	// ========== 限流模块 ==========

	// KeyRateLimitIP IP 限流计数键
	// 格式: ratelimit:ip:{ip}
	KeyRateLimitIP = "ratelimit:ip:"

	// KeyRateLimitUser 用户限流计数键
	// 格式: ratelimit:user:{userId}
	KeyRateLimitUser = "ratelimit:user:"
)

// ============================================================================
// 认证配置常量
// ============================================================================

const (
	// AuthMaxLoginAttempts 最大登录尝试次数
	AuthMaxLoginAttempts = 5

	// AuthLockoutDuration 账户锁定时长
	AuthLockoutDuration = 15 * time.Minute

	// AuthLoginAttemptsTTL 登录尝试记录过期时间 (与锁定时长相同)
	AuthLoginAttemptsTTL = 15 * time.Minute
)

// ============================================================================
// 匿名用户配置常量
// ============================================================================

const (
	// KeyGuestSession 匿名用户会话键
	// 格式: guest:session:{guestId}
	// 值: { sessionId, chatCount }
	KeyGuestSession = "guest:session:"

	// GuestMaxSessions 匿名用户最大会话数
	GuestMaxSessions = 1

	// GuestMaxChatsPerSession 匿名用户每会话最大对话数
	GuestMaxChatsPerSession = 5

	// GuestSessionTTL 匿名会话过期时间
	GuestSessionTTL = 24 * time.Hour
)

// ============================================================================
// 键生成辅助函数
// ============================================================================

// AuthRefreshTokenKey 生成刷新令牌存储键
func AuthRefreshTokenKey(tokenHash string) string {
	return KeyAuthRefreshToken + tokenHash
}

// AuthBlacklistKey 生成 Token 黑名单键
func AuthBlacklistKey(jti string) string {
	return KeyAuthBlacklist + jti
}

// AuthLoginAttemptsKey 生成登录尝试计数键
func AuthLoginAttemptsKey(identifier string) string {
	return KeyAuthLoginAttempts + identifier
}

// SessionUserKey 生成用户会话缓存键
func SessionUserKey(userID string) string {
	return KeySessionUser + userID
}

// GuestSessionKey 生成匿名用户会话键
func GuestSessionKey(guestID string) string {
	return KeyGuestSession + guestID
}
