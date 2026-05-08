package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// 中文：GenerateRandomString 执行当前包中的对应流程。
// English: GenerateRandomString executes the corresponding workflow in this package.
// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, (length+1)/2)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// 中文：ContainsString 执行当前包中的对应流程。
// English: ContainsString executes the corresponding workflow in this package.
// ContainsString 检查字符串切片是否包含某字符串
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// 中文：DefaultString 执行当前包中的对应流程。
// English: DefaultString executes the corresponding workflow in this package.
// DefaultString 返回字符串值，如果为空则返回默认值
func DefaultString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// 中文：FormatDuration 执行当前包中的对应流程。
// English: FormatDuration executes the corresponding workflow in this package.
// FormatDuration 格式化时间段
func FormatDuration(d time.Duration) string {
	switch {
	case d >= 24*time.Hour:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d >= time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d >= time.Minute:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d >= time.Second:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	default:
		return d.String()
	}
}

// 中文：MaskString 执行当前包中的对应流程。
// English: MaskString executes the corresponding workflow in this package.
// MaskString 遮蔽字符串（保留前后几位）
func MaskString(s string, head, tail int) string {
	if len(s) <= head+tail {
		return strings.Repeat("*", len(s))
	}
	return s[:head] + strings.Repeat("*", len(s)-head-tail) + s[len(s)-tail:]
}
