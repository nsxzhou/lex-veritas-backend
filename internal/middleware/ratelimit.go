package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/errors"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"golang.org/x/time/rate"
)

// ipLimiters IP 限流器映射
var (
	ipLimiters = make(map[string]*rate.Limiter)
	mu         sync.RWMutex
)

// RateLimit 限流中间件（基于令牌桶算法）
func RateLimit(cfg *config.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip, cfg.Rate, cfg.Burst)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Code:    int(errors.CodeTooManyRequests),
				Message: "请求过于频繁，请稍后重试",
			})
			return
		}

		c.Next()
	}
}

// getLimiter 获取或创建 IP 限流器
func getLimiter(ip string, rateLimit, burst int) *rate.Limiter {
	mu.RLock()
	limiter, exists := ipLimiters[ip]
	mu.RUnlock()

	if exists {
		return limiter
	}

	mu.Lock()
	defer mu.Unlock()

	// 双重检查
	if limiter, exists = ipLimiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rate.Limit(rateLimit), burst)
	ipLimiters[ip] = limiter

	return limiter
}

// CleanupLimiters 清理过期的限流器（可定期调用）
func CleanupLimiters() {
	mu.Lock()
	defer mu.Unlock()
	ipLimiters = make(map[string]*rate.Limiter)
}
