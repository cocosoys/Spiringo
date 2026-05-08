package channel

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

// 中文：AlipayChannel 定义当前包使用的数据结构或接口。
// English: AlipayChannel defines a data structure or interface used by this package.
// AlipayChannel 支付宝支付通道
type AlipayChannel struct {
	// 中文：AppID 保存当前结构中的配置或数据值。
	// English: AppID stores a configuration or data value for this struct.
	AppID string
	// 中文：PrivateKey 保存当前结构中的配置或数据值。
	// English: PrivateKey stores a configuration or data value for this struct.
	PrivateKey string // 应用私钥
	// 中文：AlipayPublicCert 保存当前结构中的配置或数据值。
	// English: AlipayPublicCert stores a configuration or data value for this struct.
	AlipayPublicCert []byte // 支付宝公钥证书内容
	// 中文：AppPublicCert 保存当前结构中的配置或数据值。
	// English: AppPublicCert stores a configuration or data value for this struct.
	AppPublicCert []byte // 应用公钥证书内容
	// 中文：AlipayRootCert 保存当前结构中的配置或数据值。
	// English: AlipayRootCert stores a configuration or data value for this struct.
	AlipayRootCert []byte // 支付宝根证书内容
	// 中文：IsProd 保存当前结构中的配置或数据值。
	// English: IsProd stores a configuration or data value for this struct.
	IsProd bool
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *alipay.Client
}

// 中文：NewAlipayChannel 创建并返回对应组件实例。
// English: NewAlipayChannel creates and returns the corresponding component instance.
func NewAlipayChannel(appID, privateKey, notifyURL string) *AlipayChannel {
	return &AlipayChannel{AppID: appID, PrivateKey: privateKey, NotifyURL: notifyURL}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c *AlipayChannel) Name() string { return "alipay" }

