package channel

import (
	"strings"
	"time"
)

// 中文：unionPayCloseOrderID 执行当前包中的对应流程。
// English: unionPayCloseOrderID executes the corresponding workflow in this package.
func unionPayCloseOrderID(origQryID string) string {
	suffix := sanitizeUnionPayOrderID(origQryID)
	if len(suffix) > 16 {
		suffix = suffix[len(suffix)-16:]
	}
	if suffix == "" {
		suffix = "VOID"
	}
	id := "V" + time.Now().Format("20060102150405") + suffix
	if len(id) > 40 {
		return id[:40]
	}
	return id
}

// 中文：sanitizeUnionPayOrderID 执行当前包中的对应流程。
// English: sanitizeUnionPayOrderID executes the corresponding workflow in this package.
func sanitizeUnionPayOrderID(value string) string {
	var b strings.Builder
	for _, r := range value {
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			b.WriteRune(r)
		}
	}
	return b.String()
}
