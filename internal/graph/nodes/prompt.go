// Package nodes 提供 Eino Graph 节点实现
package nodes

// PromptBuilder 提示词节点
// 负责构建包含验证通过信息的 Prompt
type PromptBuilder struct {
	// systemPrompt 系统提示词模板
}

// NewPromptBuilder 创建提示词节点
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// Build 构建提示词
func (p *PromptBuilder) Build(input PromptInput) string {
	// TODO: 实现提示词构建逻辑
	// 1. 组装检索到的法条内容
	// 2. 添加验证状态信息
	// 3. 构建完整的 Prompt
	return ""
}

// PromptInput 提示词构建输入
type PromptInput struct {
	Question    string
	Chunks      []ChunkWithVerification
	ChatHistory []ChatHistoryItem
}

// ChunkWithVerification 带验证信息的 Chunk
type ChunkWithVerification struct {
	Content    string
	Source     string
	Verified   bool
	VerifyInfo VerifyResult
}

// ChatHistoryItem 聊天历史项
type ChatHistoryItem struct {
	Role    string
	Content string
}
