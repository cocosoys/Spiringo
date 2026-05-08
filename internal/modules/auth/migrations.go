package auth

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/auth/model"
)

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
// Migrations 返回数据库迁移
func (m *AuthModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "auth_001_create_oauth_bindings_table",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(&model.OAuthBinding{})
			},
			Down: func(ctx context.Context) error {
				return m.migrateDB.DB().Migrator().DropTable(&model.OAuthBinding{})
			},
		},
	}
}

// 中文：Subscriptions 执行当前包中的对应流程。
// English: Subscriptions executes the corresponding workflow in this package.
// Subscriptions 返回事件订阅
func (m *AuthModule) Subscriptions() []module.EventSubscription {
	return []module.EventSubscription{
		{
			Topic: event.EventUserDeleted,
			Handler: func(ctx context.Context, e *event.Event) error {
				payload, ok := e.Payload.(*event.UserEventPayload)
				if !ok {
					return nil
				}
				// 用户删除时，清理该用户的所有OAuth绑定
				if payload.UserID != "" {
					if err := m.svc.CleanUserOAuthBindings(ctx, payload.UserID); err != nil {
						return fmt.Errorf("auth: clean oauth bindings for user %s: %w", payload.UserID, err)
					}
				}
				return nil
			},
		},
	}
}
