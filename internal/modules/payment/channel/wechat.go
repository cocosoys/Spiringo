package channel

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

// 中文：WechatChannel 定义当前包使用的数据结构或接口。
// English: WechatChannel defines a data structure or interface used by this package.
// WechatChannel 微信支付通道（V3接口）
type WechatChannel struct {
	// 中文：AppID 保存当前结构中的配置或数据值。
	// English: AppID stores a configuration or data value for this struct.
	AppID string
	// 中文：MchID 保存当前结构中的配置或数据值。
	// English: MchID stores a configuration or data value for this struct.
	MchID string
	// 中文：APIV3Key 保存当前结构中的配置或数据值。
	// English: APIV3Key stores a configuration or data value for this struct.
	APIV3Key string
	// 中文：SerialNo 保存当前结构中的配置或数据值。
	// English: SerialNo stores a configuration or data value for this struct.
	SerialNo string // 商户API证书序列号（V3必需）
	// 中文：PrivateKey 保存当前结构中的配置或数据值。
	// English: PrivateKey stores a configuration or data value for this struct.
	PrivateKey string // 商户私钥 PEM 内容（V3必需）
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *wechat.ClientV3
}

// 中文：NewWechatChannel 创建并返回对应组件实例。
// English: NewWechatChannel creates and returns the corresponding component instance.
func NewWechatChannel(appID, mchID, apiKey, notifyURL string) *WechatChannel {
	return &WechatChannel{AppID: appID, MchID: mchID, APIV3Key: apiKey, NotifyURL: notifyURL}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c *WechatChannel) Name() string { return "wechat" }

// 中文：initClient 执行当前包中的对应流程。
// English: initClient executes the corresponding workflow in this package.
// initClient 懒初始化微信支付V3客户端
func (c *WechatChannel) initClient() error {
	if c.client != nil {
		return nil
	}
	if c.PrivateKey == "" || c.SerialNo == "" {
		return fmt.Errorf("wechat: SerialNo and PrivateKey are required for V3 API")
	}
	client, err := wechat.NewClientV3(c.MchID, c.SerialNo, c.APIV3Key, c.PrivateKey)
	if err != nil {
		return fmt.Errorf("wechat: init client: %w", err)
	}
	// 自动获取微信平台证书并验签
	if err := client.AutoVerifySign(); err != nil {
		return fmt.Errorf("wechat: auto verify sign: %w", err)
	}
	c.client = client
	return nil
}

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c *WechatChannel) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	notify := notifyURL
	if notify == "" {
		notify = c.NotifyURL
	}

	bm := gopay.BodyMap{
		"appid":        c.AppID,
		"description":  subject,
		"out_trade_no": outTradeNo,
		"notify_url":   notify,
		"amount": gopay.BodyMap{
			"total":    amount,
			"currency": "CNY",
		},
	}

	switch scene {
	case "native", "qrcode":
		rsp, err := c.client.V3TransactionNative(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("wechat native pay: %w", err)
		}
		if rsp.Code != wechat.Success {
			return nil, fmt.Errorf("wechat native pay: code=%d, msg=%s", rsp.Code, rsp.Error)
		}
		return &PayResult{QrCode: rsp.Response.CodeUrl}, nil

	case "jsapi":
		bm.Set("payer", gopay.BodyMap{"openid": openID})
		rsp, err := c.client.V3TransactionJsapi(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("wechat jsapi pay: %w", err)
		}
		if rsp.Code != wechat.Success {
			return nil, fmt.Errorf("wechat jsapi pay: code=%d, msg=%s", rsp.Code, rsp.Error)
		}
		jsapiParams, err := c.client.PaySignOfJSAPI(c.AppID, rsp.Response.PrepayId)
		if err != nil {
			return nil, fmt.Errorf("wechat jsapi sign: %w", err)
		}
		return &PayResult{PrepayID: rsp.Response.PrepayId, Params: jsapiParams}, nil

	case "h5":
		bm.Set("scene_info", gopay.BodyMap{
			"payer_client_ip": "127.0.0.1",
			"h5_info": gopay.BodyMap{
				"type": "Wap",
			},
		})
		rsp, err := c.client.V3TransactionH5(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("wechat h5 pay: %w", err)
		}
		if rsp.Code != wechat.Success {
			return nil, fmt.Errorf("wechat h5 pay: code=%d, msg=%s", rsp.Code, rsp.Error)
		}
		return &PayResult{PayURL: rsp.Response.H5Url}, nil

	case "app":
		rsp, err := c.client.V3TransactionApp(ctx, bm)
		if err != nil {
			return nil, fmt.Errorf("wechat app pay: %w", err)
		}
		if rsp.Code != wechat.Success {
			return nil, fmt.Errorf("wechat app pay: code=%d, msg=%s", rsp.Code, rsp.Error)
		}
		appParams, err := c.client.PaySignOfApp(c.AppID, rsp.Response.PrepayId)
		if err != nil {
			return nil, fmt.Errorf("wechat app sign: %w", err)
		}
		return &PayResult{PrepayID: rsp.Response.PrepayId, Params: appParams}, nil

	default:
		return nil, fmt.Errorf("unsupported wechat scene: %s", scene)
	}
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c *WechatChannel) VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(rawData))
	req.Header.Set("Content-Type", "application/json")
	return c.VerifyCallbackWithRequest(ctx, req, rawData)
}

