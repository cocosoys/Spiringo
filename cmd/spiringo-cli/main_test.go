package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 中文：TestRunGenerateModuleAcceptsTrailingFlags 验证相关行为符合预期。
// English: TestRunGenerateModuleAcceptsTrailingFlags verifies the related behavior.
func TestRunGenerateModuleAcceptsTrailingFlags(t *testing.T) {
	root := t.TempDir()

	if err := run([]string{"generate", "module", "order-item", "-out", root}); err != nil {
		t.Fatalf("generate module: %v", err)
	}

	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order_item", "order_item.go"),
		"func NewOrderItemModule() *OrderItemModule",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order_item", "handler", "order_item.go"),
		"types.OKWithPage",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order_item", "handler", "order_item.go"),
		"req.GetPage(), req.GetPageSize()",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order_item", "dto", "order_item.go"),
		"func (r ListReq) GetPageSize() int",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order_item", "repository", "order_item.go"),
		"func (r *OrderItemRepository) Create",
	)
}

// 中文：TestRunGenerateDoesNotOverwriteWithoutForce 验证相关行为符合预期。
// English: TestRunGenerateDoesNotOverwriteWithoutForce verifies the related behavior.
func TestRunGenerateDoesNotOverwriteWithoutForce(t *testing.T) {
	root := t.TempDir()
	args := []string{"generate", "oauth-provider", "github", "-out", root}
	if err := run(args); err != nil {
		t.Fatalf("generate oauth provider: %v", err)
	}
	if err := run(args); err == nil {
		t.Fatalf("expected overwrite error")
	}
	if err := run(append(args, "-force")); err != nil {
		t.Fatalf("generate with force: %v", err)
	}
}

// 中文：TestGenerateProviderAndPaymentChannelMatchCurrentInterfaces 验证相关行为符合预期。
// English: TestGenerateProviderAndPaymentChannelMatchCurrentInterfaces verifies the related behavior.
func TestGenerateProviderAndPaymentChannelMatchCurrentInterfaces(t *testing.T) {
	root := t.TempDir()

	if err := run([]string{"generate", "payment-channel", "bank-pay", "--out=" + root}); err != nil {
		t.Fatalf("generate payment channel: %v", err)
	}
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"func (c *BankPayChannel) QueryPayment",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"func (c *BankPayChannel) ClosePayment",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"func (c *BankPayChannel) CallbackSuccess",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"func (c *BankPayChannel) CallbackFail",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"TradeNo:  gatewayFirst",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"QrCode:   gatewayFirst",
	)
	assertFileNotContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"not configured",
	)
	assertFileNotContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"QRCode:",
	)
	assertFileNotContains(t,
		filepath.Join(root, "internal", "modules", "payment", "channel", "bank_pay.go"),
		"TransactionID:",
	)

	if err := run([]string{"generate", "oauth-provider", "github", "--out=" + root}); err != nil {
		t.Fatalf("generate oauth provider: %v", err)
	}
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) GetUserInfo",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) GetUserInfoWithRedirect",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) AuthURL",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) GetUser",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) RefreshToken",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		`form.Set("grant_type", "refresh_token")`,
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) refreshAccessToken",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func parseGithubProviderToken",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"func (p *GithubProvider) exchangeToken",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		`form.Set("redirect_uri", redirectURL)`,
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"RawData:     raw",
	)
	assertFileNotContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"oauth provider is not configured",
	)
	assertFileNotContains(t,
		filepath.Join(root, "internal", "modules", "auth", "oauth", "github.go"),
		"ErrRefreshTokenUnsupported",
	)
}

// 中文：TestDirectModuleCommandMatchesBlueprintSyntax 验证相关行为符合预期。
// English: TestDirectModuleCommandMatchesBlueprintSyntax verifies the related behavior.
func TestDirectModuleCommandMatchesBlueprintSyntax(t *testing.T) {
	root := t.TempDir()

	if err := run([]string{"module", "inventory", "-out", root}); err != nil {
		t.Fatalf("module command: %v", err)
	}
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "inventory", "inventory.go"),
		"func NewInventoryModule() *InventoryModule",
	)
}

// 中文：TestModuleCommandUsesTargetModulePathAndUpdatesBuiltinRegistry 验证相关行为符合预期。
// English: TestModuleCommandUsesTargetModulePathAndUpdatesBuiltinRegistry verifies the related behavior.
func TestModuleCommandUsesTargetModulePathAndUpdatesBuiltinRegistry(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/acme/shop\n\ngo 1.25.1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	builtinDir := filepath.Join(root, "internal", "modules", "builtin")
	if err := os.MkdirAll(builtinDir, 0o755); err != nil {
		t.Fatal(err)
	}
	builtinPath := filepath.Join(builtinDir, "builtin.go")
	if err := os.WriteFile(builtinPath, []byte(`package builtin

import (
	"example.com/acme/shop/internal/core/app"
)

func RegisterAll(application *app.App) {
	application.RegisterModules(
	)
}
`), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := run([]string{"module", "inventory", "-out", root}); err != nil {
		t.Fatalf("module command: %v", err)
	}

	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "inventory", "inventory.go"),
		`"example.com/acme/shop/internal/modules/inventory/handler"`,
	)
	assertFileContains(t, builtinPath, `"example.com/acme/shop/internal/modules/inventory"`)
	assertFileContains(t, builtinPath, "inventory.NewInventoryModule()")
}

