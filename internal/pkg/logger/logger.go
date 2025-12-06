// Package logger 提供基于 Zap 的高性能日志系统
package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey 上下文键类型
type contextKey string

const (
	// RequestIDKey 请求 ID 上下文键
	RequestIDKey contextKey = "request_id"
)

var (
	// global 全局 logger 实例
	global *zap.Logger
	once   sync.Once
)

// Config 日志配置
type Config struct {
	Level      string // debug | info | warn | error
	Format     string // json | console
	Output     string // stdout | file
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

// Init 初始化全局日志
func Init(cfg *Config) error {
	var err error
	once.Do(func() {
		global, err = newLogger(cfg)
	})
	return err
}

// newLogger 创建新的日志实例
func newLogger(cfg *Config) (*zap.Logger, error) {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择编码器格式
	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 配置输出
	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "file" && cfg.FilePath != "" {
		file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建 core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建 logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, nil
}

// L 获取全局 logger
func L() *zap.Logger {
	if global == nil {
		// 如果未初始化，返回一个 nop logger
		return zap.NewNop()
	}
	return global
}

// S 获取全局 SugaredLogger
func S() *zap.SugaredLogger {
	return L().Sugar()
}

// WithContext 从上下文中获取带 RequestID 的 logger
func WithContext(ctx context.Context) *zap.Logger {
	logger := L()
	if ctx == nil {
		return logger
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		return logger.With(zap.String("request_id", requestID))
	}

	return logger
}

// WithRequestID 创建带 RequestID 的 logger
func WithRequestID(requestID string) *zap.Logger {
	return L().With(zap.String("request_id", requestID))
}

// Sync 同步日志缓冲
func Sync() error {
	if global != nil {
		return global.Sync()
	}
	return nil
}

// Debug 输出 debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	L().Debug(msg, fields...)
}

// Info 输出 info 级别日志
func Info(msg string, fields ...zap.Field) {
	L().Info(msg, fields...)
}

// Warn 输出 warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	L().Warn(msg, fields...)
}

// Error 输出 error 级别日志
func Error(msg string, fields ...zap.Field) {
	L().Error(msg, fields...)
}

// Fatal 输出 fatal 级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
}
