package rbac

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/rbac/model"
)

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
// Migrations 返回数据库迁移
func (m *RBACModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "rbac_001_create_roles_table",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(
					&model.Role{},
					&model.Permission{},
					&model.RolePermission{},
					&model.UserRole{},
				)
			},
			Down: func(ctx context.Context) error {
				m.migrateDB.DB().Migrator().DropTable(&model.UserRole{})
				m.migrateDB.DB().Migrator().DropTable(&model.RolePermission{})
				m.migrateDB.DB().Migrator().DropTable(&model.Permission{})
				m.migrateDB.DB().Migrator().DropTable(&model.Role{})
				return nil
			},
		},
	}
}

// 中文：Subscriptions 执行当前包中的对应流程。
// English: Subscriptions executes the corresponding workflow in this package.
// Subscriptions 返回事件订阅
func (m *RBACModule) Subscriptions() []module.EventSubscription {
	return []module.EventSubscription{
		{
			Topic: event.EventTenantCreated,
			Handler: func(ctx context.Context, e *event.Event) error {
				payload, ok := e.Payload.(*event.TenantEventPayload)
				if !ok {
					return nil
				}
				// 租户创建时，创建默认角色
				if err := m.svc.CreateDefaultRoles(ctx, payload.TenantID); err != nil {
					return fmt.Errorf("rbac: create default roles for tenant %s: %w", payload.TenantID, err)
				}
				return nil
			},
		},
		{
			Topic: event.EventUserCreated,
			Handler: func(ctx context.Context, e *event.Event) error {
				payload, ok := e.Payload.(*event.UserEventPayload)
				if !ok {
					return nil
				}
				roleCode := "viewer"
				if payload.Username == "admin" {
					roleCode = "admin"
				}
				if payload.UserID != "" && payload.TenantID != "" {
					if err := m.svc.AssignDefaultRoleByCode(ctx, payload.TenantID, payload.UserID, roleCode); err != nil {
						return fmt.Errorf("rbac: assign default role for user %s: %w", payload.UserID, err)
					}
				}
				return nil
			},
		},
	}
}
