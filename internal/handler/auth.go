package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/middleware"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"github.com/lexveritas/lex-veritas-backend/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authSvc service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login 邮箱密码登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	tokenPair, user, err := h.authSvc.LoginByEmail(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Unauthorized(c, "invalid email or password")
		case service.ErrUserDisabled:
			response.Forbidden(c, "account is disabled")
		case service.ErrAccountLocked:
			response.Forbidden(c, "account is locked, please try again later")
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, dto.LoginResponse{
		User:  user,
		Token: tokenPair,
	})
}

// LoginByPhone 手机验证码登录
// POST /api/v1/auth/login/phone
func (h *AuthHandler) LoginByPhone(c *gin.Context) {
	var req dto.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	tokenPair, user, err := h.authSvc.LoginByPhone(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, dto.LoginResponse{
		User:  user,
		Token: tokenPair,
	})
}

// SendCode 发送验证码
// POST /api/v1/auth/send-code
func (h *AuthHandler) SendCode(c *gin.Context) {
	// TODO: 实现发送验证码逻辑
	response.Success(c, gin.H{"message": "verification code sent"})
}

// Register 用户注册
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	user, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExists:
			response.BadRequest(c, "email already registered")
		default:
			// Password validation errors
			if err.Error() == "password too short" || err.Error() == "password too weak" {
				response.BadRequest(c, err.Error())
			} else {
				response.InternalError(c, err)
			}
		}
		return
	}

	response.SuccessWithMessage(c, "registration successful", user)
}

// Refresh 刷新访问令牌
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	tokenPair, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == service.ErrRefreshTokenInvalid {
			response.Unauthorized(c, "invalid or expired refresh token")
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, tokenPair)
}

// Logout 登出
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 {
		token := authHeader[7:] // 去掉 "Bearer "
		h.authSvc.Logout(c.Request.Context(), token)
	}

	response.SuccessWithMessage(c, "logged out successfully", nil)
}

// Me 获取当前用户信息
// GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "not authenticated")
		return
	}

	user, err := h.authSvc.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, user)
}
