package channel

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// 中文：gatewayHTTPClient 定义当前包使用的数据结构或接口。
// English: gatewayHTTPClient defines a data structure or interface used by this package.
type gatewayHTTPClient interface {
	// 中文：Do 声明该接口需要实现的行为。
	// English: Do declares behavior required by this interface.
	Do(*http.Request) (*http.Response, error)
}

// 中文：gatewayPostJSON 执行当前包中的对应流程。
// English: gatewayPostJSON executes the corresponding workflow in this package.
func gatewayPostJSON(ctx context.Context, client gatewayHTTPClient, endpoint string, payload any, out any) error {
	if client == nil {
		client = http.DefaultClient
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if len(bytes.TrimSpace(data)) == 0 || out == nil {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// 中文：gatewayEndpoint 执行当前包中的对应流程。
// English: gatewayEndpoint executes the corresponding workflow in this package.
func gatewayEndpoint(baseURL, path string) (string, error) {
	if strings.TrimSpace(baseURL) == "" {
		return "", fmt.Errorf("gateway_url is required")
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid gateway_url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("gateway_url must be absolute")
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return u.String(), nil
}

// 中文：gatewaySign 执行当前包中的对应流程。
// English: gatewaySign executes the corresponding workflow in this package.
func gatewaySign(apiKey string, values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key, value := range values {
		if key == "" || value == "" || key == "sign" || key == "signature" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var canonical strings.Builder
	for i, key := range keys {
		if i > 0 {
			canonical.WriteByte('&')
		}
		canonical.WriteString(key)
		canonical.WriteByte('=')
		canonical.WriteString(values[key])
	}

	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(canonical.String()))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

// 中文：gatewayVerifySignature 执行当前包中的对应流程。
// English: gatewayVerifySignature executes the corresponding workflow in this package.
func gatewayVerifySignature(apiKey string, values map[string]string, signature string) bool {
	if signature == "" {
		return false
	}
	expected := gatewaySign(apiKey, values)
	return hmac.Equal([]byte(strings.ToUpper(signature)), []byte(expected))
}

// 中文：gatewaySignedFields 执行当前包中的对应流程。
// English: gatewaySignedFields executes the corresponding workflow in this package.
func gatewaySignedFields(rawData []byte) (map[string]string, string, error) {
	decoder := json.NewDecoder(bytes.NewReader(rawData))
	decoder.UseNumber()

	var payload map[string]any
	if err := decoder.Decode(&payload); err != nil {
		return nil, "", err
	}

	values := make(map[string]string)
	signature := ""
	for key, value := range payload {
		switch key {
		case "sign", "signature":
			signature = gatewayString(value)
			continue
		}
		if scalar, ok := gatewayScalarString(value); ok {
			values[key] = scalar
		}
	}
	return values, signature, nil
}

// 中文：gatewayScalarString 执行当前包中的对应流程。
// English: gatewayScalarString executes the corresponding workflow in this package.
func gatewayScalarString(value any) (string, bool) {
	switch v := value.(type) {
	case nil:
		return "", false
	case string:
		return v, true
	case json.Number:
		return v.String(), true
	case bool:
		return strconv.FormatBool(v), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	default:
		return "", false
	}
}

// 中文：gatewayString 执行当前包中的对应流程。
// English: gatewayString executes the corresponding workflow in this package.
func gatewayString(value any) string {
	if s, ok := gatewayScalarString(value); ok {
		return s
	}
	return fmt.Sprint(value)
}

// 中文：gatewayFirst 执行当前包中的对应流程。
// English: gatewayFirst executes the corresponding workflow in this package.
func gatewayFirst(values map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values[key]); value != "" {
			return value
		}
	}
	return ""
}

// 中文：gatewayParseAmount 执行当前包中的对应流程。
// English: gatewayParseAmount executes the corresponding workflow in this package.
func gatewayParseAmount(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return int64(f*100 + 0.5)
	}
	return 0
}

// 中文：gatewayPaymentStatus 执行当前包中的对应流程。
// English: gatewayPaymentStatus executes the corresponding workflow in this package.
func gatewayPaymentStatus(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "success", "succeeded", "paid", "complete", "completed", "trade_success", "trade_finished", "00":
		return "paid"
	case "pending", "processing", "created", "notpay", "wait_buyer_pay":
		return "pending"
	default:
		return "failed"
	}
}

// 中文：gatewayRefundStatus 执行当前包中的对应流程。
// English: gatewayRefundStatus executes the corresponding workflow in this package.
func gatewayRefundStatus(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "success", "succeeded", "refunded", "complete", "completed", "00":
		return "success"
	default:
		return "refunding"
	}
}

// 中文：gatewayResponseOK 执行当前包中的对应流程。
// English: gatewayResponseOK executes the corresponding workflow in this package.
func gatewayResponseOK(success bool, code string) bool {
	code = strings.TrimSpace(code)
	if success || code == "" {
		return true
	}
	switch strings.ToLower(code) {
	case "0", "00", "success", "ok":
		return true
	default:
		return false
	}
}
