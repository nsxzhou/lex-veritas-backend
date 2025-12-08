package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/middleware"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/auth"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"github.com/lexveritas/lex-veritas-backend/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authSvc   service.AuthService
	verifySvc service.VerificationService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authSvc service.AuthService, verifySvc service.VerificationService) *AuthHandler {
	return &AuthHandler{
		authSvc:   authSvc,
		verifySvc: verifySvc,
	}
}

// Login 邮箱密码登录
// @Summary      用户登录
// @Description  使用邮箱和密码登录系统
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "登录请求"
// @Success      200 {object} response.Response{data=dto.LoginResponse} "登录成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "邮箱或密码错误"
// @Failure      403 {object} response.Response "账户被禁用或锁定"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	tokenPair, user, err := h.authSvc.LoginByEmail(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Unauthorized(c, "邮箱或密码错误")
		case service.ErrUserDisabled:
			response.Forbidden(c, "账户已被禁用")
		case service.ErrAccountLocked:
			response.Forbidden(c, "账户已被锁定，请稍后再试")
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
// @Summary      手机验证码登录
// @Description  使用手机号和验证码登录系统
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body dto.PhoneLoginRequest true "手机登录请求"
// @Success      200 {object} response.Response{data=dto.LoginResponse} "登录成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "验证码错误"
// @Router       /auth/login/phone [post]
func (h *AuthHandler) LoginByPhone(c *gin.Context) {
	var req dto.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
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
// @Summary      发送验证码
// @Description  发送邮箱验证码，用于注册或重置密码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body dto.SendCodeRequest true "发送验证码请求"
// @Success      200 {object} response.Response "验证码发送成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      429 {object} response.Response "请求过于频繁"
// @Router       /auth/send-code [post]
func (h *AuthHandler) SendCode(c *gin.Context) {
	var req dto.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.verifySvc.SendCode(c.Request.Context(), req.Email, req.Purpose); err != nil {
		switch err {
		case service.ErrCodeSendTooFrequent:
			response.TooManyRequests(c, err.Error())
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, gin.H{"message": "验证码已发送"})
}

// Register 用户注册
// @Summary      用户注册
// @Description  使用邮箱注册新用户，需先获取验证码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "注册请求"
// @Success      200 {object} response.Response{data=dto.UserResponse} "注册成功"
// @Failure      400 {object} response.Response "请求参数错误或验证码无效"
// @Failure      429 {object} response.Response "验证码尝试次数过多"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if err := h.verifySvc.VerifyCode(c.Request.Context(), req.Email, req.Code, "register"); err != nil {
		switch err {
		case service.ErrCodeInvalid, service.ErrCodeExpired:
			response.BadRequest(c, err.Error())
		case service.ErrTooManyAttempts:
			response.TooManyRequests(c, err.Error())
		default:
			response.InternalError(c, err)
		}
		return
	}

	user, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExists:
			response.BadRequest(c, "邮箱已被注册")
		case auth.ErrPasswordTooShort, auth.ErrPasswordTooWeak1, auth.ErrPasswordTooWeak2:
			response.BadRequest(c, err.Error())
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.SuccessWithMessage(c, "注册成功", user)
}

// Refresh 刷新访问令牌
// @Summary      刷新令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "刷新令牌请求"
// @Success      200 {object} response.Response{data=dto.TokenPair} "刷新成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "刷新令牌无效或已过期"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	tokenPair, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == service.ErrRefreshTokenInvalid {
			response.Unauthorized(c, "刷新令牌无效或已过期")
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, tokenPair)
}

// Logout 登出
// @Summary      用户登出
// @Description  登出当前用户，使访问令牌失效
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response "登出成功"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 {
		token := authHeader[7:] // 去掉 "Bearer "
		h.authSvc.Logout(c.Request.Context(), token)
	}

	response.SuccessWithMessage(c, "登出成功", nil)
}

// Me 获取当前用户信息
// @Summary      获取当前用户
// @Description  获取当前登录用户的详细信息
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=dto.UserResponse} "获取成功"
// @Failure      401 {object} response.Response "未授权"
// @Failure      404 {object} response.Response "用户不存在"
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "未认证")
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
