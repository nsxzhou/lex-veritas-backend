package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	var initErr error
	once.Do(func() {
		initErr = connect(cfg)
	})
	return initErr
}

// connect 建立数据库连接
func connect(cfg *config.DatabaseConfig) error {
	// 配置 GORM 日志
	gormLogger := gormlogger.Default.LogMode(gormlogger.Silent)

	// 连接数据库
	var err error
	db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	logger.Info("数据库连接成功",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("dbname", cfg.DBName),
	)

	return nil
}

// DB 获取数据库实例
func DB() *gorm.DB {
	return db
}

// WithContext 获取带上下文的数据库实例
func WithContext(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

// Health 健康检查
func Health() error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	logger.Info("关闭数据库连接")
	return sqlDB.Close()
}

// Transaction 执行事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}
