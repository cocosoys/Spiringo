package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/makiuchi-d/gozxing"
	goqrcode "github.com/makiuchi-d/gozxing/qrcode"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spiringo/spiringo/internal/modules/qrcode/dto"
	"github.com/spiringo/spiringo/internal/modules/qrcode/model"
	"github.com/spiringo/spiringo/internal/modules/qrcode/repository"
	"github.com/spiringo/spiringo/internal/pkg/storage"
	"github.com/spiringo/spiringo/internal/pkg/utils"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 二维码服务配置
type Config struct {
	// 中文：DefaultSize 保存当前结构中的配置或数据值。
	// English: DefaultSize stores a configuration or data value for this struct.
	DefaultSize int
	// 中文：DefaultLevel 保存当前结构中的配置或数据值。
	// English: DefaultLevel stores a configuration or data value for this struct.
	DefaultLevel string
	// 中文：OSSPrefix 保存当前结构中的配置或数据值。
	// English: OSSPrefix stores a configuration or data value for this struct.
	OSSPrefix string
	// 中文：BucketName 保存当前结构中的配置或数据值。
	// English: BucketName stores a configuration or data value for this struct.
	BucketName string // OSS bucket名称
}

// 中文：QRCodeService 定义当前包使用的数据结构或接口。
// English: QRCodeService defines a data structure or interface used by this package.
// QRCodeService 二维码业务逻辑
type QRCodeService struct {
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：repo 保存当前结构中的配置或数据值。
	// English: repo stores a configuration or data value for this struct.
	repo *repository.QRCodeRepository
	// 中文：storage 保存当前结构中的配置或数据值。
	// English: storage stores a configuration or data value for this struct.
	storage storage.Storage // 可选：为nil时不使用OSS
}

// 中文：ScanMeta 定义当前包使用的数据结构或接口。
// English: ScanMeta defines a data structure or interface used by this package.
type ScanMeta struct {
	// 中文：IP 保存当前结构中的配置或数据值。
	// English: IP stores a configuration or data value for this struct.
	IP string
	// 中文：UserAgent 保存当前结构中的配置或数据值。
	// English: UserAgent stores a configuration or data value for this struct.
	UserAgent string
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string
}

// 中文：NewQRCodeService 创建并返回对应组件实例。
// English: NewQRCodeService creates and returns the corresponding component instance.
// NewQRCodeService 创建二维码服务
func NewQRCodeService(config Config, repo *repository.QRCodeRepository, s storage.Storage) *QRCodeService {
	return &QRCodeService{config: config, repo: repo, storage: s}
}

// 中文：parseLevel 执行当前包中的对应流程。
// English: parseLevel executes the corresponding workflow in this package.
// parseLevel 将级别字符串转为 qrcode.RecoveryLevel
func parseLevel(level string) qrcode.RecoveryLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "l", "low":
		return qrcode.Low
	case "m", "medium":
		return qrcode.Medium
	case "q", "quartile", "high":
		return qrcode.High
	case "h", "highest":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// 中文：Generate 执行当前包中的对应流程。
// English: Generate executes the corresponding workflow in this package.
// Generate 生成二维码
func (s *QRCodeService) Generate(ctx context.Context, req dto.GenerateReq) (*dto.GenerateResp, error) {
	size := req.Size
	if size == 0 {
		size = s.config.DefaultSize
	}
	if size == 0 {
		size = 256
	}
	level := req.Level
	if level == "" {
		level = s.config.DefaultLevel
	}
	if level == "" {
		level = "medium"
	}

	// 先生成图片，避免颜色、Logo 等输入无效时留下不可用记录。
	pngData, err := buildQRCodePNG(ctx, req, size, level)
	if err != nil {
		return nil, types.ErrBadRequest.WithMessagef("generate qr code image: %v", err)
	}

	shortCode := utils.GenerateRandomString(8)

	record := &model.QRCodeRecord{
		Content:         req.Content,
		ShortCode:       shortCode,
		Size:            size,
		Level:           level,
		ForegroundColor: req.ForegroundColor,
		BackgroundColor: req.BackgroundColor,
		LogoURL:         req.LogoURL,
		LogoSize:        req.LogoSize,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("create qrcode record: %w", err)
	}

	// 上传到OSS
	imageURL := ""
	if s.storage != nil && s.config.BucketName != "" {
		key := qrcodeObjectKey(s.config.OSSPrefix, shortCode)
		if err := s.storage.Upload(ctx, s.config.BucketName, key, pngData, "image/png"); err != nil {
			slog.Warn("failed to upload qrcode to OSS", "key", key, "error", err)
		} else {
			imageURL = s.storage.GetURL(s.config.BucketName, key)
		}
	} else if s.config.OSSPrefix != "" {
		// 仅拼接URL（无Storage客户端时的降级方案）
		imageURL = qrcodeFallbackURL(s.config.OSSPrefix, shortCode)
	}

	// 更新记录的图片URL
	if imageURL != "" {
		record.ImageURL = imageURL
		if err := s.repo.Update(ctx, record); err != nil {
			slog.Warn("failed to update qrcode image URL", "error", err)
		}
	}

	resp := &dto.GenerateResp{
		ID:        record.ID,
		Content:   record.Content,
		ShortCode: record.ShortCode,
		ShortURL:  fmt.Sprintf("/qrcode/s/%s", shortCode),
		ImageURL:  imageURL,
		ImageData: pngData,
	}

	return resp, nil
}

