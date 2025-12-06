// Package response 提供统一的 HTTP 响应格式
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/errors"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CursorPageData 游标分页数据结构
type CursorPageData struct {
	List       interface{} `json:"list"`
	NextCursor string      `json:"nextCursor,omitempty"` // 下一页游标，为空表示没有更多数据
	HasMore    bool        `json:"hasMore"`              // 是否有更多数据
	Limit      int         `json:"limit"`                // 每页数量
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.CodeSuccess),
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.CodeSuccess),
		Message: message,
		Data:    data,
	})
}

// SuccessCursorPage 游标分页成功响应
func SuccessCursorPage(c *gin.Context, list interface{}, nextCursor string, hasMore bool, limit int) {
	Success(c, CursorPageData{
		List:       list,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      limit,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	requestID := getRequestID(c)
	appErr := errors.FromError(err)

	// 记录错误日志（内部使用 requestID）
	logger.WithRequestID(requestID).Error(appErr.Error())

	c.JSON(appErr.HTTPStatus(), Response{
		Code:    int(appErr.Code),
		Message: appErr.Message,
	})
}

// ErrorWithCode 指定错误码响应
func ErrorWithCode(c *gin.Context, code errors.Code) {
	Error(c, errors.New(code))
}

// ErrorWithMessage 指定错误码和消息响应
func ErrorWithMessage(c *gin.Context, code errors.Code, message string) {
	Error(c, errors.NewWithMessage(code, message))
}

// BadRequest 参数错误响应
func BadRequest(c *gin.Context, message string) {
	ErrorWithMessage(c, errors.CodeInvalidParam, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context) {
	ErrorWithCode(c, errors.CodeUnauthorized)
}

// Forbidden 权限不足响应
func Forbidden(c *gin.Context) {
	ErrorWithCode(c, errors.CodeForbidden)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context) {
	ErrorWithCode(c, errors.CodeNotFound)
}

// InternalError 服务器内部错误响应
func InternalError(c *gin.Context, err error) {
	Error(c, errors.Wrap(errors.CodeInternalError, err))
}

// getRequestID 从上下文获取请求 ID（仅用于日志）
func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}
