package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/notification/dto"
	"github.com/spiringo/spiringo/internal/modules/notification/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：NotificationHandler 定义当前包使用的数据结构或接口。
// English: NotificationHandler defines a data structure or interface used by this package.
type NotificationHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.Service
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
func New(svc *service.Service) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// 中文：Send 执行当前包中的对应流程。
// English: Send executes the corresponding workflow in this package.
func (h *NotificationHandler) Send(c *gin.Context) {
	var req dto.SendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}
	severity := req.Severity
	if severity == "" {
		severity = "info"
	}

	err := h.svc.Notify(c.Request.Context(), service.Message{
		Event:       req.Event,
		Severity:    severity,
		Subject:     req.Subject,
		Content:     req.Content,
		TenantID:    req.TenantID,
		RecipientID: req.RecipientID,
		Payload:     req.Payload,
	})
	if err != nil {
		types.Fail(c, types.ErrInternal.WithMessage(err.Error()))
		return
	}
	types.OK(c, gin.H{"sent": true})
}

// 中文：Inbox 执行当前包中的对应流程。
// English: Inbox executes the corresponding workflow in this package.
func (h *NotificationHandler) Inbox(c *gin.Context) {
	var req dto.InboxListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	items, total, err := h.svc.ListInbox(c.Request.Context(), service.InboxFilter{
		Page:        req.GetPage(),
		PageSize:    req.GetPageSize(),
		Event:       req.Event,
		RecipientID: req.RecipientID,
		UnreadOnly:  req.UnreadOnly,
	})
	if err != nil {
		types.Fail(c, types.ErrInternal.WithMessage(err.Error()))
		return
	}
	types.OKWithPage(c, items, total, req.GetPage(), req.GetPageSize())
}

// 中文：MarkRead 执行当前包中的对应流程。
// English: MarkRead executes the corresponding workflow in this package.
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	if err := h.svc.MarkRead(c.Request.Context(), c.Param("id"), c.Query("recipient_id")); err != nil {
		types.Fail(c, types.ErrInternal.WithMessage(err.Error()))
		return
	}
	types.OK(c, gin.H{"read": true})
}
