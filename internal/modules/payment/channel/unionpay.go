package channel

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// 中文：UnionPayChannel 定义当前包使用的数据结构或接口。
// English: UnionPayChannel defines a data structure or interface used by this package.
// UnionPayChannel 银联支付通道
// 银联无成熟Go SDK，使用HTTP直连全渠道API
type UnionPayChannel struct {
	// 中文：MchID 保存当前结构中的配置或数据值。
	// English: MchID stores a configuration or data value for this struct.
	MchID string
	// 中文：APIKey 保存当前结构中的配置或数据值。
	// English: APIKey stores a configuration or data value for this struct.
	APIKey string
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string
	// 中文：Sandbox 保存当前结构中的配置或数据值。
	// English: Sandbox stores a configuration or data value for this struct.
	Sandbox bool
}

// 中文：NewUnionPayChannel 创建并返回对应组件实例。
// English: NewUnionPayChannel creates and returns the corresponding component instance.
func NewUnionPayChannel(mchID, apiKey, notifyURL string) *UnionPayChannel {
	return &UnionPayChannel{MchID: mchID, APIKey: apiKey, NotifyURL: notifyURL}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c *UnionPayChannel) Name() string { return "unionpay" }

// 中文：gatewayURL 执行当前包中的对应流程。
// English: gatewayURL executes the corresponding workflow in this package.
func (c *UnionPayChannel) gatewayURL() string {
	if c.Sandbox {
		return "https://gateway.test.95516.com/gateway/api/frontTransReq.do"
	}
	return "https://gateway.95516.com/gateway/api/frontTransReq.do"
}

// 中文：queryURL 执行当前包中的对应流程。
// English: queryURL executes the corresponding workflow in this package.
func (c *UnionPayChannel) queryURL() string {
	if c.Sandbox {
		return "https://gateway.test.95516.com/gateway/api/singleQuery.do"
	}
	return "https://gateway.95516.com/gateway/api/singleQuery.do"
}

// 中文：refundURL 执行当前包中的对应流程。
// English: refundURL executes the corresponding workflow in this package.
func (c *UnionPayChannel) refundURL() string {
	if c.Sandbox {
		return "https://gateway.test.95516.com/gateway/api/refund.do"
	}
	return "https://gateway.95516.com/gateway/api/refund.do"
}

// 中文：voidURL 执行当前包中的对应流程。
// English: voidURL executes the corresponding workflow in this package.
func (c *UnionPayChannel) voidURL() string {
	if c.Sandbox {
		return "https://gateway.test.95516.com/gateway/api/backTransReq.do"
	}
	return "https://gateway.95516.com/gateway/api/backTransReq.do"
}

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c *UnionPayChannel) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	notify := notifyURL
	if notify == "" {
		notify = c.NotifyURL
	}

	params := map[string]string{
		"version":      "5.1.0",
		"encoding":     "UTF-8",
		"signMethod":   "01", // SHA256
		"txnType":      "01", // 消费
		"txnSubType":   "01",
		"bizType":      "000201",
		"channelType":  "07", // PC
		"merId":        c.MchID,
		"orderId":      outTradeNo,
		"txnTime":      time.Now().Format("20060102150405"),
		"txnAmt":       fmt.Sprintf("%d", amount),
		"currencyCode": "156", // CNY
		"frontUrl":     returnURL,
		"backUrl":      notify,
		"accessType":   "0",
		"subject":      subject,
	}

	if scene == "qrcode" || scene == "native" {
		params["qrPayMode"] = "2"
	}

	params["signature"] = c.sign(params)

	// 构造网关跳转URL
	v := make(url.Values, len(params))
	for k, val := range params {
		v.Set(k, val)
	}
	payURL := c.gatewayURL() + "?" + v.Encode()

	return &PayResult{PayURL: payURL}, nil
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c *UnionPayChannel) VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error) {
	values, err := url.ParseQuery(string(rawData))
	if err != nil {
		return nil, fmt.Errorf("unionpay: parse callback: %w", err)
	}

	// 提取签名参数
	signature := values.Get("signature")

	// 构造待签名字符串
	params := make(map[string]string)
	for k := range values {
		if k == "signature" || values.Get(k) == "" {
			continue
		}
		params[k] = values.Get(k)
	}

	// 验签
	expectedSign := c.sign(params)
	if !strings.EqualFold(signature, expectedSign) {
		return nil, fmt.Errorf("unionpay: signature mismatch")
	}

	status := "failed"
	if values.Get("respCode") == "00" {
		status = "paid"
	}

	return &CallbackResult{
		OutTradeNo: values.Get("orderId"),
		TradeNo:    values.Get("queryId"),
		Status:     status,
		RawData:    rawData,
	}, nil
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c *UnionPayChannel) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	params := map[string]string{
		"version":     "5.1.0",
		"encoding":    "UTF-8",
		"signMethod":  "01",
		"txnType":     "04", // 退货
		"txnSubType":  "01",
		"bizType":     "000201",
		"channelType": "07",
		"merId":       c.MchID,
		"orderId":     outRefundNo,
		"origQryId":   outTradeNo,
		"txnTime":     time.Now().Format("20060102150405"),
		"txnAmt":      fmt.Sprintf("%d", refundAmount),
		"backUrl":     c.NotifyURL,
	}
	if reason != "" {
		params["reqReserved"] = reason
	}
	params["signature"] = c.sign(params)

	// 发送退款请求
	resp, err := unionPayHTTPPost(ctx, c.refundURL(), params)
	if err != nil {
		return nil, fmt.Errorf("unionpay refund: %w", err)
	}

	status := "refunding"
	if resp["respCode"] == "00" {
		status = "success"
	}

	return &RefundResult{
		RefundNo: resp["queryId"],
		Status:   status,
	}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c *UnionPayChannel) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	params := map[string]string{
		"version":    "5.1.0",
		"encoding":   "UTF-8",
		"signMethod": "01",
		"txnType":    "00", // 查询
		"txnSubType": "00",
		"bizType":    "000201",
		"merId":      c.MchID,
		"orderId":    outTradeNo,
		"txnTime":    time.Now().Format("20060102150405"),
	}
	params["signature"] = c.sign(params)

	resp, err := unionPayHTTPPost(ctx, c.queryURL(), params)
	if err != nil {
		return nil, fmt.Errorf("unionpay query: %w", err)
	}

	status := "pending"
	switch resp["origRespCode"] {
	case "00":
		status = "paid"
	case "03", "04", "05":
		status = "pending"
	case "01", "02", "06":
		status = "failed"
	}

	return &CallbackResult{
		OutTradeNo: resp["orderId"],
		TradeNo:    resp["queryId"],
		Status:     status,
	}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
