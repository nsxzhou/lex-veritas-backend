package client

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	milvusClient, err = client.NewClient(ctx, client.Config{
		Address: cfg.Addr(),
	})
	if err != nil {
		return fmt.Errorf("milvus connection failed: %w", err)
	}

	// 验证连接是否真正可用（gRPC 是延迟连接）
	_, err = milvusClient.ListCollections(ctx)
	if err != nil {
		milvusClient.Close()
		milvusClient = nil
		return fmt.Errorf("milvus connection validation failed: %w", err)
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
		return fmt.Errorf("milvus not initialized")
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
		return nil, fmt.Errorf("milvus not initialized")
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