// 中文：Parse 执行当前包中的对应流程。
// English: Parse executes the corresponding workflow in this package.
// Parse 解析二维码图片
func (s *QRCodeService) Parse(ctx context.Context, imageData []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", types.ErrBadRequest.WithMessagef("decode image: %v", err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", types.ErrBadRequest.WithMessagef("process image: %v", err)
	}

	reader := goqrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		return "", types.ErrBadRequest.WithMessagef("parse qr code: %v", err)
	}

	return result.String(), nil
}

// 中文：ResolveShortCode 执行当前包中的对应流程。
// English: ResolveShortCode executes the corresponding workflow in this package.
// ResolveShortCode 解析短码
func (s *QRCodeService) ResolveShortCode(ctx context.Context, code string) (string, error) {
	return s.ResolveShortCodeWithMeta(ctx, code, ScanMeta{})
}

// 中文：ResolveShortCodeWithMeta 执行当前包中的对应流程。
// English: ResolveShortCodeWithMeta executes the corresponding workflow in this package.
func (s *QRCodeService) ResolveShortCodeWithMeta(ctx context.Context, code string, meta ScanMeta) (string, error) {
	record, err := s.repo.GetByShortCode(ctx, code)
	if err != nil {
		return "", types.ErrNotFound.WithMessage("short code not found")
	}

	// 增加扫码计数
	if err := s.repo.IncrScanCount(ctx, code); err != nil {
		slog.Warn("failed to increment scan count", "code", code, "error", err)
	}

	// 记录扫码日志
	if err := s.repo.CreateScanLog(ctx, &model.ScanLog{
		ShortCode: code,
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
		UserID:    meta.UserID,
		TenantID:  meta.TenantID,
	}); err != nil {
		slog.Warn("failed to create scan log", "code", code, "error", err)
	}

	return record.Content, nil
}

// 中文：GetStats 执行当前包中的对应流程。
// English: GetStats executes the corresponding workflow in this package.
// GetStats 获取统计
func (s *QRCodeService) GetStats(ctx context.Context, code string) (*dto.StatsResp, error) {
	record, err := s.repo.GetByShortCode(ctx, code)
	if err != nil {
		return nil, types.ErrNotFound.WithMessage("short code not found")
	}

	return &dto.StatsResp{
		ShortCode: record.ShortCode,
		ScanCount: record.ScanCount,
		Content:   record.Content,
	}, nil
}

// 中文：qrcodeObjectKey 执行当前包中的对应流程。
// English: qrcodeObjectKey executes the corresponding workflow in this package.
func qrcodeObjectKey(prefix, shortCode string) string {
	prefix = strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "" {
		prefix = "qrcode"
	}
	return prefix + "/" + shortCode + ".png"
}

// 中文：qrcodeFallbackURL 执行当前包中的对应流程。
// English: qrcodeFallbackURL executes the corresponding workflow in this package.
func qrcodeFallbackURL(prefix, shortCode string) string {
	prefix = strings.TrimRight(strings.TrimSpace(prefix), "/")
	if prefix == "" {
		return ""
	}
	return prefix + "/" + shortCode + ".png"
}

