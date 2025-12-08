// Package main 提供数据库迁移工具
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/database"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
	"go.uber.org/zap"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件路径")
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	_ = logger.Init(&logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	})

	// 连接数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}
	defer database.Close()

	// 执行自动迁移
	logger.Info("开始数据库迁移...")

	db := database.DB()
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}

	logger.Info("数据库迁移完成")
}
