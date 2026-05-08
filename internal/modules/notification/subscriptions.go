package notification

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/notification/model"
	"github.com/spiringo/spiringo/internal/modules/notification/service"
)

// 中文：Subscriptions 执行当前包中的对应流程。
// English: Subscriptions executes the corresponding workflow in this package.
func (m *Module) Subscriptions() []module.EventSubscription {
	return []module.EventSubscription{
		{
			Topic:   event.EventPaymentFailed,
			Handler: m.handlePaymentFailed,
		},
		{
			Topic:   event.EventPaymentSuccess,
			Handler: m.handlePaymentSuccess,
		},
		{
			Topic:   event.EventTenantSuspended,
			Handler: m.handleTenantSuspended,
		},
	}
}

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
func (m *Module) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "notification_001_create_messages_table",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(&model.Notification{})
			},
			Down: func(ctx context.Context) error {
				m.migrateDB.DB().Migrator().DropTable(&model.Notification{})
				return nil
			},
		},
	}
}

// 中文：handlePaymentFailed 执行当前包中的对应流程。
// English: handlePaymentFailed executes the corresponding workflow in this package.
func (m *Module) handlePaymentFailed(ctx context.Context, e *event.Event) error {
	payload, ok := e.Payload.(*event.PaymentEventPayload)
	if !ok {
		return nil
	}
	return m.svc.Notify(ctx, service.Message{
		Event:    event.EventPaymentFailed,
		Severity: "error",
		Subject:  "Payment failed",
		Content:  fmt.Sprintf("Payment order %s failed on channel %s.", payload.OutTradeNo, payload.Channel),
		TenantID: payload.TenantID,
		Payload: map[string]any{
			"out_trade_no": payload.OutTradeNo,
			"trade_no":     payload.TradeNo,
			"channel":      payload.Channel,
			"amount":       payload.Amount,
		},
	})
}

// 中文：handlePaymentSuccess 执行当前包中的对应流程。
// English: handlePaymentSuccess executes the corresponding workflow in this package.
func (m *Module) handlePaymentSuccess(ctx context.Context, e *event.Event) error {
	payload, ok := e.Payload.(*event.PaymentEventPayload)
	if !ok {
		return nil
	}
	return m.svc.Notify(ctx, service.Message{
		Event:    event.EventPaymentSuccess,
		Severity: "info",
		Subject:  "Payment succeeded",
		Content:  fmt.Sprintf("Payment order %s succeeded on channel %s.", payload.OutTradeNo, payload.Channel),
		TenantID: payload.TenantID,
		Payload: map[string]any{
			"out_trade_no": payload.OutTradeNo,
			"trade_no":     payload.TradeNo,
			"channel":      payload.Channel,
			"amount":       payload.Amount,
		},
	})
}

// 中文：handleTenantSuspended 执行当前包中的对应流程。
// English: handleTenantSuspended executes the corresponding workflow in this package.
func (m *Module) handleTenantSuspended(ctx context.Context, e *event.Event) error {
	payload, ok := e.Payload.(*event.TenantEventPayload)
	if !ok {
		return nil
	}
	return m.svc.Notify(ctx, service.Message{
		Event:    event.EventTenantSuspended,
		Severity: "warning",
		Subject:  "Tenant suspended",
		Content:  fmt.Sprintf("Tenant %s has been suspended.", payload.TenantName),
		TenantID: payload.TenantID,
		Payload: map[string]any{
			"tenant_id":   payload.TenantID,
			"tenant_name": payload.TenantName,
			"strategy":    payload.Strategy,
		},
	})
}
