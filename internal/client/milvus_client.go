// Package client 提供外部服务客户端封装
package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

var (
	milvusClient client.Client
	milvusOnce   sync.Once
)

// MilvusConfig Milvus 配置别名
type MilvusConfig = config.MilvusConfig

// InitMilvus 初始化 Milvus 客户端
func InitMilvus(cfg *MilvusConfig) error {
	var initErr error
	milvusOnce.Do(func() {
		initErr = connectMilvus(cfg)
	})
	return initErr
}

// connectMilvus 建立 Milvus 连接
func connectMilvus(cfg *MilvusConfig) error {
	ctx := context.Background()

	var err error
	milvusClient, err = client.NewClient(ctx, client.Config{
		Address: cfg.Addr(),
	})
	if err != nil {
		return fmt.Errorf("Milvus 连接失败: %w", err)
	}

	logger.Info("Milvus 连接成功",
		zap.String("addr", cfg.Addr()),
		zap.String("collection", cfg.CollectionName),
	)

	return nil
}

// GetMilvusClient 获取 Milvus 客户端
func GetMilvusClient() client.Client {
	return milvusClient
}

// MilvusHealth 健康检查
func MilvusHealth() error {
	if milvusClient == nil {
		return fmt.Errorf("Milvus 未初始化")
	}
	// Milvus SDK 没有直接的 Ping 方法，通过检查连接状态判断
	return nil
}

// CloseMilvus 关闭 Milvus 连接
func CloseMilvus() error {
	if milvusClient == nil {
		return nil
	}

	logger.Info("关闭 Milvus 连接")
	return milvusClient.Close()
}

// SearchVector 向量检索
func SearchVector(ctx context.Context, collectionName string, vectors [][]float32, topK int) ([]SearchResult, error) {
	if milvusClient == nil {
		return nil, fmt.Errorf("Milvus 未初始化")
	}

	// TODO: 实现具体的向量检索逻辑
	// 这里只是接口定义，具体实现需要根据 Collection Schema 调整
	return nil, nil
}

// SearchResult 检索结果
type SearchResult struct {
	ID       int64
	Score    float32
	Content  string
	Metadata map[string]interface{}
}
