package builtin

import (
	"github.com/spiringo/spiringo/internal/core/app"
	"github.com/spiringo/spiringo/internal/modules/auth"
	"github.com/spiringo/spiringo/internal/modules/notification"
	"github.com/spiringo/spiringo/internal/modules/payment"
	"github.com/spiringo/spiringo/internal/modules/qrcode"
	"github.com/spiringo/spiringo/internal/modules/rbac"
	"github.com/spiringo/spiringo/internal/modules/tenant"
	"github.com/spiringo/spiringo/internal/modules/user"
)

// 中文：RegisterAll 执行当前包中的对应流程。
// English: RegisterAll executes the corresponding workflow in this package.
func RegisterAll(application *app.App) {
	application.RegisterModules(
		tenant.NewTenantModule(),
		user.NewUserModule(),
		auth.NewAuthModule(),
		rbac.NewRBACModule(),
		payment.NewPaymentModule(),
		qrcode.NewQRCodeModule(),
		notification.NewNotificationModule(),
	)
}
