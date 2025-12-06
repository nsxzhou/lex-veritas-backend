// Package nodes 提供 Eino Graph 节点实现
package nodes

import "context"

// LLMNode 大模型节点
// 负责调用 LLM API 生成回答
type LLMNode struct {
	// client LLM 客户端
	// model  模型名称
}

// NewLLMNode 创建大模型节点
func NewLLMNode() *LLMNode {
	return &LLMNode{}
}

// Generate 生成回答
func (l *LLMNode) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: 实现 LLM 调用逻辑
	return "", nil
}

// GenerateStream 流式生成回答
func (l *LLMNode) GenerateStream(ctx context.Context, prompt string) (<-chan string, error) {
	// TODO: 实现流式 LLM 调用逻辑
	return nil, nil
}
