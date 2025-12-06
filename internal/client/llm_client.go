package client

import (
	"context"
)

// LLMClient LLM 客户端（占位实现）
// TODO: 实现与 OpenAI GPT-4 的交互逻辑
type LLMClient struct {
	apiKey  string
	model   string
	baseURL string
}

// LLMConfig LLM 配置
type LLMConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// NewLLMClient 创建 LLM 客户端
func NewLLMClient(cfg *LLMConfig) *LLMClient {
	return &LLMClient{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: cfg.BaseURL,
	}
}

// ChatCompletion 聊天补全
func (c *LLMClient) ChatCompletion(ctx context.Context, messages []Message) (*ChatResponse, error) {
	// TODO: 实现 OpenAI API 调用
	return nil, nil
}

// ChatCompletionStream 流式聊天补全
func (c *LLMClient) ChatCompletionStream(ctx context.Context, messages []Message) (<-chan StreamChunk, error) {
	// TODO: 实现流式响应
	return nil, nil
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"` // system | user | assistant
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content      string
	FinishReason string
	Usage        TokenUsage
}

// TokenUsage Token 使用量
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// StreamChunk 流式响应块
type StreamChunk struct {
	Content string
	Done    bool
	Error   error
}