// 中文：initClient 执行当前包中的对应流程。
// English: initClient executes the corresponding workflow in this package.
// initClient 懒初始化支付宝客户端
func (c *AlipayChannel) initClient() error {
	if c.client != nil {
		return nil
	}

	client, err := alipay.NewClient(c.AppID, c.PrivateKey, c.IsProd)
	if err != nil {
		return fmt.Errorf("alipay: init client: %w", err)
	}

	// 设置证书SN（证书模式）
	if len(c.AppPublicCert) > 0 && len(c.AlipayPublicCert) > 0 && len(c.AlipayRootCert) > 0 {
		if err := client.SetCertSnByContent(c.AppPublicCert, c.AlipayPublicCert, c.AlipayRootCert); err != nil {
			return fmt.Errorf("alipay: set cert sn: %w", err)
		}
		// 开启自动验签
		client.AutoVerifySign(c.AlipayPublicCert)
	}

	c.client = client
	return nil
}

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c *AlipayChannel) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	notify := notifyURL
	if notify == "" {
		notify = c.NotifyURL
	}
	c.client.SetNotifyUrl(notify)
	if returnURL != "" {
		c.client.SetReturnUrl(returnURL)
	}

	// 金额：分转元
	amountYuan := fmt.Sprintf("%.2f", float64(amount)/100.0)

	bm := gopay.BodyMap{
		"out_trade_no": outTradeNo,
		"total_amount": amountYuan,
		"subject":      subject,
	}

	switch scene {
	case "qrcode", "native":
		// 当面付（扫码）
		rsp, err := c.client.TradePrecreate(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("alipay precreate: %w", err)
		}
		if rsp.Response.Code != "10000" {
			return nil, fmt.Errorf("alipay precreate: code=%s, msg=%s", rsp.Response.Code, rsp.Response.Msg)
		}
		return &PayResult{QrCode: rsp.Response.QrCode}, nil

	case "h5", "wap":
		// 手机网站支付
		bm.Set("product_code", "QUICK_WAP_WAY")
		payURL, err := c.client.TradeWapPay(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("alipay wap pay: %w", err)
		}
		return &PayResult{PayURL: payURL}, nil

	case "page", "pc":
		// PC网站支付
		bm.Set("product_code", "FAST_INSTANT_TRADE_PAY")
		payURL, err := c.client.TradePagePay(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("alipay page pay: %w", err)
		}
		return &PayResult{PayURL: payURL}, nil

	case "app":
		// APP支付
		orderStr, err := c.client.TradeAppPay(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("alipay app pay: %w", err)
		}
		return &PayResult{Params: orderStr}, nil

	default:
		return nil, fmt.Errorf("unsupported alipay scene: %s", scene)
	}
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c *AlipayChannel) VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	// 解析回调form数据
	values, err := url.ParseQuery(string(rawData))
	if err != nil {
		return nil, fmt.Errorf("alipay: parse callback: %w", err)
	}

	// 验签
	bm := gopay.BodyMap{}
	for k, vs := range values {
		if len(vs) > 0 {
			bm.Set(k, vs[0])
		}
	}

	if len(c.AlipayPublicCert) > 0 {
		ok, err := alipay.VerifySignWithCert(c.AlipayPublicCert, bm)
		if err != nil {
			return nil, fmt.Errorf("alipay: verify sign: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("alipay: sign verification failed")
		}
	}

	outTradeNo := values.Get("out_trade_no")
	tradeNo := values.Get("trade_no")
	tradeStatus := values.Get("trade_status")
	totalAmountStr := values.Get("total_amount")

	status := "failed"
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		status = "paid"
	}

	// 金额(元)转为分
	var amountCents int64
	if f, err := strconv.ParseFloat(totalAmountStr, 64); err == nil {
		amountCents = int64(f * 100)
	}

	return &CallbackResult{
		OutTradeNo: outTradeNo,
		TradeNo:    tradeNo,
		Status:     status,
		Amount:     amountCents,
		RawData:    rawData,
	}, nil
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c *AlipayChannel) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	amountYuan := fmt.Sprintf("%.2f", float64(refundAmount)/100.0)

	bm := gopay.BodyMap{
		"out_trade_no":   outTradeNo,
		"out_request_no": outRefundNo,
		"refund_amount":  amountYuan,
		"refund_reason":  reason,
	}

	rsp, err := c.client.TradeRefund(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("alipay refund: %w", err)
	}
	if rsp.Response.Code != "10000" {
		return nil, fmt.Errorf("alipay refund: code=%s, msg=%s", rsp.Response.Code, rsp.Response.Msg)
	}

	status := "refunding"
	if rsp.Response.RefundFee != "" {
		status = "success"
	}

	return &RefundResult{
		RefundNo: rsp.Response.TradeNo,
		Status:   status,
	}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c *AlipayChannel) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	bm := gopay.BodyMap{"out_trade_no": outTradeNo}
	rsp, err := c.client.TradeQuery(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("alipay query: %w", err)
	}
	if rsp.Response.Code != "10000" {
		return nil, fmt.Errorf("alipay query: code=%s, msg=%s", rsp.Response.Code, rsp.Response.Msg)
	}

	status := "pending"
	switch rsp.Response.TradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		status = "paid"
	case "TRADE_CLOSED":
		status = "failed"
	}

	var amountCents int64
	if f, err := strconv.ParseFloat(rsp.Response.TotalAmount, 64); err == nil {
		amountCents = int64(f * 100)
	}

	return &CallbackResult{
		OutTradeNo: rsp.Response.OutTradeNo,
		TradeNo:    rsp.Response.TradeNo,
		Status:     status,
		Amount:     amountCents,
	}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
// SetCerts 设置证书（证书模式必需）
func (c *AlipayChannel) ClosePayment(ctx context.Context, outTradeNo string) error {
	if err := c.initClient(); err != nil {
		return err
	}
	if outTradeNo == "" {
		return fmt.Errorf("alipay: out_trade_no is required")
	}

	rsp, err := c.client.TradeClose(ctx, gopay.BodyMap{"out_trade_no": outTradeNo})
	if err != nil {
		return fmt.Errorf("alipay close payment: %w", err)
	}
	if rsp.Response == nil {
		return fmt.Errorf("alipay close payment: empty response")
	}
	if rsp.Response.Code != "10000" {
		return fmt.Errorf("alipay close payment: code=%s, msg=%s", rsp.Response.Code, rsp.Response.Msg)
	}
	return nil
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c *AlipayChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c *AlipayChannel) CallbackFail() any { return "failure" }

// 中文：SetCerts 执行当前包中的对应流程。
// English: SetCerts executes the corresponding workflow in this package.
func (c *AlipayChannel) SetCerts(appPublicCert, alipayPublicCert, alipayRootCert []byte) {
	c.AppPublicCert = appPublicCert
	c.AlipayPublicCert = alipayPublicCert
	c.AlipayRootCert = alipayRootCert
}

// 中文：SetProd 执行当前包中的对应流程。
// English: SetProd executes the corresponding workflow in this package.
// SetProd 设置是否生产环境
func (c *AlipayChannel) SetProd(isProd bool) {
	c.IsProd = isProd
}
