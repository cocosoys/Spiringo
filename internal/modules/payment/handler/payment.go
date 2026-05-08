package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/payment/dto"
	"github.com/spiringo/spiringo/internal/modules/payment/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：PaymentHandler 定义当前包使用的数据结构或接口。
// English: PaymentHandler defines a data structure or interface used by this package.
// PaymentHandler 支付HTTP处理器
type PaymentHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.PaymentService
}

// 中文：NewPaymentHandler 创建并返回对应组件实例。
// English: NewPaymentHandler creates and returns the corresponding component instance.
// NewPaymentHandler 创建支付处理器
func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (h *PaymentHandler) Create(c *gin.Context) {
	var req dto.CreatePayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	resp, err := h.svc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	c.JSON(201, types.Response{Code: 0, Message: "created", Data: resp})
}

// 中文：Callback 执行当前包中的对应流程。
// English: Callback executes the corresponding workflow in this package.
func (h *PaymentHandler) Callback(c *gin.Context) {
	channel := c.Param("channel")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.writeCallback(c, channel, http.StatusBadRequest, false)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	if err := h.svc.HandleCallbackRequest(c.Request.Context(), channel, c.Request, body); err != nil {
		h.writeCallback(c, channel, http.StatusInternalServerError, false)
		return
	}

	h.writeCallback(c, channel, http.StatusOK, true)
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (h *PaymentHandler) Refund(c *gin.Context) {
	var req dto.RefundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	refund, err := h.svc.Refund(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, refund)
}

// 中文：Query 执行当前包中的对应流程。
// English: Query executes the corresponding workflow in this package.
func (h *PaymentHandler) Query(c *gin.Context) {
	id := c.Param("id")
	order, err := h.svc.QueryOrder(c.Request.Context(), id)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, order)
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (h *PaymentHandler) Close(c *gin.Context) {
	id := c.Param("id")
	order, err := h.svc.CloseOrder(c.Request.Context(), id)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, order)
}

// 中文：writeCallback 执行当前包中的对应流程。
// English: writeCallback executes the corresponding workflow in this package.
func (h *PaymentHandler) writeCallback(c *gin.Context, channel string, status int, success bool) {
	payload := h.svc.CallbackResponse(channel, success)
	switch v := payload.(type) {
	case nil:
		c.Status(status)
	case string:
		c.String(status, v)
	case []byte:
		c.Data(status, "application/octet-stream", v)
	default:
		c.JSON(status, v)
	}
}
