package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/user/dto"
	"github.com/spiringo/spiringo/internal/modules/user/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：UserHandler 定义当前包使用的数据结构或接口。
// English: UserHandler defines a data structure or interface used by this package.
// UserHandler 用户HTTP处理器
type UserHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.UserService
}

// 中文：NewUserHandler 创建并返回对应组件实例。
// English: NewUserHandler creates and returns the corresponding component instance.
// NewUserHandler 创建用户处理器
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (h *UserHandler) List(c *gin.Context) {
	var req types.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	users, total, err := h.svc.List(c.Request.Context(), req.GetPage(), req.GetPageSize())
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OKWithPage(c, users, total, req.GetPage(), req.GetPageSize())
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, user)
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	user, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.Response{
		Code:    0,
		Message: "created",
		Data:    user,
	})
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	user, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, user)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, nil)
}
