package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spiringo/spiringo/internal/modules/qrcode/dto"
)

// 中文：QRCodeStyle 定义当前包使用的数据结构或接口。
// English: QRCodeStyle defines a data structure or interface used by this package.
type QRCodeStyle struct {
	// 中文：Size 保存当前结构中的配置或数据值。
	// English: Size stores a configuration or data value for this struct.
	Size int
	// 中文：ForegroundColor 保存当前结构中的配置或数据值。
	// English: ForegroundColor stores a configuration or data value for this struct.
	ForegroundColor string
	// 中文：BackgroundColor 保存当前结构中的配置或数据值。
	// English: BackgroundColor stores a configuration or data value for this struct.
	BackgroundColor string
	// 中文：LogoURL 保存当前结构中的配置或数据值。
	// English: LogoURL stores a configuration or data value for this struct.
	LogoURL string
	// 中文：LogoSize 保存当前结构中的配置或数据值。
	// English: LogoSize stores a configuration or data value for this struct.
	LogoSize int
	// 中文：Level 保存当前结构中的配置或数据值。
	// English: Level stores a configuration or data value for this struct.
	Level string
}

// 中文：QRCodeResult 定义当前包使用的数据结构或接口。
// English: QRCodeResult defines a data structure or interface used by this package.
type QRCodeResult struct {
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content []byte
	// 中文：ImageURL 保存当前结构中的配置或数据值。
	// English: ImageURL stores a configuration or data value for this struct.
	ImageURL string
	// 中文：ImageBytes 保存当前结构中的配置或数据值。
	// English: ImageBytes stores a configuration or data value for this struct.
	ImageBytes []byte
	// 中文：ExpireAt 保存当前结构中的配置或数据值。
	// English: ExpireAt stores a configuration or data value for this struct.
	ExpireAt time.Time
}

// 中文：ShortLinkResult 定义当前包使用的数据结构或接口。
// English: ShortLinkResult defines a data structure or interface used by this package.
type ShortLinkResult struct {
	// 中文：ShortCode 保存当前结构中的配置或数据值。
	// English: ShortCode stores a configuration or data value for this struct.
	ShortCode string
	// 中文：ShortURL 保存当前结构中的配置或数据值。
	// English: ShortURL stores a configuration or data value for this struct.
	ShortURL string
	// 中文：TargetURL 保存当前结构中的配置或数据值。
	// English: TargetURL stores a configuration or data value for this struct.
	TargetURL string
	// 中文：ScanCount 保存当前结构中的配置或数据值。
	// English: ScanCount stores a configuration or data value for this struct.
	ScanCount int64
	// 中文：ExpireAt 保存当前结构中的配置或数据值。
	// English: ExpireAt stores a configuration or data value for this struct.
	ExpireAt time.Time
}

// 中文：Generator 定义当前包使用的数据结构或接口。
// English: Generator defines a data structure or interface used by this package.
type Generator interface {
	// 中文：Generate 声明该接口需要实现的行为。
	// English: Generate declares behavior required by this interface.
	Generate(ctx context.Context, content string, style *QRCodeStyle) (*QRCodeResult, error)
	// 中文：Parse 声明该接口需要实现的行为。
	// English: Parse declares behavior required by this interface.
	Parse(ctx context.Context, imageData []byte) (string, error)
}

// 中文：ShortLinkService 定义当前包使用的数据结构或接口。
// English: ShortLinkService defines a data structure or interface used by this package.
type ShortLinkService interface {
	// 中文：Create 声明该接口需要实现的行为。
	// English: Create declares behavior required by this interface.
	Create(ctx context.Context, targetURL string, expire time.Duration) (*ShortLinkResult, error)
	// 中文：Resolve 声明该接口需要实现的行为。
	// English: Resolve declares behavior required by this interface.
	Resolve(ctx context.Context, shortCode string) (string, error)
	// 中文：Stats 声明该接口需要实现的行为。
	// English: Stats declares behavior required by this interface.
	Stats(ctx context.Context, shortCode string) (*ShortLinkResult, error)
}

