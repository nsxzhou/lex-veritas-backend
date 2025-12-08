package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/auth"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
)

// 上下文键
const (
	ContextKeyUserID = "userID"
	ContextKeyRole   = "role"
	ContextKeyClaims = "claims"
)

// TokenValidator Token 验证接口 (由 service.AuthService 实现)
type TokenValidator interface {
	ValidateAccessToken(token string) (*auth.Claims, error)
	IsTokenBlacklisted(ctx context.Context, jti string) bool
}

// JWTAuth JWT 认证中间件
func JWTAuth(validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证头")
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 token
		claims, err := validator.ValidateAccessToken(tokenString)
		if err != nil {
			if err == auth.ErrTokenExpired {
				response.Unauthorized(c, "令牌已过期")
			} else {
				response.Unauthorized(c, "无效的令牌")
			}
			c.Abort()
			return
		}

		// 检查 token 是否在黑名单中
		if validator.IsTokenBlacklisted(c.Request.Context(), claims.ID) {
			response.Unauthorized(c, "令牌已被撤销")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件 (不强制要求登录)
func OptionalAuth(validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := validator.ValidateAccessToken(tokenString)
		if err == nil && !validator.IsTokenBlacklisted(c.Request.Context(), claims.ID) {
			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyRole, claims.Role)
			c.Set(ContextKeyClaims, claims)
		}

		c.Next()
	}
}

// RequireRole 角色验证中间件
func RequireRole(roles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, exists := c.Get(ContextKeyRole)
		if !exists {
			response.Forbidden(c, "上下文中未找到角色信息")
			c.Abort()
			return
		}

		currentRole := model.UserRole(roleStr.(string))
		for _, role := range roles {
			if currentRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}

// ========================= 辅助函数 =========================

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get(ContextKeyUserID); exists {
		return userID.(string)
	}
	return ""
}

// GetRole 从上下文获取角色
func GetRole(c *gin.Context) model.UserRole {
	if role, exists := c.Get(ContextKeyRole); exists {
		return model.UserRole(role.(string))
	}
	return ""
}

// GetClaims 从上下文获取完整 Claims
func GetClaims(c *gin.Context) *auth.Claims {
	if claims, exists := c.Get(ContextKeyClaims); exists {
		return claims.(*auth.Claims)
	}
	return nil
}

// GetUserIDFromContext 从标准 context 获取用户 ID
func GetUserIDFromContext(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return GetUserID(ginCtx)
	}
	return ""
}
