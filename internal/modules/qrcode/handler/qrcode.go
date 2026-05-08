package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/qrcode/dto"
	"github.com/spiringo/spiringo/internal/modules/qrcode/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：maxParseImageBytes 声明当前包使用的常量。
// English: maxParseImageBytes declares constants used by this package.
const maxParseImageBytes = 5 << 20

// 中文：QRCodeHandler 定义当前包使用的数据结构或接口。
// English: QRCodeHandler defines a data structure or interface used by this package.
type QRCodeHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.QRCodeService
}

// 中文：NewQRCodeHandler 创建并返回对应组件实例。
// English: NewQRCodeHandler creates and returns the corresponding component instance.
func NewQRCodeHandler(svc *service.QRCodeService) *QRCodeHandler {
	return &QRCodeHandler{svc: svc}
}

// 中文：Generate 执行当前包中的对应流程。
// English: Generate executes the corresponding workflow in this package.
func (h *QRCodeHandler) Generate(c *gin.Context) {
	var req dto.GenerateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	resp, err := h.svc.Generate(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	if c.Query("format") == "image" && len(resp.ImageData) > 0 {
		c.Data(http.StatusOK, "image/png", resp.ImageData)
		return
	}

	c.JSON(http.StatusCreated, types.Response{Code: 0, Message: "created", Data: resp})
}

// 中文：Parse 执行当前包中的对应流程。
// English: Parse executes the corresponding workflow in this package.
func (h *QRCodeHandler) Parse(c *gin.Context) {
	imageData, err := readParseImageData(c)
	if err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	content, err := h.svc.Parse(c.Request.Context(), imageData)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, gin.H{"content": content})
}

// 中文：Redirect 执行当前包中的对应流程。
// English: Redirect executes the corresponding workflow in this package.
func (h *QRCodeHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	targetURL, err := h.svc.ResolveShortCodeWithMeta(c.Request.Context(), code, service.ScanMeta{
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		UserID:    types.GetUserID(c.Request.Context()),
		TenantID:  types.GetTenantID(c.Request.Context()),
	})
	if err != nil {
		types.Fail(c, err)
		return
	}

	if isURL(targetURL) {
		c.Redirect(http.StatusFound, targetURL)
		return
	}

	types.OK(c, gin.H{"content": targetURL})
}

// 中文：Stats 执行当前包中的对应流程。
// English: Stats executes the corresponding workflow in this package.
func (h *QRCodeHandler) Stats(c *gin.Context) {
	code := c.Param("code")
	stats, err := h.svc.GetStats(c.Request.Context(), code)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, stats)
}

// 中文：readParseImageData 执行当前包中的对应流程。
// English: readParseImageData executes the corresponding workflow in this package.
func readParseImageData(c *gin.Context) ([]byte, error) {
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		file, _, err := c.Request.FormFile("image")
		if err != nil {
			return nil, fmt.Errorf("image file required")
		}
		defer file.Close()
		return readLimitedImage(file)
	}

	var req struct {
		ImageURL    string `json:"image_url"`
		ImageBase64 string `json:"image_base64"`
		ImageData   string `json:"image_data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	if req.ImageURL != "" {
		return downloadParseImage(c, req.ImageURL)
	}
	if req.ImageBase64 != "" {
		return decodeParseImageBase64(req.ImageBase64)
	}
	if req.ImageData != "" {
		return decodeParseImageBase64(req.ImageData)
	}
	return nil, fmt.Errorf("image_url, image_base64 or image file required")
}

// 中文：downloadParseImage 执行当前包中的对应流程。
// English: downloadParseImage executes the corresponding workflow in this package.
func downloadParseImage(c *gin.Context, rawURL string) ([]byte, error) {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return nil, fmt.Errorf("image_url must start with http:// or https://")
	}
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create image request: %w", err)
	}
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode < http.StatusOK || httpResp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("download image http %d", httpResp.StatusCode)
	}
	return readLimitedImage(httpResp.Body)
}

// 中文：readLimitedImage 执行当前包中的对应流程。
// English: readLimitedImage executes the corresponding workflow in this package.
func readLimitedImage(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxParseImageBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read image failed: %w", err)
	}
	if len(data) > maxParseImageBytes {
		return nil, fmt.Errorf("image is too large")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("image is empty")
	}
	return data, nil
}

// 中文：decodeParseImageBase64 执行当前包中的对应流程。
// English: decodeParseImageBase64 executes the corresponding workflow in this package.
func decodeParseImageBase64(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("image_base64 is empty")
	}
	if idx := strings.Index(value, ","); idx >= 0 && strings.Contains(value[:idx], "base64") {
		value = value[idx+1:]
	}
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("decode image_base64: %w", err)
	}
	if len(data) > maxParseImageBytes {
		return nil, fmt.Errorf("image is too large")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("image is empty")
	}
	return data, nil
}

// 中文：isURL 执行当前包中的对应流程。
// English: isURL executes the corresponding workflow in this package.
func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