// 中文：QRCodeGenerator 定义当前包使用的数据结构或接口。
// English: QRCodeGenerator defines a data structure or interface used by this package.
// QRCodeGenerator adapts QRCodeService to the blueprint Generator contract.
// It falls back to pure image generation when no repository-backed service is
// available.
type QRCodeGenerator struct {
	// 中文：service 保存当前结构中的配置或数据值。
	// English: service stores a configuration or data value for this struct.
	service *QRCodeService
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
var _ Generator = (*QRCodeGenerator)(nil)

// 中文：NewQRCodeGenerator 创建并返回对应组件实例。
// English: NewQRCodeGenerator creates and returns the corresponding component instance.
func NewQRCodeGenerator(service *QRCodeService) *QRCodeGenerator {
	return &QRCodeGenerator{service: service}
}

// 中文：AsGenerator 执行当前包中的对应流程。
// English: AsGenerator executes the corresponding workflow in this package.
func (s *QRCodeService) AsGenerator() Generator {
	return NewQRCodeGenerator(s)
}

// 中文：Generate 执行当前包中的对应流程。
// English: Generate executes the corresponding workflow in this package.
func (g *QRCodeGenerator) Generate(ctx context.Context, content string, style *QRCodeStyle) (*QRCodeResult, error) {
	req := styleToGenerateReq(content, style)
	if g != nil && g.service != nil && g.service.repo != nil {
		resp, err := g.service.Generate(ctx, req)
		if err != nil {
			return nil, err
		}
		return &QRCodeResult{
			Content:    []byte(resp.Content),
			ImageURL:   resp.ImageURL,
			ImageBytes: resp.ImageData,
		}, nil
	}

	size := req.Size
	if size == 0 {
		size = 256
	}
	level := req.Level
	if level == "" {
		level = "medium"
	}
	pngData, err := buildQRCodePNG(ctx, req, size, level)
	if err != nil {
		return nil, err
	}
	return &QRCodeResult{
		Content:    []byte(content),
		ImageBytes: pngData,
	}, nil
}

// 中文：Parse 执行当前包中的对应流程。
// English: Parse executes the corresponding workflow in this package.
func (g *QRCodeGenerator) Parse(ctx context.Context, imageData []byte) (string, error) {
	if g != nil && g.service != nil {
		return g.service.Parse(ctx, imageData)
	}
	return (&QRCodeService{}).Parse(ctx, imageData)
}

// 中文：ShortLinkAdapter 定义当前包使用的数据结构或接口。
// English: ShortLinkAdapter defines a data structure or interface used by this package.
type ShortLinkAdapter struct {
	// 中文：service 保存当前结构中的配置或数据值。
	// English: service stores a configuration or data value for this struct.
	service *QRCodeService
	// 中文：baseURL 保存当前结构中的配置或数据值。
	// English: baseURL stores a configuration or data value for this struct.
	baseURL string
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
var _ ShortLinkService = (*ShortLinkAdapter)(nil)

// 中文：NewShortLinkAdapter 创建并返回对应组件实例。
// English: NewShortLinkAdapter creates and returns the corresponding component instance.
func NewShortLinkAdapter(service *QRCodeService, baseURL string) *ShortLinkAdapter {
	return &ShortLinkAdapter{service: service, baseURL: baseURL}
}

// 中文：AsShortLinkService 执行当前包中的对应流程。
// English: AsShortLinkService executes the corresponding workflow in this package.
func (s *QRCodeService) AsShortLinkService(baseURL string) ShortLinkService {
	return NewShortLinkAdapter(s, baseURL)
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (a *ShortLinkAdapter) Create(ctx context.Context, targetURL string, expire time.Duration) (*ShortLinkResult, error) {
	if a == nil || a.service == nil {
		return nil, fmt.Errorf("qrcode service is required")
	}
	resp, err := a.service.Generate(ctx, dto.GenerateReq{Content: targetURL})
	if err != nil {
		return nil, err
	}
	result := &ShortLinkResult{
		ShortCode: resp.ShortCode,
		ShortURL:  resp.ShortURL,
		TargetURL: resp.Content,
	}
	if a.baseURL != "" && result.ShortCode != "" {
		result.ShortURL = joinBaseURL(a.baseURL, "/qrcode/s/"+result.ShortCode)
	}
	if expire > 0 {
		result.ExpireAt = time.Now().Add(expire)
	}
	return result, nil
}

// 中文：Resolve 执行当前包中的对应流程。
// English: Resolve executes the corresponding workflow in this package.
func (a *ShortLinkAdapter) Resolve(ctx context.Context, shortCode string) (string, error) {
	if a == nil || a.service == nil {
		return "", fmt.Errorf("qrcode service is required")
	}
	return a.service.ResolveShortCode(ctx, shortCode)
}

// 中文：Stats 执行当前包中的对应流程。
// English: Stats executes the corresponding workflow in this package.
func (a *ShortLinkAdapter) Stats(ctx context.Context, shortCode string) (*ShortLinkResult, error) {
	if a == nil || a.service == nil {
		return nil, fmt.Errorf("qrcode service is required")
	}
	stats, err := a.service.GetStats(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	return &ShortLinkResult{
		ShortCode: stats.ShortCode,
		ShortURL:  joinBaseURL(a.baseURL, "/qrcode/s/"+stats.ShortCode),
		TargetURL: stats.Content,
		ScanCount: stats.ScanCount,
	}, nil
}

// 中文：styleToGenerateReq 执行当前包中的对应流程。
// English: styleToGenerateReq executes the corresponding workflow in this package.
func styleToGenerateReq(content string, style *QRCodeStyle) dto.GenerateReq {
	req := dto.GenerateReq{Content: content}
	if style == nil {
		return req
	}
	req.Size = style.Size
	req.Level = style.Level
	req.ForegroundColor = style.ForegroundColor
	req.BackgroundColor = style.BackgroundColor
	req.LogoURL = style.LogoURL
	req.LogoSize = style.LogoSize
	return req
}

// 中文：joinBaseURL 执行当前包中的对应流程。
// English: joinBaseURL executes the corresponding workflow in this package.
func joinBaseURL(baseURL, path string) string {
	if baseURL == "" {
		return path
	}
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}
