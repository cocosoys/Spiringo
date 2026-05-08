package qrcode

import (
	"context"

	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/qrcode/model"
)

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
// Migrations 返回数据库迁移
func (m *QRCodeModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "qrcode_001_create_qrcode_tables",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(
					&model.QRCodeRecord{},
					&model.ScanLog{},
				)
			},
			Down: func(ctx context.Context) error {
				m.migrateDB.DB().Migrator().DropTable(&model.ScanLog{})
				m.migrateDB.DB().Migrator().DropTable(&model.QRCodeRecord{})
				return nil
			},
		},
	}
}

// 中文：Subscriptions 执行当前包中的对应流程。
// English: Subscriptions executes the corresponding workflow in this package.
// Subscriptions 返回事件订阅
func (m *QRCodeModule) Subscriptions() []module.EventSubscription {
	return []module.EventSubscription{}
}
