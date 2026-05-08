package middleware

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：GinLanguageKey 声明当前包使用的常量。
// English: GinLanguageKey declares constants used by this package.
const GinLanguageKey = "language"

// 中文：I18nConfig 定义当前包使用的数据结构或接口。
// English: I18nConfig defines a data structure or interface used by this package.
type I18nConfig struct {
	// 中文：DefaultLanguage 保存当前结构中的配置或数据值。
	// English: DefaultLanguage stores a configuration or data value for this struct.
	DefaultLanguage string
	// 中文：Supported 保存当前结构中的配置或数据值。
	// English: Supported stores a configuration or data value for this struct.
	Supported []string
	// 中文：Header 保存当前结构中的配置或数据值。
	// English: Header stores a configuration or data value for this struct.
	Header string
}

// 中文：DefaultI18nConfig 执行当前包中的对应流程。
// English: DefaultI18nConfig executes the corresponding workflow in this package.
func DefaultI18nConfig() I18nConfig {
	return I18nConfig{
		DefaultLanguage: "zh-CN",
		Supported:       []string{"zh-CN", "en-US"},
		Header:          "X-Language",
	}
}

// 中文：I18n 执行当前包中的对应流程。
// English: I18n executes the corresponding workflow in this package.
// I18n negotiates the request language and exposes it in Gin and request context.
func I18n(cfg I18nConfig) gin.HandlerFunc {
	if cfg.DefaultLanguage == "" {
		cfg.DefaultLanguage = "zh-CN"
	}
	if cfg.Header == "" {
		cfg.Header = "X-Language"
	}
	supported := normalizeSupported(cfg.Supported, cfg.DefaultLanguage)

	return func(c *gin.Context) {
		language := negotiateLanguage(c.GetHeader(cfg.Header), c.GetHeader("Accept-Language"), supported, cfg.DefaultLanguage)
		c.Set(GinLanguageKey, language)
		c.Header("Content-Language", language)
		ctx := types.WithLanguage(c.Request.Context(), language)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// 中文：Language 执行当前包中的对应流程。
// English: Language executes the corresponding workflow in this package.
func Language(c *gin.Context) string {
	if v, ok := c.Get(GinLanguageKey); ok {
		if language, ok := v.(string); ok {
			return language
		}
	}
	return types.GetLanguage(c.Request.Context())
}

// 中文：normalizeSupported 执行当前包中的对应流程。
// English: normalizeSupported executes the corresponding workflow in this package.
func normalizeSupported(values []string, fallback string) map[string]string {
	if len(values) == 0 {
		values = []string{fallback}
	}
	out := make(map[string]string, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		out[strings.ToLower(value)] = value
		if base, _, ok := strings.Cut(value, "-"); ok && base != "" {
			if _, exists := out[strings.ToLower(base)]; !exists {
				out[strings.ToLower(base)] = value
			}
		}
	}
	if len(out) == 0 {
		out[strings.ToLower(fallback)] = fallback
	}
	return out
}

// 中文：negotiateLanguage 执行当前包中的对应流程。
// English: negotiateLanguage executes the corresponding workflow in this package.
func negotiateLanguage(headerLanguage, acceptLanguage string, supported map[string]string, fallback string) string {
	for _, candidate := range append(languageCandidates(headerLanguage), languageCandidates(acceptLanguage)...) {
		if language, ok := supported[strings.ToLower(candidate.tag)]; ok {
			return language
		}
		if base, _, ok := strings.Cut(candidate.tag, "-"); ok {
			if language, ok := supported[strings.ToLower(base)]; ok {
				return language
			}
		}
	}
	if language, ok := supported[strings.ToLower(fallback)]; ok {
		return language
	}
	return fallback
}

// 中文：languageCandidate 定义当前包使用的数据结构或接口。
// English: languageCandidate defines a data structure or interface used by this package.
type languageCandidate struct {
	// 中文：tag 保存当前结构中的配置或数据值。
	// English: tag stores a configuration or data value for this struct.
	tag string
	// 中文：q 保存当前结构中的配置或数据值。
	// English: q stores a configuration or data value for this struct.
	q float64
}

// 中文：languageCandidates 执行当前包中的对应流程。
// English: languageCandidates executes the corresponding workflow in this package.
func languageCandidates(value string) []languageCandidate {
	parts := strings.Split(value, ",")
	candidates := make([]languageCandidate, 0, len(parts))
	for _, part := range parts {
		tag, q := parseLanguagePart(part)
		if tag == "" || tag == "*" || q <= 0 {
			continue
		}
		candidates = append(candidates, languageCandidate{tag: tag, q: q})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].q > candidates[j].q
	})
	return candidates
}

// 中文：parseLanguagePart 执行当前包中的对应流程。
// English: parseLanguagePart executes the corresponding workflow in this package.
func parseLanguagePart(part string) (string, float64) {
	part = strings.TrimSpace(part)
	if part == "" {
		return "", 0
	}
	pieces := strings.Split(part, ";")
	tag := strings.TrimSpace(pieces[0])
	q := 1.0
	for _, piece := range pieces[1:] {
		key, value, ok := strings.Cut(strings.TrimSpace(piece), "=")
		if !ok || strings.TrimSpace(key) != "q" {
			continue
		}
		parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err == nil {
			q = parsed
		}
	}
	return tag, q
}
