// Package main 提供 LexVeritas 后端服务入口
//
// @title           LexVeritas API
// @version         1.0.0
// @description     法律智能问答系统 API
//
// @contact.name    LexVeritas nsxzhou
// @contact.email   1790146932@qq.com
//
// @host            localhost:8080
// @BasePath        /api/v1
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lexveritas/lex-veritas-backend/internal/client"
	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/cache"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/database"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
	"github.com/lexveritas/lex-veritas-backend/internal/router"
	"go.uber.org/zap"
)

var (
	configPath string
	version    = "dev"
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件路径")
}

func main() {
	flag.Parse()

	// 1. 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := logger.Init(&logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
	}); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("LexVeritas Backend 启动中...",
		zap.String("version", version),
		zap.String("mode", cfg.Server.Mode),
	)

	// 3. 初始化数据库（可选，允许启动时不连接）
	if err := database.Init(&cfg.Database); err != nil {
		logger.Warn("数据库连接失败（后续请求可能受影响）",
			zap.Error(err),
		)
	}

	// 4. 初始化 Redis（可选，允许启动时不连接）
	if err := cache.Init(&cfg.Redis); err != nil {
		logger.Warn("Redis 连接失败（后续请求可能受影响）",
			zap.Error(err),
		)
	}

	// 5. 初始化 Milvus（可选，允许启动时不连接）
	if err := client.InitMilvus(&cfg.Milvus); err != nil {
		logger.Warn("Milvus 连接失败（后续请求可能受影响）",
			zap.Error(err),
		)
	}

	// 6. 设置路由
	r := router.Setup(cfg)

	// 7. 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 8. 启动服务器（异步）
	go func() {
		baseURL := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
		logger.Info("HTTP 服务器启动",
			zap.Int("port", cfg.Server.Port),
			zap.String("url", baseURL),
		)
		logger.Info("API 文档地址",
			zap.String("scalar", baseURL+"/docs"),
			zap.String("swagger", baseURL+"/swagger/index.html"),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 10. 优雅关闭（等待处理中的请求完成）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	// 11. 关闭资源连接
	if err := database.Close(); err != nil {
		logger.Error("关闭数据库连接失败", zap.Error(err))
	}

	if err := cache.Close(); err != nil {
		logger.Error("关闭 Redis 连接失败", zap.Error(err))
	}

	if err := client.CloseMilvus(); err != nil {
		logger.Error("关闭 Milvus 连接失败", zap.Error(err))
	}

	logger.Info("服务器已安全关闭")
}