// 中文：buildQRCodePNG 执行当前包中的对应流程。
// English: buildQRCodePNG executes the corresponding workflow in this package.
func buildQRCodePNG(ctx context.Context, req dto.GenerateReq, size int, level string) ([]byte, error) {
	q, err := qrcode.New(req.Content, parseLevel(level))
	if err != nil {
		return nil, err
	}

	background := color.Color(color.White)
	if fg, ok, err := parseHexColor(req.ForegroundColor); err != nil {
		return nil, err
	} else if ok {
		q.ForegroundColor = fg
	}
	if bg, ok, err := parseHexColor(req.BackgroundColor); err != nil {
		return nil, err
	} else if ok {
		q.BackgroundColor = bg
		background = bg
	}

	img := imageToRGBA(q.Image(size))
	if strings.TrimSpace(req.LogoURL) != "" {
		logo, err := fetchLogoImage(ctx, req.LogoURL)
		if err != nil {
			return nil, err
		}
		drawCenteredLogo(img, logo, chooseLogoSize(size, req.LogoSize), background)
	}

	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// 中文：parseHexColor 执行当前包中的对应流程。
// English: parseHexColor executes the corresponding workflow in this package.
func parseHexColor(value string) (color.Color, bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false, nil
	}
	value = strings.TrimPrefix(value, "#")
	if len(value) == 3 {
		value = strings.Repeat(value[0:1], 2) + strings.Repeat(value[1:2], 2) + strings.Repeat(value[2:3], 2)
	}
	if len(value) != 6 {
		return nil, false, fmt.Errorf("invalid hex color %q", value)
	}
	n, err := strconv.ParseUint(value, 16, 32)
	if err != nil {
		return nil, false, fmt.Errorf("invalid hex color %q: %w", value, err)
	}
	return color.RGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 0xff}, true, nil
}

// 中文：fetchLogoImage 执行当前包中的对应流程。
// English: fetchLogoImage executes the corresponding workflow in this package.
func fetchLogoImage(ctx context.Context, rawURL string) (image.Image, error) {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return nil, fmt.Errorf("logo_url must start with http:// or https://")
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create logo request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch logo: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("fetch logo http %d", resp.StatusCode)
	}
	img, _, err := image.Decode(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("decode logo image: %w", err)
	}
	return img, nil
}

// 中文：imageToRGBA 执行当前包中的对应流程。
// English: imageToRGBA executes the corresponding workflow in this package.
func imageToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

// 中文：chooseLogoSize 执行当前包中的对应流程。
// English: chooseLogoSize executes the corresponding workflow in this package.
func chooseLogoSize(qrSize, requested int) int {
	if qrSize <= 0 {
		qrSize = 256
	}
	maxLogo := qrSize / 3
	if maxLogo < 16 {
		maxLogo = 16
	}
	if requested <= 0 {
		requested = qrSize / 5
	}
	if requested < 16 {
		return 16
	}
	if requested > maxLogo {
		return maxLogo
	}
	return requested
}

// 中文：drawCenteredLogo 执行当前包中的对应流程。
// English: drawCenteredLogo executes the corresponding workflow in this package.
func drawCenteredLogo(dst *image.RGBA, logo image.Image, size int, background color.Color) {
	if dst == nil || logo == nil || size <= 0 {
		return
	}
	bounds := dst.Bounds()
	x0 := bounds.Min.X + (bounds.Dx()-size)/2
	y0 := bounds.Min.Y + (bounds.Dy()-size)/2
	logoRect := image.Rect(x0, y0, x0+size, y0+size)

	padding := size / 8
	if padding < 4 {
		padding = 4
	}
	bgRect := logoRect.Inset(-padding).Intersect(bounds)
	draw.Draw(dst, bgRect, &image.Uniform{C: background}, image.Point{}, draw.Src)

	scaled := resizeNearest(logo, size, size)
	draw.Draw(dst, logoRect, scaled, image.Point{}, draw.Over)
}

// 中文：resizeNearest 执行当前包中的对应流程。
// English: resizeNearest executes the corresponding workflow in this package.
func resizeNearest(src image.Image, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	sb := src.Bounds()
	for y := 0; y < height; y++ {
		sy := sb.Min.Y + y*sb.Dy()/height
		for x := 0; x < width; x++ {
			sx := sb.Min.X + x*sb.Dx()/width
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}
