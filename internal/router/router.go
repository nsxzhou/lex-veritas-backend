package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/config"
	"github.com/lexveritas/lex-veritas-backend/internal/handler"
	"github.com/lexveritas/lex-veritas-backend/internal/middleware"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/auth"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/cache"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/database"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/email"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"github.com/lexveritas/lex-veritas-backend/internal/service"
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

	// 初始化认证服务
	authSvc := service.NewAuthService(
		&auth.JWTConfig{
			Secret:        cfg.JWT.Secret,
			AccessExpire:  cfg.JWT.AccessExpire,
			RefreshExpire: cfg.JWT.RefreshExpire,
			Issuer:        cfg.JWT.Issuer,
		},
		&auth.PasswordConfig{
			BcryptCost: cfg.Auth.BcryptCost,
		},
	)

	// 初始化邮件发送器
	emailSender := email.NewSMTPSender(&cfg.Email)

	// 初始化验证码服务
	verifySvc := service.NewVerificationService(emailSender, &cfg.Verification)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc, verifySvc)

	// 健康检查端点
	r.GET("/health", healthHandler)
	r.GET("/ready", readyHandler)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// ======== 认证路由 (无需登录) ========
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/login/phone", authHandler.LoginByPhone)
			authRoutes.POST("/send-code", authHandler.SendCode)
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/refresh", authHandler.Refresh)
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