// 中文：TestNewCommandCopiesTemplateAndRewritesModule 验证相关行为符合预期。
// English: TestNewCommandCopiesTemplateAndRewritesModule verifies the related behavior.
func TestNewCommandCopiesTemplateAndRewritesModule(t *testing.T) {
	template := t.TempDir()
	if err := os.MkdirAll(filepath.Join(template, "cmd", "spiringo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(template, "go.mod"), []byte("module github.com/spiringo/spiringo\n\ngo 1.25.1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(template, "README.md"), []byte("# Template\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(template, "cmd", "spiringo", "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(template, ".omx"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(template, ".omx", "state.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(template, ".gocache"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(template, ".gocache", "cache.bin"), []byte("cache"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := t.TempDir()
	if err := run([]string{"new", "Billing_App", "-out", root, "-module", "example.com/acme/billing", "-template", template}); err != nil {
		t.Fatalf("new command: %v", err)
	}

	project := filepath.Join(root, "billing-app")
	assertFileContains(t, filepath.Join(project, "go.mod"), "module example.com/acme/billing")
	assertFileContains(t, filepath.Join(project, "README.md"), "Module: `example.com/acme/billing`")
	assertFileContains(t, filepath.Join(project, "cmd", "spiringo", "main.go"), "package main")
	if _, err := os.Stat(filepath.Join(project, ".omx")); !os.IsNotExist(err) {
		t.Fatalf("expected .omx to be skipped, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(project, ".gocache")); !os.IsNotExist(err) {
		t.Fatalf("expected .gocache to be skipped, err=%v", err)
	}
}

// 中文：TestCRUDCommandGeneratesLayeredScaffold 验证相关行为符合预期。
// English: TestCRUDCommandGeneratesLayeredScaffold verifies the related behavior.
func TestCRUDCommandGeneratesLayeredScaffold(t *testing.T) {
	root := t.TempDir()

	if err := run([]string{"crud", "order", "Product", "-out", root}); err != nil {
		t.Fatalf("crud command: %v", err)
	}

	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order", "model", "product.go"),
		"type Product struct",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order", "repository", "product.go"),
		"func (r *ProductRepository) GetByID",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order", "handler", "product.go"),
		"func RegisterProductRoutes",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order", "handler", "product.go"),
		"req.GetPage(), req.GetPageSize()",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order", "dto", "product.go"),
		"func (r ProductListReq) GetPageSize() int",
	)
	assertFileContains(t,
		filepath.Join(root, "internal", "modules", "order", "migrations_product.go"),
		"func ProductMigrations(db *orm.DB) []module.Migration",
	)
}

// 中文：TestMigrateCreateCommandGeneratesMigration 验证相关行为符合预期。
// English: TestMigrateCreateCommandGeneratesMigration verifies the related behavior.
func TestMigrateCreateCommandGeneratesMigration(t *testing.T) {
	root := t.TempDir()

	if err := run([]string{"migrate", "create", "add-product-table", "-out", root}); err != nil {
		t.Fatalf("migrate create command: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(root, "internal", "migrations", "*_add_product_table.go"))
	if err != nil {
		t.Fatalf("glob migration: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one migration, got %d", len(matches))
	}
	assertFileContains(t, matches[0], "func NewAddProductTableMigration() module.Migration")
}

// 中文：TestPackageNameValidation 验证相关行为符合预期。
// English: TestPackageNameValidation verifies the related behavior.
func TestPackageNameValidation(t *testing.T) {
	got, err := packageName("Order-Item")
	if err != nil {
		t.Fatalf("packageName: %v", err)
	}
	if got != "order_item" {
		t.Fatalf("packageName = %q", got)
	}
	if _, err := packageName("1bad"); err == nil {
		t.Fatalf("expected invalid package name")
	}
}

// 中文：assertFileContains 执行当前包中的对应流程。
// English: assertFileContains executes the corresponding workflow in this package.
func assertFileContains(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), want) {
		t.Fatalf("%s does not contain %q", path, want)
	}
}

// 中文：assertFileNotContains 执行当前包中的对应流程。
// English: assertFileNotContains executes the corresponding workflow in this package.
func assertFileNotContains(t *testing.T, path, unwanted string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if strings.Contains(string(data), unwanted) {
		t.Fatalf("%s contains unwanted %q", path, unwanted)
	}
}
