// Package errors 提供统一的错误处理机制
package errors

import (
	"fmt"
	"net/http"
)

// Code 错误码类型
type Code int

// 错误码定义
const (
	// 通用错误 (1xxx)
	CodeSuccess         Code = 0
	CodeUnknown         Code = 1000
	CodeInvalidParam    Code = 1001
	CodeUnauthorized    Code = 1002
	CodeForbidden       Code = 1003
	CodeNotFound        Code = 1004
	CodeTooManyRequests Code = 1005

	// 用户相关错误 (2xxx)
	CodeUserNotFound      Code = 2001
	CodeUserAlreadyExists Code = 2002
	CodeInvalidPassword   Code = 2003
	CodeInvalidToken      Code = 2004
	CodeTokenExpired      Code = 2005

	// 业务错误 (3xxx)
	CodeChatSessionNotFound Code = 3001
	CodeDocumentNotFound    Code = 3002
	CodeDocumentUploadFail  Code = 3003

	// 系统错误 (5xxx)
	CodeInternalError Code = 5000
	CodeDatabaseError Code = 5001
	CodeRedisError    Code = 5002
	CodeMilvusError   Code = 5003
	CodeLLMError      Code = 5004
)

// 错误码信息映射
var codeMessages = map[Code]string{
	CodeSuccess:         "成功",
	CodeUnknown:         "未知错误",
	CodeInvalidParam:    "参数错误",
	CodeUnauthorized:    "未授权",
	CodeForbidden:       "权限不足",
	CodeNotFound:        "资源不存在",
	CodeTooManyRequests: "请求过于频繁",

	CodeUserNotFound:      "用户不存在",
	CodeUserAlreadyExists: "用户已存在",
	CodeInvalidPassword:   "密码错误",
	CodeInvalidToken:      "无效的令牌",
	CodeTokenExpired:      "令牌已过期",

	CodeChatSessionNotFound: "会话不存在",
	CodeDocumentNotFound:    "文档不存在",
	CodeDocumentUploadFail:  "文档上传失败",

	CodeInternalError: "服务器内部错误",
	CodeDatabaseError: "数据库错误",
	CodeRedisError:    "缓存错误",
	CodeMilvusError:   "向量数据库错误",
	CodeLLMError:      "大模型服务错误",
}

// HTTP 状态码映射
var codeHTTPStatus = map[Code]int{
	CodeSuccess:         http.StatusOK,
	CodeUnknown:         http.StatusInternalServerError,
	CodeInvalidParam:    http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,
	CodeTooManyRequests: http.StatusTooManyRequests,

	CodeUserNotFound:      http.StatusNotFound,
	CodeUserAlreadyExists: http.StatusConflict,
	CodeInvalidPassword:   http.StatusUnauthorized,
	CodeInvalidToken:      http.StatusUnauthorized,
	CodeTokenExpired:      http.StatusUnauthorized,

	CodeChatSessionNotFound: http.StatusNotFound,
	CodeDocumentNotFound:    http.StatusNotFound,
	CodeDocumentUploadFail:  http.StatusInternalServerError,

	CodeInternalError: http.StatusInternalServerError,
	CodeDatabaseError: http.StatusInternalServerError,
	CodeRedisError:    http.StatusInternalServerError,
	CodeMilvusError:   http.StatusInternalServerError,
	CodeLLMError:      http.StatusInternalServerError,
}

// AppError 应用错误
type AppError struct {
	Code    Code   // 错误码
	Message string // 错误信息
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus 返回对应的 HTTP 状态码
func (e *AppError) HTTPStatus() int {
	if status, ok := codeHTTPStatus[e.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// New 创建新的应用错误
func New(code Code) *AppError {
	msg := codeMessages[code]
	if msg == "" {
		msg = "未知错误"
	}
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

// NewWithMessage 创建自定义消息的应用错误
func NewWithMessage(code Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装原始错误
func Wrap(code Code, err error) *AppError {
	msg := codeMessages[code]
	if msg == "" {
		msg = "未知错误"
	}
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

// WrapWithMessage 使用自定义消息包装错误
func WrapWithMessage(code Code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// FromError 从 error 转换为 AppError
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return Wrap(CodeUnknown, err)
}
