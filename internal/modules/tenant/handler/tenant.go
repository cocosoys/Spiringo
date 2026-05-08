package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/tenant/dto"
	"github.com/spiringo/spiringo/internal/modules/tenant/model"
	"github.com/spiringo/spiringo/internal/modules/tenant/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TenantHandler 定义当前包使用的数据结构或接口。
// English: TenantHandler defines a data structure or interface used by this package.
// TenantHandler 租户HTTP处理器
type TenantHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.TenantService
}

// 中文：NewTenantHandler 创建并返回对应组件实例。
// English: NewTenantHandler creates and returns the corresponding component instance.
// NewTenantHandler 创建租户处理器
func NewTenantHandler(svc *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (h *TenantHandler) List(c *gin.Context) {
	var req types.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	tenants, total, err := h.svc.List(c.Request.Context(), req.GetPage(), req.GetPageSize())
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OKWithPage(c, tenants, total, req.GetPage(), req.GetPageSize())
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (h *TenantHandler) Get(c *gin.Context) {
	id := c.Param("id")
	tenant, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, tenant)
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (h *TenantHandler) Create(c *gin.Context) {
	var req dto.CreateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	m := &model.Tenant{
		Name:     req.Name,
		Code:     req.Code,
		Strategy: req.Strategy,
		Domain:   req.Domain,
	}

	if m.Strategy == "" {
		m.Strategy = "shared_db"
	}

	if err := h.svc.Create(c.Request.Context(), m); err != nil {
		types.Fail(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.Response{
		Code:    0,
		Message: "created",
		Data:    m,
	})
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (h *TenantHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	m, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		types.Fail(c, err)
		return
	}

	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Domain != "" {
		m.Domain = req.Domain
	}
	if req.Status != "" {
		m.Status = req.Status
	}

	if err := h.svc.Update(c.Request.Context(), m); err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, m)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (h *TenantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, nil)
}
