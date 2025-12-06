// Package nodes 提供 Eino Graph 节点实现
package nodes

import "context"

// Retriever 检索节点
// 负责调用 Milvus 检索相关 Chunks
type Retriever struct {
	// milvusClient Milvus 客户端
	// topK         检索数量
}

// NewRetriever 创建检索节点
func NewRetriever() *Retriever {
	return &Retriever{}
}

// Retrieve 执行向量检索
func (r *Retriever) Retrieve(ctx context.Context, query string) ([]RetrievalResult, error) {
	// TODO: 实现向量检索逻辑
	// 1. 将 query 转换为向量
	// 2. 调用 Milvus 进行检索
	// 3. 返回检索结果
	return nil, nil
}

// RetrievalResult 检索结果
type RetrievalResult struct {
	ChunkID     int64
	Content     string
	Score       float32
	MerkleProof []string
	Metadata    map[string]interface{}
}
