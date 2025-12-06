package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	// TODO: 注入依赖
}

// NewChatHandler 创建聊天处理器
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

// Chat 处理聊天请求（SSE 流式响应）
func (h *ChatHandler) Chat(c *gin.Context) {
	// TODO: 实现 SSE 流式聊天响应
	response.Success(c, gin.H{
		"message": "聊天接口尚未实现",
	})
}

// GetSessions 获取会话列表
func (h *ChatHandler) GetSessions(c *gin.Context) {
	// TODO: 实现获取会话列表
	response.Success(c, []interface{}{})
}

// GetSession 获取会话详情
func (h *ChatHandler) GetSession(c *gin.Context) {
	// TODO: 实现获取会话详情
	response.Success(c, nil)
}

// CreateSession 创建新会话
func (h *ChatHandler) CreateSession(c *gin.Context) {
	// TODO: 实现创建会话
	response.Success(c, nil)
}

// DeleteSession 删除会话
func (h *ChatHandler) DeleteSession(c *gin.Context) {
	// TODO: 实现删除会话
	response.Success(c, nil)
}
