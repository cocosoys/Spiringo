package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// 中文：SHA256 执行当前包中的对应流程。
// English: SHA256 executes the corresponding workflow in this package.
// SHA256 计算SHA256哈希
func SHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// 中文：SHA256String 执行当前包中的对应流程。
// English: SHA256String executes the corresponding workflow in this package.
// SHA256String 计算字符串的SHA256
func SHA256String(s string) string {
	return SHA256([]byte(s))
}
