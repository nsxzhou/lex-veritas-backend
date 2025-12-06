package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/errors"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"go.uber.org/zap"
)

// Recovery Panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求 ID
				requestID := ""
				if id, exists := c.Get("request_id"); exists {
					requestID, _ = id.(string)
				}

				// 记录 panic 日志
				logger.WithRequestID(requestID).Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:      int(errors.CodeInternalError),
					Message:   "服务器内部错误",
					RequestID: requestID,
				})
			}
		}()

		c.Next()
	}
}
