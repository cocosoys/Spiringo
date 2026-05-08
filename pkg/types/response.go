package types

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 中文：Response 定义当前包使用的数据结构或接口。
// English: Response defines a data structure or interface used by this package.
// Response 统一响应结构
type Response struct {
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code int `json:"code"`
	// 中文：Message 保存当前结构中的配置或数据值。
	// English: Message stores a configuration or data value for this struct.
	Message string `json:"message"`
	// 中文：Data 保存当前结构中的配置或数据值。
	// English: Data stores a configuration or data value for this struct.
	Data any `json:"data,omitempty"`
}

// 中文：PageData 定义当前包使用的数据结构或接口。
// English: PageData defines a data structure or interface used by this package.
// PageData 分页数据结构
type PageData struct {
	// 中文：List 保存当前结构中的配置或数据值。
	// English: List stores a configuration or data value for this struct.
	List any `json:"list"`
	// 中文：Total 保存当前结构中的配置或数据值。
	// English: Total stores a configuration or data value for this struct.
	Total int64 `json:"total"`
	// 中文：Page 保存当前结构中的配置或数据值。
	// English: Page stores a configuration or data value for this struct.
	Page int `json:"page"`
	// 中文：PageSize 保存当前结构中的配置或数据值。
	// English: PageSize stores a configuration or data value for this struct.
	PageSize int `json:"page_size"`
}

// 中文：OK 执行当前包中的对应流程。
// English: OK executes the corresponding workflow in this package.
// OK 成功响应
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// 中文：OKWithMessage 执行当前包中的对应流程。
// English: OKWithMessage executes the corresponding workflow in this package.
// OKWithMessage 带自定义消息的成功响应
func OKWithMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// 中文：OKWithPage 执行当前包中的对应流程。
// English: OKWithPage executes the corresponding workflow in this package.
// OKWithPage 分页成功响应
func OKWithPage(c *gin.Context, list any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// 中文：Fail 执行当前包中的对应流程。
// English: Fail executes the corresponding workflow in this package.
// Fail 失败响应
func Fail(c *gin.Context, err error) {
	if e, ok := err.(*Error); ok {
		c.JSON(e.HTTPStatus, Response{
			Code:    e.Code,
			Message: e.Message,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    ErrInternal.Code,
		Message: err.Error(),
	})
}

// 中文：FailWithCode 执行当前包中的对应流程。
// English: FailWithCode executes the corresponding workflow in this package.
// FailWithCode 带错误码的失败响应
func FailWithCode(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}
