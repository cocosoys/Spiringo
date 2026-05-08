package payment

import (
	"context"

	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/payment/model"
)

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
func (m *PaymentModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "payment_001_create_orders_tables",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(
					&model.PaymentOrder{},
					&model.RefundOrder{},
					&model.CallbackLog{},
				)
			},
			Down: func(ctx context.Context) error {
				m.migrateDB.DB().Migrator().DropTable(&model.CallbackLog{})
				m.migrateDB.DB().Migrator().DropTable(&model.RefundOrder{})
				m.migrateDB.DB().Migrator().DropTable(&model.PaymentOrder{})
				return nil
			},
		},
	}
}
