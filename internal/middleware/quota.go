package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/database"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
)

// QuotaCheck 额度检查中间件
// 检查登录用户是否有足够的 Token 额度
func QuotaCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == "" {
			// 匿名用户由 GuestLimit 处理
			c.Next()
			return
		}

		// 管理员无限制
		if GetRole(c) == model.RoleAdmin {
			c.Next()
			return
		}

		// 查询用户额度
		var user model.User
		if err := database.DB().Select("token_quota", "token_used").
			First(&user, "id = ?", userID).Error; err != nil {
			response.InternalError(c, err)
			c.Abort()
			return
		}

		// 检查额度
		if user.TokenUsed >= user.TokenQuota {
			response.Forbidden(c, "您的额度已用完，请联系管理员")
			c.Abort()
			return
		}

		// 将剩余额度存入上下文
		c.Set("tokenRemaining", user.TokenQuota-user.TokenUsed)
		c.Next()
	}
}

// GetTokenRemaining 获取剩余额度
func GetTokenRemaining(c *gin.Context) int64 {
	if remaining, exists := c.Get("tokenRemaining"); exists {
		return remaining.(int64)
	}
	return 0
}
