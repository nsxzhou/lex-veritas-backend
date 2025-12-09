package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lexveritas/lex-veritas-backend/internal/dto"
	"github.com/lexveritas/lex-veritas-backend/internal/middleware"
	"github.com/lexveritas/lex-veritas-backend/internal/model"
	"github.com/lexveritas/lex-veritas-backend/internal/pkg/response"
	"github.com/lexveritas/lex-veritas-backend/internal/service"
)

// AdminHandler 管理员功能处理器
type AdminHandler struct {
	userSvc service.UserService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(userSvc service.UserService) *AdminHandler {
	return &AdminHandler{
		userSvc: userSvc,
	}
}

// ListUsers 获取用户列表
// @Summary      获取用户列表
// @Description  分页查询用户列表,支持按邮箱、角色、状态筛选
// @Tags         管理员
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "页码" default(1)
// @Param        pageSize query int false "每页数量" default(10)
// @Param        email query string false "邮箱筛选"
// @Param        role query string false "角色筛选" Enums(user,admin,super_admin)
// @Param        status query string false "状态筛选" Enums(active,inactive,banned)
// @Param        keyword query string false "关键词搜索(邮箱或姓名)"
// @Success      200 {object} response.Response{data=dto.UserListResponse} "获取成功"
// @Failure      401 {object} response.Response "未授权"
// @Failure      403 {object} response.Response "权限不足"
// @Router       /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := h.userSvc.ListUsers(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetUser 获取用户详情
// @Summary      获取用户详情
// @Description  根据用户ID获取详细信息
// @Tags         管理员
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Success      200 {object} response.Response{data=dto.UserResponse} "获取成功"
// @Failure      401 {object} response.Response "未授权"
// @Failure      403 {object} response.Response "权限不足"
// @Failure      404 {object} response.Response "用户不存在"
// @Router       /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	userResp := dto.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Phone:      user.Phone,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Role:       string(user.Role),
		Status:     string(user.Status),
		TokenQuota: user.TokenQuota,
		TokenUsed:  user.TokenUsed,
		LastLogin:  user.LastLoginAt,
		CreatedAt:  user.CreatedAt,
	}

	response.Success(c, userResp)
}

// UpdateStatus 修改用户状态
// @Summary      修改用户状态
// @Description  修改用户的激活状态
// @Tags         管理员
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Param        request body dto.UpdateStatusRequest true "状态更新请求"
// @Success      200 {object} response.Response "修改成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Failure      403 {object} response.Response "权限不足"
// @Failure      404 {object} response.Response "用户不存在"
// @Router       /admin/users/{id}/status [put]
func (h *AdminHandler) UpdateStatus(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if userID == currentUserID {
		response.BadRequest(c, "不能修改自己的状态")
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userSvc.UpdateStatus(c.Request.Context(), userID, req.Status); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "状态修改成功", nil)
}

// UpdateRole 修改用户角色
// @Summary      修改用户角色
// @Description  修改用户角色,需要super_admin权限
// @Tags         管理员
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Param        request body dto.UpdateRoleRequest true "角色更新请求"
// @Success      200 {object} response.Response "修改成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Failure      403 {object} response.Response "权限不足"
// @Failure      404 {object} response.Response "用户不存在"
// @Router       /admin/users/{id}/role [put]
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if userID == currentUserID {
		response.BadRequest(c, "不能修改自己的角色")
		return
	}

	currentRole := middleware.GetRole(c)
	if currentRole != model.RoleSuperAdmin {
		response.Forbidden(c, "只有超级管理员才能修改角色")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userSvc.UpdateRole(c.Request.Context(), userID, req.Role); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "角色修改成功", nil)
}

// AdjustQuota 调整用户Token配额
// @Summary      调整用户Token配额
// @Description  调整指定用户的Token配额
// @Tags         管理员
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Param        request body dto.AdjustQuotaRequest true "配额调整请求"
// @Success      200 {object} response.Response "调整成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Failure      403 {object} response.Response "权限不足"
// @Failure      404 {object} response.Response "用户不存在"
// @Router       /admin/users/{id}/quota [put]
func (h *AdminHandler) AdjustQuota(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	var req dto.AdjustQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	req.UserID = userID

	if err := h.userSvc.AdjustQuota(c.Request.Context(), userID, req.NewQuota); err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "配额调整成功", nil)
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  软删除指定用户
// @Tags         管理员
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Success      200 {object} response.Response "删除成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Failure      403 {object} response.Response "权限不足"
// @Failure      404 {object} response.Response "用户不存在"
// @Router       /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if userID == currentUserID {
		response.BadRequest(c, "不能删除自己")
		return
	}

	if err := h.userSvc.DeleteUser(c.Request.Context(), userID); err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.SuccessWithMessage(c, "用户删除成功", nil)
}