// sign 银联SHA256签名
func (c *UnionPayChannel) ClosePayment(ctx context.Context, origQryID string) error {
	if origQryID == "" {
		return fmt.Errorf("unionpay: origQryId is required")
	}

	params := map[string]string{
		"version":     "5.1.0",
		"encoding":    "UTF-8",
		"signMethod":  "01",
		"txnType":     "31",
		"txnSubType":  "00",
		"bizType":     "000201",
		"channelType": "07",
		"merId":       c.MchID,
		"orderId":     unionPayCloseOrderID(origQryID),
		"origQryId":   origQryID,
		"txnTime":     time.Now().Format("20060102150405"),
		"backUrl":     c.NotifyURL,
	}
	params["signature"] = c.sign(params)

	resp, err := unionPayHTTPPost(ctx, c.voidURL(), params)
	if err != nil {
		return fmt.Errorf("unionpay close payment: %w", err)
	}
	if resp["respCode"] != "00" {
		return fmt.Errorf("unionpay close payment failed: %s %s", resp["respCode"], firstNonEmptyString(resp["respMsg"], resp["origRespMsg"]))
	}
	return nil
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c *UnionPayChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c *UnionPayChannel) CallbackFail() any { return "fail" }

// 中文：sign 执行当前包中的对应流程。
// English: sign executes the corresponding workflow in this package.
func (c *UnionPayChannel) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "signature" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}

	h := sha256.New()
	h.Write([]byte(buf.String() + "&key=" + c.APIKey))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// 中文：unionPayHTTPPost 声明当前包使用的变量。
// English: unionPayHTTPPost declares variables used by this package.
// unionPayHTTPPost 发送银联HTTP POST请求
var unionPayHTTPPost = defaultUnionPayHTTPPost

// 中文：defaultUnionPayHTTPPost 执行当前包中的对应流程。
// English: defaultUnionPayHTTPPost executes the corresponding workflow in this package.
func defaultUnionPayHTTPPost(ctx context.Context, endpoint string, params map[string]string) (map[string]string, error) {
	v := make(url.Values, len(params))
	for k, val := range params {
		v.Set(k, val)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	pairs := strings.Split(string(body), "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[kv[0]], _ = url.QueryUnescape(kv[1])
		}
	}

	return result, nil
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
// 确保 bytes 被使用
var _ = bytes.NewReader
