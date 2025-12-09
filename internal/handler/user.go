package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/middleware"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"github.com/lexveritas/lex-veritas-backend/internal/service"
)

// UserHandler 用户自服务处理器
type UserHandler struct {
	userSvc service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

// GetMyQuota 获取我的Token配额
// @Summary      获取当前用户Token配额
// @Description  查看当前用户的Token配额使用情况
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=dto.UsageStats} "获取成功"
// @Failure      401 {object} response.Response "未授权"
// @Router       /users/me/quota [get]
func (h *UserHandler) GetMyQuota(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "未认证")
		return
	}

	usage, err := h.userSvc.GetUsage(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, usage)
}

// UpdateProfile 更新个人资料
// @Summary      更新个人资料
// @Description  更新当前用户的姓名、手机、头像等信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body dto.UpdateProfileRequest true "更新请求"
// @Success      200 {object} response.Response "更新成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Router       /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "未认证")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userSvc.UpdateProfile(c.Request.Context(), userID, &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// ChangePassword 修改密码
// @Summary      修改密码
// @Description  修改当前用户密码
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body dto.ChangePasswordRequest true "修改密码请求"
// @Success      200 {object} response.Response "密码修改成功"
// @Failure      400 {object} response.Response "旧密码错误或请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Router       /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "未认证")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userSvc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err.Error() == "旧密码错误" {
			response.BadRequest(c, "旧密码错误")
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}
