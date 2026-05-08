package channel

import (
	"context"
	"testing"
)

// 中文：registryTestChannel 定义当前包使用的数据结构或接口。
// English: registryTestChannel defines a data structure or interface used by this package.
type registryTestChannel struct {
	// 中文：name 保存当前结构中的配置或数据值。
	// English: name stores a configuration or data value for this struct.
	name string
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c registryTestChannel) Name() string { return c.name }

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c registryTestChannel) CreatePayment(context.Context, string, string, int64, string, string, string, string) (*PayResult, error) {
	return &PayResult{}, nil
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c registryTestChannel) VerifyCallback(context.Context, []byte) (*CallbackResult, error) {
	return &CallbackResult{}, nil
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c registryTestChannel) Refund(context.Context, string, string, int64, int64, string) (*RefundResult, error) {
	return &RefundResult{}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c registryTestChannel) QueryPayment(context.Context, string) (*CallbackResult, error) {
	return &CallbackResult{}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
func (c registryTestChannel) ClosePayment(context.Context, string) error { return nil }

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c registryTestChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c registryTestChannel) CallbackFail() any { return "fail" }

// 中文：TestRegistryListReturnsChannelsByName 验证相关行为符合预期。
// English: TestRegistryListReturnsChannelsByName verifies the related behavior.
func TestRegistryListReturnsChannelsByName(t *testing.T) {
	reg := NewRegistry()
	reg.Register(nil)
	reg.Register(registryTestChannel{name: "stripe"})
	reg.Register(registryTestChannel{name: "alipay"})
	reg.Register(registryTestChannel{name: "wechat"})

	channels := reg.List()
	if len(channels) != 3 {
		t.Fatalf("len(List()) = %d, want 3", len(channels))
	}

	got := []string{channels[0].Name(), channels[1].Name(), channels[2].Name()}
	want := []string{"alipay", "stripe", "wechat"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("List() = %v, want %v", got, want)
		}
	}
}
