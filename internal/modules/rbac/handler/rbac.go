package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/rbac/dto"
	"github.com/spiringo/spiringo/internal/modules/rbac/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：RBACHandler 定义当前包使用的数据结构或接口。
// English: RBACHandler defines a data structure or interface used by this package.
// RBACHandler 权限管理HTTP处理器
type RBACHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.RBACService
}

// 中文：NewRBACHandler 创建并返回对应组件实例。
// English: NewRBACHandler creates and returns the corresponding component instance.
// NewRBACHandler 创建权限管理处理器
func NewRBACHandler(svc *service.RBACService) *RBACHandler {
	return &RBACHandler{svc: svc}
}

// 中文：ListRoles 执行当前包中的对应流程。
// English: ListRoles executes the corresponding workflow in this package.
func (h *RBACHandler) ListRoles(c *gin.Context) {
	var req types.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	roles, total, err := h.svc.ListRoles(c.Request.Context(), req.GetPage(), req.GetPageSize())
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OKWithPage(c, roles, total, req.GetPage(), req.GetPageSize())
}

// 中文：CreateRole 执行当前包中的对应流程。
// English: CreateRole executes the corresponding workflow in this package.
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	role, err := h.svc.CreateRole(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	c.JSON(201, types.Response{Code: 0, Message: "created", Data: role})
}

// 中文：UpdateRole 执行当前包中的对应流程。
// English: UpdateRole executes the corresponding workflow in this package.
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	role, err := h.svc.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, role)
}

// 中文：DeleteRole 执行当前包中的对应流程。
// English: DeleteRole executes the corresponding workflow in this package.
func (h *RBACHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRole(c.Request.Context(), id); err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, nil)
}

// 中文：ListPermissions 执行当前包中的对应流程。
// English: ListPermissions executes the corresponding workflow in this package.
func (h *RBACHandler) ListPermissions(c *gin.Context) {
	var req types.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	perms, total, err := h.svc.ListPermissions(c.Request.Context(), req.GetPage(), req.GetPageSize())
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OKWithPage(c, perms, total, req.GetPage(), req.GetPageSize())
}

// 中文：AssignPermissions 执行当前包中的对应流程。
// English: AssignPermissions executes the corresponding workflow in this package.
func (h *RBACHandler) AssignPermissions(c *gin.Context) {
	roleID := c.Param("id")
	var req dto.AssignPermissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	if err := h.svc.AssignPermissions(c.Request.Context(), roleID, req.PermissionIDs); err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, nil)
}
