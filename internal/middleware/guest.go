package middleware

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/cache"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
)

// GuestSession 匿名用户会话数据
type GuestSession struct {
	SessionID string `json:"sessionId"` // 当前聊天会话 ID
	ChatCount int    `json:"chatCount"` // 已发送对话次数
}

const guestCookieName = "lex_guest_id"

// GuestLimit 匿名用户限制中间件
// - 检查会话数量 (最多 1 个)
// - 检查对话次数 (每会话最多 5 次)
func GuestLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果已登录，跳过限制
		if GetUserID(c) != "" {
			c.Next()
			return
		}

		guestID := GetOrCreateGuestID(c)
		ctx := c.Request.Context()

		// 获取匿名用户会话数据
		key := cache.GuestSessionKey(guestID)
		session, err := getGuestSession(ctx, key)
		if err != nil {
			// 新的匿名用户，允许通过
			c.Set("guestId", guestID)
			c.Set("guestSession", &GuestSession{})
			c.Next()
			return
		}

		// 检查对话次数限制
		if session.ChatCount >= cache.GuestMaxChatsPerSession {
			response.Forbidden(c, "您已达到免费对话限制，请登录以继续使用")
			c.Abort()
			return
		}

		c.Set("guestId", guestID)
		c.Set("guestSession", session)
		c.Next()
	}
}

// GetOrCreateGuestID 获取或创建匿名用户 ID
func GetOrCreateGuestID(c *gin.Context) string {
	// 尝试从 Cookie 获取
	if guestID, err := c.Cookie(guestCookieName); err == nil && guestID != "" {
		return guestID
	}

	// 生成新的 ID
	guestID := uuid.New().String()

	// 设置 Cookie (24 小时过期)
	c.SetCookie(
		guestCookieName,
		guestID,
		int(cache.GuestSessionTTL.Seconds()),
		"/",
		"",    // domain
		false, // secure (生产环境应为 true)
		true,  // httpOnly
	)

	return guestID
}

// GetGuestID 从上下文获取匿名用户 ID
func GetGuestID(c *gin.Context) string {
	if id, exists := c.Get("guestId"); exists {
		return id.(string)
	}
	return ""
}

// GetGuestSession 从上下文获取匿名会话
func GetGuestSession(c *gin.Context) *GuestSession {
	if session, exists := c.Get("guestSession"); exists {
		return session.(*GuestSession)
	}
	return nil
}

// IncrementGuestChatCount 增加匿名用户对话次数
func IncrementGuestChatCount(c *gin.Context, sessionID string) error {
	guestID := GetGuestID(c)
	if guestID == "" {
		return nil // 不是匿名用户
	}

	ctx := c.Request.Context()
	key := cache.GuestSessionKey(guestID)

	session, _ := getGuestSession(ctx, key)
	if session == nil {
		session = &GuestSession{}
	}

	session.SessionID = sessionID
	session.ChatCount++

	// 保存到 Redis
	data, _ := json.Marshal(session)
	return cache.Set(ctx, key, string(data), cache.GuestSessionTTL)
}

// getGuestSession 从 Redis 获取匿名会话数据
func getGuestSession(ctx context.Context, key string) (*GuestSession, error) {
	val, err := cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var session GuestSession
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}

	return &session, nil
}