// 中文：VerifyCallbackWithRequest 执行当前包中的对应流程。
// English: VerifyCallbackWithRequest executes the corresponding workflow in this package.
func (c *WechatChannel) VerifyCallbackWithRequest(ctx context.Context, req *http.Request, rawData []byte) (*CallbackResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("wechat: original http request is required")
	}
	req = req.WithContext(ctx)
	req.Body = http.NoBody
	if rawData != nil {
		req.Body = io.NopCloser(bytes.NewReader(rawData))
	}
	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		return nil, fmt.Errorf("wechat: parse notify: %w", err)
	}

	// 验签
	if err := notifyReq.VerifySignByPKMap(c.client.WxPublicKeyMap()); err != nil {
		return nil, fmt.Errorf("wechat: verify sign: %w", err)
	}

	// 解密
	payResult, err := notifyReq.DecryptPayCipherText(c.APIV3Key)
	if err != nil {
		return nil, fmt.Errorf("wechat: decrypt notify: %w", err)
	}

	status := "failed"
	if payResult.TradeState == "SUCCESS" {
		status = "paid"
	}

	return &CallbackResult{
		OutTradeNo: payResult.OutTradeNo,
		TradeNo:    payResult.TransactionId,
		Status:     status,
		Amount:     int64(payResult.Amount.Total),
		RawData:    rawData,
	}, nil
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c *WechatChannel) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	bm := gopay.BodyMap{
		"out_trade_no":  outTradeNo,
		"out_refund_no": outRefundNo,
		"amount": gopay.BodyMap{
			"refund":   refundAmount,
			"total":    totalAmount,
			"currency": "CNY",
		},
		"reason": reason,
	}

	rsp, err := c.client.V3Refund(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("wechat refund: %w", err)
	}
	if rsp.Code != wechat.Success {
		return nil, fmt.Errorf("wechat refund: code=%d, msg=%s", rsp.Code, rsp.Error)
	}

	status := "refunding"
	if rsp.Response.Status == "SUCCESS" {
		status = "success"
	}

	return &RefundResult{
		RefundNo: rsp.Response.RefundId,
		Status:   status,
	}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c *WechatChannel) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	rsp, err := c.client.V3TransactionQueryOrder(ctx, wechat.OutTradeNo, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("wechat query: %w", err)
	}
	if rsp.Code != wechat.Success {
		return nil, fmt.Errorf("wechat query: code=%d, msg=%s", rsp.Code, rsp.Error)
	}

	status := "failed"
	if rsp.Response.TradeState == "SUCCESS" {
		status = "paid"
	} else if rsp.Response.TradeState == "NOTPAY" {
		status = "pending"
	}

	return &CallbackResult{
		OutTradeNo: rsp.Response.OutTradeNo,
		TradeNo:    rsp.Response.TransactionId,
		Status:     status,
		Amount:     int64(rsp.Response.Amount.Total),
	}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
// SetSerialNo 设置商户API证书序列号（V3必需）
func (c *WechatChannel) ClosePayment(ctx context.Context, outTradeNo string) error {
	if err := c.initClient(); err != nil {
		return err
	}
	if outTradeNo == "" {
		return fmt.Errorf("wechat: out_trade_no is required")
	}

	rsp, err := c.client.V3TransactionCloseOrder(ctx, outTradeNo)
	if err != nil {
		return fmt.Errorf("wechat close payment: %w", err)
	}
	if rsp.Code != wechat.Success {
		return fmt.Errorf("wechat close payment: code=%d, msg=%s", rsp.Code, rsp.Error)
	}
	return nil
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c *WechatChannel) CallbackSuccess() any {
	return map[string]string{"code": "SUCCESS", "message": "success"}
}

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c *WechatChannel) CallbackFail() any {
	return map[string]string{"code": "FAIL", "message": "fail"}
}

// 中文：SetSerialNo 执行当前包中的对应流程。
// English: SetSerialNo executes the corresponding workflow in this package.
func (c *WechatChannel) SetSerialNo(serialNo string) {
	c.SerialNo = serialNo
}

// 中文：SetPrivateKey 执行当前包中的对应流程。
// English: SetPrivateKey executes the corresponding workflow in this package.
// SetPrivateKey 设置商户私钥PEM内容（V3必需）
func (c *WechatChannel) SetPrivateKey(pem string) {
	c.PrivateKey = pem
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
var _ = strconv.Itoa
