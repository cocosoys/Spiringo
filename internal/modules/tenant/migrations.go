package tenant

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/tenant/model"
)

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
// Migrations 返回数据库迁移
func (m *TenantModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "tenant_001_create_tenants_table",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(&model.Tenant{})
			},
			Down: func(ctx context.Context) error {
				return m.migrateDB.DB().Migrator().DropTable(&model.Tenant{})
			},
		},
	}
}

// 中文：Subscriptions 执行当前包中的对应流程。
// English: Subscriptions executes the corresponding workflow in this package.
// Subscriptions 返回事件订阅
func (m *TenantModule) Subscriptions() []module.EventSubscription {
	return []module.EventSubscription{
		{
			Topic: event.EventUserCreated,
			Handler: func(ctx context.Context, e *event.Event) error {
				payload, ok := e.Payload.(*event.UserEventPayload)
				if !ok {
					return nil
				}
				// 用户创建时，确保其所属租户状态正常
				if payload.TenantID != "" {
					tenant, err := m.svc.GetByID(ctx, payload.TenantID)
					if err != nil {
						return fmt.Errorf("tenant: check tenant on user created: %w", err)
					}
					if tenant.Status != "active" {
						return fmt.Errorf("tenant: tenant %s is not active", payload.TenantID)
					}
				}
				return nil
			},
		},
	}
}
