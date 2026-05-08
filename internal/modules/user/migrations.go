package user

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/user/model"
)

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
// Migrations 返回数据库迁移
func (m *UserModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "user_001_create_users_table",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(&model.User{})
			},
			Down: func(ctx context.Context) error {
				return m.migrateDB.DB().Migrator().DropTable(&model.User{})
			},
		},
	}
}

// 中文：Subscriptions 执行当前包中的对应流程。
// English: Subscriptions executes the corresponding workflow in this package.
// Subscriptions 返回事件订阅
func (m *UserModule) Subscriptions() []module.EventSubscription {
	return []module.EventSubscription{
		{
			Topic: event.EventTenantCreated,
			Handler: func(ctx context.Context, e *event.Event) error {
				payload, ok := e.Payload.(*event.TenantEventPayload)
				if !ok {
					return nil
				}
				// 租户创建时，创建默认管理员
				if err := m.svc.CreateAdminForTenant(ctx, payload.TenantID); err != nil {
					return fmt.Errorf("user: create default admin for tenant %s: %w", payload.TenantID, err)
				}
				return nil
			},
		},
	}
}
