package service

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spiringo/spiringo/internal/modules/qrcode/dto"
	"github.com/spiringo/spiringo/internal/modules/qrcode/model"
	"github.com/spiringo/spiringo/internal/modules/qrcode/repository"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：TestBuildQRCodePNGAppliesColors 验证相关行为符合预期。
// English: TestBuildQRCodePNGAppliesColors verifies the related behavior.
func TestBuildQRCodePNGAppliesColors(t *testing.T) {
	data, err := buildQRCodePNG(context.Background(), dto.GenerateReq{
		Content:         "https://example.com",
		ForegroundColor: "#123456",
		BackgroundColor: "#fefefe",
	}, 128, "medium")
	if err != nil {
		t.Fatalf("build qrcode: %v", err)
	}
	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("decode generated png: %v", err)
	}
}

// 中文：TestBuildQRCodePNGRejectsInvalidColor 验证相关行为符合预期。
// English: TestBuildQRCodePNGRejectsInvalidColor verifies the related behavior.
func TestBuildQRCodePNGRejectsInvalidColor(t *testing.T) {
	_, err := buildQRCodePNG(context.Background(), dto.GenerateReq{
		Content:         "https://example.com",
		ForegroundColor: "not-a-color",
	}, 128, "medium")
	if err == nil {
		t.Fatal("expected invalid color error")
	}
}

// 中文：TestBuildQRCodePNGOverlaysLogo 验证相关行为符合预期。
// English: TestBuildQRCodePNGOverlaysLogo verifies the related behavior.
func TestBuildQRCodePNGOverlaysLogo(t *testing.T) {
	var logo bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	if err := png.Encode(&logo, img); err != nil {
		t.Fatalf("encode logo: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(logo.Bytes())
	}))
	defer srv.Close()

	data, err := buildQRCodePNG(context.Background(), dto.GenerateReq{
		Content:  "https://example.com",
		LogoURL:  srv.URL,
		LogoSize: 24,
	}, 128, "high")
	if err != nil {
		t.Fatalf("build qrcode with logo: %v", err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode generated png: %v", err)
	}
	r, g, b, _ := decoded.At(decoded.Bounds().Dx()/2, decoded.Bounds().Dy()/2).RGBA()
	if r <= g || r <= b {
		t.Fatalf("center pixel did not contain the red logo: r=%d g=%d b=%d", r, g, b)
	}
}

// 中文：TestParseLevelAcceptsBlueprintRecoveryCodes 验证相关行为符合预期。
// English: TestParseLevelAcceptsBlueprintRecoveryCodes verifies the related behavior.
func TestParseLevelAcceptsBlueprintRecoveryCodes(t *testing.T) {
	if parseLevel("L") != parseLevel("low") {
		t.Fatal("L should map to low recovery level")
	}
	if parseLevel("M") != parseLevel("medium") {
		t.Fatal("M should map to medium recovery level")
	}
	if parseLevel("Q") != parseLevel("high") {
		t.Fatal("Q should map to quartile/high recovery level")
	}
	if parseLevel("H") != parseLevel("highest") {
		t.Fatal("H should map to highest recovery level")
	}
}

// 中文：TestGenerateDoesNotPersistRecordWhenImageBuildFails 验证相关行为符合预期。
// English: TestGenerateDoesNotPersistRecordWhenImageBuildFails verifies the related behavior.
func TestGenerateDoesNotPersistRecordWhenImageBuildFails(t *testing.T) {
	db, err := orm.New(orm.Config{Driver: "sqlite", DSN: filepath.Join(t.TempDir(), "qrcode.db")})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.QRCodeRecord{}, &model.ScanLog{}); err != nil {
		t.Fatalf("migrate qrcode: %v", err)
	}

	repo := repository.NewQRCodeRepository(orm.NewTenantDB(db), db)
	svc := NewQRCodeService(Config{DefaultSize: 128, DefaultLevel: "M"}, repo, nil)
	if _, err := svc.Generate(context.Background(), dto.GenerateReq{
		Content:         "https://example.com",
		ForegroundColor: "not-a-color",
	}); err == nil {
		t.Fatal("expected invalid color error")
	}

	var total int64
	if err := db.DB().Model(&model.QRCodeRecord{}).Count(&total).Error; err != nil {
		t.Fatalf("count records: %v", err)
	}
	if total != 0 {
		t.Fatalf("record count = %d, want 0", total)
	}
}

// 中文：TestQRCodeObjectPrefixHelpers 验证相关行为符合预期。
// English: TestQRCodeObjectPrefixHelpers verifies the related behavior.
func TestQRCodeObjectPrefixHelpers(t *testing.T) {
	if got := qrcodeObjectKey("custom/prefix/", "abc123"); got != "custom/prefix/abc123.png" {
		t.Fatalf("object key = %q", got)
	}
	if got := qrcodeObjectKey("", "abc123"); got != "qrcode/abc123.png" {
		t.Fatalf("default object key = %q", got)
	}
	if got := qrcodeFallbackURL("https://cdn.example.com/qrcode/", "abc123"); got != "https://cdn.example.com/qrcode/abc123.png" {
		t.Fatalf("fallback url = %q", got)
	}
}

// 中文：TestQRCodeGeneratorFallbackImplementsBlueprintContract 验证相关行为符合预期。
// English: TestQRCodeGeneratorFallbackImplementsBlueprintContract verifies the related behavior.
func TestQRCodeGeneratorFallbackImplementsBlueprintContract(t *testing.T) {
	generator := NewQRCodeGenerator(nil)
	result, err := generator.Generate(context.Background(), "https://example.com", &QRCodeStyle{
		Size:            128,
		ForegroundColor: "#123456",
		BackgroundColor: "#ffffff",
		Level:           "medium",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if string(result.Content) != "https://example.com" {
		t.Fatalf("content = %q", result.Content)
	}
	if len(result.ImageBytes) == 0 {
		t.Fatal("empty image bytes")
	}
	parsed, err := generator.Parse(context.Background(), result.ImageBytes)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if parsed != "https://example.com" {
		t.Fatalf("parsed = %q", parsed)
	}
}
