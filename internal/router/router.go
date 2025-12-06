// Package router 提供 HTTP 路由配置
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/middleware"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/cache"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/database"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
)

// Setup 配置并返回 Gin 引擎
func Setup(cfg *config.Config) *gin.Engine {
	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	// 创建引擎
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(&cfg.CORS))
	r.Use(middleware.RateLimit(&cfg.RateLimit))

	// 健康检查端点
	r.GET("/health", healthHandler)
	r.GET("/ready", readyHandler)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证路由
		auth := v1.Group("/auth")
		{
			auth.POST("/login", notImplemented)
			auth.POST("/login/phone", notImplemented)
			auth.POST("/send-code", notImplemented)
			auth.POST("/register", notImplemented)
			auth.GET("/me", notImplemented)
		}

		// 聊天路由
		chat := v1.Group("/chat")
		{
			chat.POST("", notImplemented)
			chat.GET("/sessions", notImplemented)
			chat.GET("/sessions/:id", notImplemented)
			chat.POST("/sessions", notImplemented)
			chat.DELETE("/sessions/:id", notImplemented)
		}

		// 知识库路由
		knowledge := v1.Group("/knowledge")
		{
			knowledge.GET("/documents", notImplemented)
			knowledge.POST("/documents", notImplemented)
			knowledge.DELETE("/documents/:id", notImplemented)
			knowledge.POST("/documents/mint", notImplemented)
		}

		// 统计路由
		stats := v1.Group("/stats")
		{
			stats.GET("/overview", notImplemented)
			stats.GET("/query-volume", notImplemented)
			stats.GET("/recent-activity", notImplemented)
		}

		// 区块链存证路由
		proof := v1.Group("/proof")
		{
			proof.GET("/merkle-tree", notImplemented)
			proof.GET("/recent", notImplemented)
			proof.POST("/verify/:chunkId", notImplemented)
			proof.GET("/stats", notImplemented)
		}

		// 审计日志路由
		audit := v1.Group("/audit")
		{
			audit.GET("/logs", notImplemented)
			audit.POST("/logs/:id/resolve", notImplemented)
		}

		// 用户管理路由
		users := v1.Group("/users")
		{
			users.GET("", notImplemented)
			users.GET("/:id", notImplemented)
			users.PUT("/:id", notImplemented)
			users.DELETE("/:id", notImplemented)
			users.GET("/:id/history", notImplemented)
			users.GET("/:id/token-usage", notImplemented)
		}

		// 系统设置路由
		settings := v1.Group("/settings")
		{
			settings.GET("", notImplemented)
			settings.PUT("", notImplemented)
		}
	}

	return r
}

// healthHandler 健康检查处理器
func healthHandler(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "healthy",
	})
}

// readyHandler 就绪检查处理器（检查所有依赖服务）
func readyHandler(c *gin.Context) {
	checks := make(map[string]string)
	allReady := true

	// 检查数据库
	if err := database.Health(); err != nil {
		checks["database"] = "unhealthy"
		allReady = false
	} else {
		checks["database"] = "healthy"
	}

	// 检查 Redis
	if err := cache.Health(); err != nil {
		checks["redis"] = "unhealthy"
		allReady = false
	} else {
		checks["redis"] = "healthy"
	}

	// TODO: 检查 Milvus

	status := "ready"
	if !allReady {
		status = "not_ready"
	}

	response.Success(c, gin.H{
		"status": status,
		"checks": checks,
	})
}

// notImplemented 未实现接口的占位处理器
func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "接口尚未实现",
	})
}
