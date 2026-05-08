package main

import (
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

// 中文：generator 定义当前包使用的数据结构或接口。
// English: generator defines a data structure or interface used by this package.
type generator struct {
	// 中文：root 保存当前结构中的配置或数据值。
	// English: root stores a configuration or data value for this struct.
	root string
	// 中文：force 保存当前结构中的配置或数据值。
	// English: force stores a configuration or data value for this struct.
	force bool
}

// 中文：cliOptions 定义当前包使用的数据结构或接口。
// English: cliOptions defines a data structure or interface used by this package.
type cliOptions struct {
	// 中文：root 保存当前结构中的配置或数据值。
	// English: root stores a configuration or data value for this struct.
	root string
	// 中文：force 保存当前结构中的配置或数据值。
	// English: force stores a configuration or data value for this struct.
	force bool
	// 中文：modulePath 保存当前结构中的配置或数据值。
	// English: modulePath stores a configuration or data value for this struct.
	modulePath string
	// 中文：templateRoot 保存当前结构中的配置或数据值。
	// English: templateRoot stores a configuration or data value for this struct.
	templateRoot string
	// 中文：rest 保存当前结构中的配置或数据值。
	// English: rest stores a configuration or data value for this struct.
	rest []string
}

// 中文：main 是当前命令的程序入口。
// English: main is the entry point for this command.
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// 中文：run 执行当前包中的对应流程。
// English: run executes the corresponding workflow in this package.
func run(args []string) error {
	if len(args) == 0 {
		return usageError("")
	}

	switch args[0] {
	case "generate", "gen":
		return runGenerate(args[1:])
	case "new":
		return runNew(args[1:])
	case "module":
		return runDirectGenerate("module", args[1:])
	case "payment-channel":
		return runDirectGenerate("payment-channel", args[1:])
	case "oauth-provider":
		return runDirectGenerate("oauth-provider", args[1:])
	case "crud":
		return runCRUD(args[1:])
	case "migrate":
		return runMigrate(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return usageError("unknown command: " + args[0])
	}
}

// 中文：runNew 执行当前包中的对应流程。
// English: runNew executes the corresponding workflow in this package.
func runNew(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if len(opts.rest) != 1 {
		return usageError("new requires: <project-name>")
	}
	return generator{root: opts.root, force: opts.force}.newProject(opts.rest[0], opts.modulePath, opts.templateRoot)
}

// 中文：runGenerate 执行当前包中的对应流程。
// English: runGenerate executes the corresponding workflow in this package.
func runGenerate(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}

	if len(opts.rest) != 2 {
		return usageError("generate requires: <module|payment-channel|oauth-provider|crud> <name>")
	}

	g := generator{root: opts.root, force: opts.force}
	kind, name := opts.rest[0], opts.rest[1]
	switch kind {
	case "module":
		return g.module(name)
	case "payment-channel":
		return g.paymentChannel(name)
	case "oauth-provider":
		return g.oauthProvider(name)
	default:
		return usageError("unknown generator kind: " + kind)
	}
}

// 中文：runDirectGenerate 执行当前包中的对应流程。
// English: runDirectGenerate executes the corresponding workflow in this package.
func runDirectGenerate(kind string, args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if len(opts.rest) != 1 {
		return usageError(kind + " requires: <name>")
	}
	g := generator{root: opts.root, force: opts.force}
	switch kind {
	case "module":
		return g.module(opts.rest[0])
	case "payment-channel":
		return g.paymentChannel(opts.rest[0])
	case "oauth-provider":
		return g.oauthProvider(opts.rest[0])
	default:
		return usageError("unknown generator kind: " + kind)
	}
}

// 中文：runCRUD 执行当前包中的对应流程。
// English: runCRUD executes the corresponding workflow in this package.
func runCRUD(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if len(opts.rest) != 2 {
		return usageError("crud requires: <module> <model>")
	}
	return generator{root: opts.root, force: opts.force}.crud(opts.rest[0], opts.rest[1])
}

// 中文：runMigrate 执行当前包中的对应流程。
// English: runMigrate executes the corresponding workflow in this package.
func runMigrate(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if len(opts.rest) != 2 || opts.rest[0] != "create" {
		return usageError("migrate requires: create <name>")
	}
	return generator{root: opts.root, force: opts.force}.migration(opts.rest[1])
}

// 中文：parseOptions 执行当前包中的对应流程。
// English: parseOptions executes the corresponding workflow in this package.
func parseOptions(args []string) (cliOptions, error) {
	opts := cliOptions{root: ".", rest: make([]string, 0, len(args))}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-force" || arg == "--force":
			opts.force = true
		case arg == "-out" || arg == "--out":
			if i+1 >= len(args) {
				return opts, usageError(arg + " requires a value")
			}
			i++
			opts.root = args[i]
		case strings.HasPrefix(arg, "-out="):
			opts.root = strings.TrimPrefix(arg, "-out=")
		case strings.HasPrefix(arg, "--out="):
			opts.root = strings.TrimPrefix(arg, "--out=")
		case arg == "-module" || arg == "--module":
			if i+1 >= len(args) {
				return opts, usageError(arg + " requires a value")
			}
			i++
			opts.modulePath = args[i]
		case strings.HasPrefix(arg, "-module="):
			opts.modulePath = strings.TrimPrefix(arg, "-module=")
		case strings.HasPrefix(arg, "--module="):
			opts.modulePath = strings.TrimPrefix(arg, "--module=")
		case arg == "-template" || arg == "--template":
			if i+1 >= len(args) {
				return opts, usageError(arg + " requires a value")
			}
			i++
			opts.templateRoot = args[i]
		case strings.HasPrefix(arg, "-template="):
			opts.templateRoot = strings.TrimPrefix(arg, "-template=")
		case strings.HasPrefix(arg, "--template="):
			opts.templateRoot = strings.TrimPrefix(arg, "--template=")
		case strings.HasPrefix(arg, "-"):
			return opts, usageError("unknown flag: " + arg)
		default:
			opts.rest = append(opts.rest, arg)
		}
	}
	return opts, nil
}

// 中文：newProject 执行当前包中的对应流程。
// English: newProject executes the corresponding workflow in this package.
func (g generator) newProject(name, modulePath, templateRoot string) error {
	slug, err := projectSlug(name)
	if err != nil {
		return err
	}
	if strings.TrimSpace(templateRoot) == "" {
		templateRoot = "."
	}
	if strings.TrimSpace(modulePath) == "" {
		modulePath = "github.com/spiringo/" + slug
	}

	src, err := filepath.Abs(templateRoot)
	if err != nil {
		return err
	}
	dst, err := filepath.Abs(filepath.Join(g.root, slug))
	if err != nil {
		return err
	}
	if !g.force {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("project already exists: %s", dst)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	if err := copyTemplateTree(src, dst, modulePath, slug); err != nil {
		return err
	}
	fmt.Println("created project", dst)
	return nil
}

// 中文：module 执行当前包中的对应流程。
// English: module executes the corresponding workflow in this package.
func (g generator) module(name string) error {
	pkg, err := packageName(name)
	if err != nil {
		return err
	}
	modulePath, err := g.modulePath()
	if err != nil {
		return err
	}
	typ := exportName(pkg)
	dir := filepath.Join(g.root, "internal", "modules", pkg)
	data := map[string]string{
		filepath.Join(dir, pkg+".go"):               fmt.Sprintf(moduleTemplate, modulePath, pkg, typ),
		filepath.Join(dir, "migrations.go"):         fmt.Sprintf(moduleMigrationsTemplate, modulePath, pkg, typ),
		filepath.Join(dir, "handler", pkg+".go"):    fmt.Sprintf(moduleHandlerTemplate, modulePath, pkg, typ),
		filepath.Join(dir, "service", pkg+".go"):    fmt.Sprintf(moduleServiceTemplate, modulePath, pkg, typ),
		filepath.Join(dir, "repository", pkg+".go"): fmt.Sprintf(moduleRepositoryTemplate, modulePath, pkg, typ),
		filepath.Join(dir, "model", pkg+".go"):      fmt.Sprintf(moduleModelTemplate, modulePath, pkg, typ),
		filepath.Join(dir, "dto", pkg+".go"):        moduleDTOTemplate,
		filepath.Join(dir, "README.md"):             fmt.Sprintf(moduleReadmeTemplate, pkg, typ),
	}
	if err := g.writeFiles(data); err != nil {
		return err
	}
	if err := g.registerBuiltinModule(modulePath, pkg, typ); err != nil {
		return err
	}
	return nil
}

// 中文：crud 执行当前包中的对应流程。
// English: crud executes the corresponding workflow in this package.
func (g generator) crud(moduleName, modelName string) error {
	modPkg, err := packageName(moduleName)
	if err != nil {
		return err
	}
	modulePath, err := g.modulePath()
	if err != nil {
		return err
	}
	modelPkg, err := packageName(modelName)
	if err != nil {
		return err
	}
	modelType := exportName(modelPkg)
	tableName := modPkg + "_" + modelPkg
	base := filepath.Join(g.root, "internal", "modules", modPkg)
	data := map[string]string{
		filepath.Join(base, "model", modelPkg+".go"):      fmt.Sprintf(crudModelTemplate, modulePath, tableName, modelType),
		filepath.Join(base, "dto", modelPkg+".go"):        fmt.Sprintf(crudDTOTemplate, modelType),
		filepath.Join(base, "repository", modelPkg+".go"): fmt.Sprintf(crudRepositoryTemplate, modulePath, modPkg, modelType),
		filepath.Join(base, "service", modelPkg+".go"):    fmt.Sprintf(crudServiceTemplate, modulePath, modPkg, modelType),
		filepath.Join(base, "handler", modelPkg+".go"):    fmt.Sprintf(crudHandlerTemplate, modulePath, modPkg, modelType, modelPkg),
		filepath.Join(base, "migrations_"+modelPkg+".go"): fmt.Sprintf(crudMigrationTemplate, modulePath, modPkg, modelType, modelPkg),
	}
	return g.writeFiles(data)
}

// 中文：paymentChannel 执行当前包中的对应流程。
// English: paymentChannel executes the corresponding workflow in this package.
func (g generator) paymentChannel(name string) error {
	pkg, err := packageName(name)
	if err != nil {
		return err
	}
	typ := exportName(pkg) + "Channel"
	dir := filepath.Join(g.root, "internal", "modules", "payment", "channel")
	return g.writeFiles(map[string]string{
		filepath.Join(dir, pkg+".go"): fmt.Sprintf(paymentChannelTemplate, typ, typ, pkg, typ, typ, typ, typ, typ),
	})
}

// 中文：oauthProvider 执行当前包中的对应流程。
// English: oauthProvider executes the corresponding workflow in this package.
func (g generator) oauthProvider(name string) error {
	pkg, err := packageName(name)
	if err != nil {
		return err
	}
	typ := exportName(pkg) + "Provider"
	dir := filepath.Join(g.root, "internal", "modules", "auth", "oauth")
	return g.writeFiles(map[string]string{
		filepath.Join(dir, pkg+".go"): fmt.Sprintf(oauthProviderTemplate, typ, typ, pkg),
	})
}

// 中文：migration 执行当前包中的对应流程。
// English: migration executes the corresponding workflow in this package.
func (g generator) migration(name string) error {
	pkg, err := packageName(name)
	if err != nil {
		return err
	}
	stamp := time.Now().UTC().Format("20060102150405")
	id := stamp + "_" + pkg
	typ := exportName(pkg)
	path := filepath.Join(g.root, "internal", "migrations", id+".go")
	return g.writeFiles(map[string]string{
		path: fmt.Sprintf(migrationTemplate, id, typ),
	})
}

// 中文：writeFiles 执行当前包中的对应流程。
// English: writeFiles executes the corresponding workflow in this package.
func (g generator) writeFiles(files map[string]string) error {
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		content := files[path]
		if !g.force {
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("file already exists: %s", path)
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
		fmt.Println("created", path)
	}
	return nil
}

// 中文：modulePath 执行当前包中的对应流程。
// English: modulePath executes the corresponding workflow in this package.
func (g generator) modulePath() (string, error) {
	goMod := filepath.Join(g.root, "go.mod")
	data, err := os.ReadFile(goMod)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "github.com/spiringo/spiringo", nil
		}
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if value != "" {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("go.mod does not declare module path")
}

// 中文：registerBuiltinModule 执行当前包中的对应流程。
// English: registerBuiltinModule executes the corresponding workflow in this package.
func (g generator) registerBuiltinModule(modulePath, pkg, typ string) error {
	path := filepath.Join(g.root, "internal", "modules", "builtin", "builtin.go")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	content := string(data)
	importPath := fmt.Sprintf("%s/internal/modules/%s", modulePath, pkg)
	if strings.Contains(content, `"`+importPath+`"`) {
		return nil
	}

	content, err = insertImport(content, importPath)
	if err != nil {
		return fmt.Errorf("register module import: %w", err)
	}
	content, err = insertModuleRegistration(content, pkg, typ)
	if err != nil {
		return fmt.Errorf("register module constructor: %w", err)
	}
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("format builtin module registration: %w", err)
	}
	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		return err
	}
	fmt.Println("updated", path)
	return nil
}

// 中文：insertImport 执行当前包中的对应流程。
// English: insertImport executes the corresponding workflow in this package.
func insertImport(content, importPath string) (string, error) {
	marker := "import (\n"
	idx := strings.Index(content, marker)
	if idx < 0 {
		return "", fmt.Errorf("import block not found")
	}
	insertAt := idx + len(marker)
	return content[:insertAt] + fmt.Sprintf("\t%q\n", importPath) + content[insertAt:], nil
}

// 中文：insertModuleRegistration 执行当前包中的对应流程。
// English: insertModuleRegistration executes the corresponding workflow in this package.
func insertModuleRegistration(content, pkg, typ string) (string, error) {
	call := "application.RegisterModules(\n"
	idx := strings.Index(content, call)
	if idx < 0 {
		return "", fmt.Errorf("RegisterModules call not found")
	}
	start := idx + len(call)
	offset := 0
	for _, line := range strings.SplitAfter(content[start:], "\n") {
		if strings.TrimSpace(line) == ")" {
			insertAt := start + offset
			registration := fmt.Sprintf("\t\t%s.New%sModule(),\n", pkg, typ)
			return content[:insertAt] + registration + content[insertAt:], nil
		}
		offset += len(line)
	}
	return "", fmt.Errorf("RegisterModules closing marker not found")
}

// 中文：packageName 执行当前包中的对应流程。
// English: packageName executes the corresponding workflow in this package.
func packageName(raw string) (string, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return "", fmt.Errorf("name is required")
	}
	raw = strings.ReplaceAll(raw, "-", "_")
	if !regexp.MustCompile(`^[a-z][a-z0-9_]*$`).MatchString(raw) {
		return "", fmt.Errorf("invalid name %q: use letters, numbers, hyphen or underscore, starting with a letter", raw)
	}
	return raw, nil
}

// 中文：projectSlug 执行当前包中的对应流程。
// English: projectSlug executes the corresponding workflow in this package.
func projectSlug(raw string) (string, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return "", fmt.Errorf("project name is required")
	}
	raw = strings.ReplaceAll(raw, "_", "-")
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(raw) {
		return "", fmt.Errorf("invalid project name %q: use letters, numbers, hyphen or underscore, starting with a letter", raw)
	}
	return raw, nil
}

// 中文：copyTemplateTree 执行当前包中的对应流程。
// English: copyTemplateTree executes the corresponding workflow in this package.
func copyTemplateTree(src, dst, modulePath, slug string) error {
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if shouldSkipTemplatePath(src, dst, absPath, entry) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(src, absPath)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(absPath)
		if err != nil {
			return err
		}
		switch filepath.ToSlash(rel) {
		case "go.mod":
			data = rewriteGoModModule(data, modulePath)
		case "README.md":
			data = rewriteProjectReadme(data, slug, modulePath)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// 中文：shouldSkipTemplatePath 执行当前包中的对应流程。
// English: shouldSkipTemplatePath executes the corresponding workflow in this package.
func shouldSkipTemplatePath(src, dst, path string, entry os.DirEntry) bool {
	if path == src {
		return false
	}
	if sameOrInside(path, dst) {
		return true
	}
	name := entry.Name()
	if entry.IsDir() {
		switch name {
		case ".git", ".idea", ".omx", ".qwen", ".gocache", "bin", "dist", "tmp", "node_modules", "vendor":
			return true
		}
		return false
	}
	if strings.HasSuffix(name, ".exe") || strings.HasSuffix(name, ".test") {
		return true
	}
	switch name {
	case "coverage.txt", ".DS_Store":
		return true
	default:
		return false
	}
}

// 中文：sameOrInside 执行当前包中的对应流程。
// English: sameOrInside executes the corresponding workflow in this package.
func sameOrInside(path, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

// 中文：rewriteGoModModule 执行当前包中的对应流程。
// English: rewriteGoModModule executes the corresponding workflow in this package.
func rewriteGoModModule(data []byte, modulePath string) []byte {
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "module ") {
			lines[i] = "module " + modulePath
			return []byte(strings.Join(lines, "\n"))
		}
	}
	return []byte("module " + modulePath + "\n\n" + string(data))
}

// 中文：rewriteProjectReadme 执行当前包中的对应流程。
// English: rewriteProjectReadme executes the corresponding workflow in this package.
func rewriteProjectReadme(data []byte, slug, modulePath string) []byte {
	content := strings.TrimSpace(string(data))
	if content == "" {
		return []byte(fmt.Sprintf("# %s\n\nGenerated from the Spiringo template.\n\nModule: `%s`\n", slug, modulePath))
	}
	return []byte(fmt.Sprintf("# %s\n\nModule: `%s`\n\n%s\n", slug, modulePath, content))
}

// 中文：exportName 执行当前包中的对应流程。
// English: exportName executes the corresponding workflow in this package.
func exportName(pkg string) string {
	parts := strings.FieldsFunc(pkg, func(r rune) bool { return r == '_' || r == '-' })
	var b strings.Builder
	for _, part := range parts {
		for i, r := range part {
			if i == 0 {
				b.WriteRune(unicode.ToUpper(r))
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

// 中文：usageError 执行当前包中的对应流程。
// English: usageError executes the corresponding workflow in this package.
func usageError(msg string) error {
	if msg != "" {
		return fmt.Errorf("%s\n\n%s", msg, usageText)
	}
	return errors.New(usageText)
}

// 中文：printUsage 执行当前包中的对应流程。
// English: printUsage executes the corresponding workflow in this package.
func printUsage() {
	fmt.Fprintln(os.Stdout, usageText)
}

// 中文：usageText 声明当前包使用的常量。
// English: usageText declares constants used by this package.
const usageText = `Usage:
  spiringo-cli new <project-name> [-out <directory>] [-module <module-path>] [-template <template-root>] [-force]
  spiringo-cli module <name> [-out <project-root>] [-force]
  spiringo-cli crud <module> <model> [-out <project-root>] [-force]
  spiringo-cli migrate create <name> [-out <project-root>] [-force]
  spiringo-cli payment-channel <name> [-out <project-root>] [-force]
  spiringo-cli oauth-provider <name> [-out <project-root>] [-force]

Compatibility:
  spiringo-cli generate module <name> [-out <project-root>] [-force]
  spiringo-cli generate payment-channel <name> [-out <project-root>] [-force]
  spiringo-cli generate oauth-provider <name> [-out <project-root>] [-force]`

// 中文：moduleTemplate 声明当前包使用的常量。
// English: moduleTemplate declares constants used by this package.
const moduleTemplate = `package %[2]s

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"%[1]s/internal/core/di"
	"%[1]s/internal/core/module"
	"%[1]s/internal/modules/%[2]s/handler"
	"%[1]s/internal/modules/%[2]s/repository"
	"%[1]s/internal/modules/%[2]s/service"
	"%[1]s/internal/pkg/orm"
)

type %[3]sModule struct {
	*module.BaseModule
	svc       *service.%[3]sService
	migrateDB *orm.DB
}

func New%[3]sModule() *%[3]sModule {
	return &%[3]sModule{BaseModule: module.NewBaseModule("%[2]s", "tenant")}
}

func (m *%[3]sModule) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("%[2]s module init: %%w", err)
	}
	m.migrateDB = db
	repo := repository.New%[3]sRepository(orm.NewTenantDB(db), db)
	m.svc = service.New%[3]sService(repo)
	return nil
}

func (m *%[3]sModule) Routes(r *gin.RouterGroup) {
	h := handler.New%[3]sHandler(m.svc)
	r.GET("", h.List)
	r.POST("", h.Create)
}

func (m *%[3]sModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("%[2]s service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("%[2]s migration database is not initialized")
	}
	return nil
}

func (m *%[3]sModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}
`

// 中文：moduleMigrationsTemplate 声明当前包使用的常量。
// English: moduleMigrationsTemplate declares constants used by this package.
const moduleMigrationsTemplate = `package %[2]s

import (
	"context"

	"%[1]s/internal/core/module"
	"%[1]s/internal/modules/%[2]s/model"
)

func (m *%[3]sModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "%[2]s_001_create_table",
			Up: func(ctx context.Context) error {
				return m.migrateDB.AutoMigrate(&model.%[3]s{})
			},
			Down: func(ctx context.Context) error {
				return m.migrateDB.DB().Migrator().DropTable(&model.%[3]s{})
			},
		},
	}
}
`

// 中文：moduleHandlerTemplate 声明当前包使用的常量。
// English: moduleHandlerTemplate declares constants used by this package.
const moduleHandlerTemplate = `package handler

import (
	"github.com/gin-gonic/gin"
	"%[1]s/internal/modules/%[2]s/dto"
	"%[1]s/internal/modules/%[2]s/service"
	"%[1]s/pkg/types"
)

type %[3]sHandler struct {
	svc *service.%[3]sService
}

func New%[3]sHandler(svc *service.%[3]sService) *%[3]sHandler {
	return &%[3]sHandler{svc: svc}
}

func (h *%[3]sHandler) List(c *gin.Context) {
	var req dto.ListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}
	items, total, err := h.svc.List(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OKWithPage(c, items, total, req.GetPage(), req.GetPageSize())
}

func (h *%[3]sHandler) Create(c *gin.Context) {
	var req dto.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}
	item, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, item)
}
`

// 中文：moduleServiceTemplate 声明当前包使用的常量。
// English: moduleServiceTemplate declares constants used by this package.
const moduleServiceTemplate = `package service

import (
	"context"

	"%[1]s/internal/modules/%[2]s/dto"
	"%[1]s/internal/modules/%[2]s/model"
	"%[1]s/internal/modules/%[2]s/repository"
	"%[1]s/pkg/types"
)

type %[3]sService struct {
	repo *repository.%[3]sRepository
}

func New%[3]sService(repo *repository.%[3]sRepository) *%[3]sService {
	return &%[3]sService{repo: repo}
}

func (s *%[3]sService) List(ctx context.Context, req dto.ListReq) ([]model.%[3]s, int64, error) {
	return s.repo.List(ctx, req.GetPage(), req.GetPageSize())
}

func (s *%[3]sService) Create(ctx context.Context, req dto.CreateReq) (*model.%[3]s, error) {
	if req.Name == "" {
		return nil, types.ErrBadRequest.WithMessage("name is required")
	}
	item := &model.%[3]s{Name: req.Name}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
`

// 中文：moduleRepositoryTemplate 声明当前包使用的常量。
// English: moduleRepositoryTemplate declares constants used by this package.
const moduleRepositoryTemplate = `package repository

import (
	"context"

	"%[1]s/internal/modules/%[2]s/model"
	"%[1]s/internal/pkg/orm"
)

type %[3]sRepository struct {
	tdb *orm.TenantDB
	db  *orm.DB
}

func New%[3]sRepository(tdb *orm.TenantDB, db *orm.DB) *%[3]sRepository {
	return &%[3]sRepository{tdb: tdb, db: db}
}

func (r *%[3]sRepository) List(ctx context.Context, page, pageSize int) ([]model.%[3]s, int64, error) {
	var items []model.%[3]s
	total, err := r.tdb.Paginate(ctx, &items, page, pageSize)
	return items, total, err
}

func (r *%[3]sRepository) Create(ctx context.Context, item *model.%[3]s) error {
	return r.tdb.Create(ctx, item)
}
`

// 中文：moduleModelTemplate 声明当前包使用的常量。
// English: moduleModelTemplate declares constants used by this package.
const moduleModelTemplate = `package model

import "%[1]s/internal/pkg/orm"

type %[3]s struct {
	orm.TenantBaseModel
	Name string ` + "`gorm:\"size:128;not null\" json:\"name\"`" + `
}

func (%[3]s) TableName() string { return "%[2]s" }
`

// 中文：moduleDTOTemplate 声明当前包使用的常量。
// English: moduleDTOTemplate declares constants used by this package.
const moduleDTOTemplate = `package dto

type ListReq struct {
	Page     int ` + "`form:\"page\"`" + `
	PageSize int ` + "`form:\"page_size\"`" + `
}

func (r ListReq) GetPage() int {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

func (r ListReq) GetPageSize() int {
	if r.PageSize <= 0 {
		return 20
	}
	if r.PageSize > 100 {
		return 100
	}
	return r.PageSize
}

type CreateReq struct {
	Name string ` + "`json:\"name\" binding:\"required\"`" + `
}
`

// 中文：moduleReadmeTemplate 声明当前包使用的常量。
// English: moduleReadmeTemplate declares constants used by this package.
const moduleReadmeTemplate = `# %[2]s Module

Generated scaffold for the %[1]s module.

The CLI updates internal/modules/builtin/builtin.go automatically when that file exists.
If you use a custom entrypoint, register the module explicitly:

` + "```go" + `
module.RegisterModules(
	%[1]s.New%[2]sModule(),
)
` + "```" + `
`

// 中文：crudModelTemplate 声明当前包使用的常量。
// English: crudModelTemplate declares constants used by this package.
const crudModelTemplate = `package model

import "%[1]s/internal/pkg/orm"

type %[3]s struct {
	orm.TenantBaseModel
	Name        string ` + "`gorm:\"size:128;not null\" json:\"name\"`" + `
	Description string ` + "`gorm:\"size:512\" json:\"description,omitempty\"`" + `
}

func (%[3]s) TableName() string { return "%[2]s" }
`

// 中文：crudDTOTemplate 声明当前包使用的常量。
// English: crudDTOTemplate declares constants used by this package.
const crudDTOTemplate = `package dto

type %[1]sListReq struct {
	Page     int ` + "`form:\"page\"`" + `
	PageSize int ` + "`form:\"page_size\"`" + `
}

func (r %[1]sListReq) GetPage() int {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

func (r %[1]sListReq) GetPageSize() int {
	if r.PageSize <= 0 {
		return 20
	}
	if r.PageSize > 100 {
		return 100
	}
	return r.PageSize
}

type %[1]sCreateReq struct {
	Name        string ` + "`json:\"name\" binding:\"required\"`" + `
	Description string ` + "`json:\"description\"`" + `
}

type %[1]sUpdateReq struct {
	Name        *string ` + "`json:\"name\"`" + `
	Description *string ` + "`json:\"description\"`" + `
}
`

// 中文：crudRepositoryTemplate 声明当前包使用的常量。
// English: crudRepositoryTemplate declares constants used by this package.
const crudRepositoryTemplate = `package repository

import (
	"context"

	"%[1]s/internal/modules/%[2]s/model"
	"%[1]s/internal/pkg/orm"
)

type %[3]sRepository struct {
	tdb *orm.TenantDB
	db  *orm.DB
}

func New%[3]sRepository(tdb *orm.TenantDB, db *orm.DB) *%[3]sRepository {
	return &%[3]sRepository{tdb: tdb, db: db}
}

func (r *%[3]sRepository) List(ctx context.Context, page, pageSize int) ([]model.%[3]s, int64, error) {
	var items []model.%[3]s
	total, err := r.tdb.Paginate(ctx, &items, page, pageSize)
	return items, total, err
}

func (r *%[3]sRepository) GetByID(ctx context.Context, id string) (*model.%[3]s, error) {
	var item model.%[3]s
	if err := r.tdb.First(ctx, &item, "id = ?", id); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *%[3]sRepository) Create(ctx context.Context, item *model.%[3]s) error {
	return r.tdb.Create(ctx, item)
}

func (r *%[3]sRepository) Update(ctx context.Context, item *model.%[3]s) error {
	return r.tdb.Update(ctx, item)
}

func (r *%[3]sRepository) Delete(ctx context.Context, id string) error {
	var item model.%[3]s
	return r.tdb.Delete(ctx, &item, "id = ?", id)
}
`

// 中文：crudServiceTemplate 声明当前包使用的常量。
// English: crudServiceTemplate declares constants used by this package.
const crudServiceTemplate = `package service

import (
	"context"

	"%[1]s/internal/modules/%[2]s/dto"
	"%[1]s/internal/modules/%[2]s/model"
	"%[1]s/internal/modules/%[2]s/repository"
	"%[1]s/pkg/types"
)

type %[3]sService struct {
	repo *repository.%[3]sRepository
}

func New%[3]sService(repo *repository.%[3]sRepository) *%[3]sService {
	return &%[3]sService{repo: repo}
}

func (s *%[3]sService) List(ctx context.Context, req dto.%[3]sListReq) ([]model.%[3]s, int64, error) {
	return s.repo.List(ctx, req.GetPage(), req.GetPageSize())
}

func (s *%[3]sService) GetByID(ctx context.Context, id string) (*model.%[3]s, error) {
	if id == "" {
		return nil, types.ErrBadRequest.WithMessage("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *%[3]sService) Create(ctx context.Context, req dto.%[3]sCreateReq) (*model.%[3]s, error) {
	if req.Name == "" {
		return nil, types.ErrBadRequest.WithMessage("name is required")
	}
	item := &model.%[3]s{Name: req.Name, Description: req.Description}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *%[3]sService) Update(ctx context.Context, id string, req dto.%[3]sUpdateReq) (*model.%[3]s, error) {
	item, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *%[3]sService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return types.ErrBadRequest.WithMessage("id is required")
	}
	return s.repo.Delete(ctx, id)
}
`

// 中文：crudHandlerTemplate 声明当前包使用的常量。
// English: crudHandlerTemplate declares constants used by this package.
const crudHandlerTemplate = `package handler

import (
	"github.com/gin-gonic/gin"
	"%[1]s/internal/modules/%[2]s/dto"
	"%[1]s/internal/modules/%[2]s/service"
	"%[1]s/pkg/types"
)

type %[3]sHandler struct {
	svc *service.%[3]sService
}

func New%[3]sHandler(svc *service.%[3]sService) *%[3]sHandler {
	return &%[3]sHandler{svc: svc}
}

func Register%[3]sRoutes(r *gin.RouterGroup, svc *service.%[3]sService) {
	h := New%[3]sHandler(svc)
	g := r.Group("/%[4]s")
	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *%[3]sHandler) List(c *gin.Context) {
	var req dto.%[3]sListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}
	items, total, err := h.svc.List(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OKWithPage(c, items, total, req.GetPage(), req.GetPageSize())
}

func (h *%[3]sHandler) Get(c *gin.Context) {
	item, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, item)
}

func (h *%[3]sHandler) Create(c *gin.Context) {
	var req dto.%[3]sCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}
	item, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, item)
}

func (h *%[3]sHandler) Update(c *gin.Context) {
	var req dto.%[3]sUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}
	item, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, item)
}

func (h *%[3]sHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		types.Fail(c, err)
		return
	}
	types.OK(c, gin.H{"deleted": true})
}
`

// 中文：crudMigrationTemplate 声明当前包使用的常量。
// English: crudMigrationTemplate declares constants used by this package.
const crudMigrationTemplate = `package %[2]s

import (
	"context"

	"%[1]s/internal/core/module"
	"%[1]s/internal/modules/%[2]s/model"
	"%[1]s/internal/pkg/orm"
)

func %[3]sMigrations(db *orm.DB) []module.Migration {
	return []module.Migration{
		{
			ID: "%[2]s_%[4]s_001_create_table",
			Up: func(ctx context.Context) error {
				return db.AutoMigrate(&model.%[3]s{})
			},
			Down: func(ctx context.Context) error {
				return db.DB().Migrator().DropTable(&model.%[3]s{})
			},
		},
	}
}
`

// 中文：paymentChannelTemplate 声明当前包使用的常量。
// English: paymentChannelTemplate declares constants used by this package.
const paymentChannelTemplate = `package channel

import (
	"context"
	"encoding/json"
	"fmt"
)

type %[1]s struct {
	mchID      string
	apiKey     string
	gatewayURL string
	notifyURL  string
	client     gatewayHTTPClient
}

func New%[2]s(mchID, apiKey, gatewayURL, notifyURL string) *%[1]s {
	return &%[1]s{
		mchID:      mchID,
		apiKey:     apiKey,
		gatewayURL: gatewayURL,
		notifyURL:  notifyURL,
	}
}

func (c *%[1]s) Name() string { return "%[3]s" }

func (c *%[1]s) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	if notifyURL == "" {
		notifyURL = c.notifyURL
	}
	endpoint, err := gatewayEndpoint(c.gatewayURL, "/payment/create")
	if err != nil {
		return nil, err
	}
	payload := map[string]string{
		"mch_id":       c.mchID,
		"out_trade_no": outTradeNo,
		"subject":      subject,
		"amount":       fmt.Sprintf("%%d", amount),
		"scene":        scene,
		"notify_url":   notifyURL,
		"return_url":   returnURL,
		"open_id":      openID,
	}
	payload["sign"] = gatewaySign(c.apiKey, payload)
	var resp %[1]sResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return nil, err
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return nil, fmt.Errorf("%[3]s create payment failed: %%s", gatewayFirst(resp.values(), "message", "msg", "code"))
	}
	return &PayResult{
		PayURL:   gatewayFirst(resp.values(), "pay_url", "payUrl", "url"),
		QrCode:   gatewayFirst(resp.values(), "qr_code", "qrCode", "code_url"),
		PrepayID: gatewayFirst(resp.values(), "prepay_id", "prepayId"),
		TradeNo:  gatewayFirst(resp.values(), "transaction_id", "trade_no", "payment_id"),
		Params:   resp.values(),
	}, nil
}

func (c *%[1]s) VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error) {
	values, signature, err := gatewaySignedFields(rawData)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" && !gatewayVerifySignature(c.apiKey, values, signature) {
		return nil, fmt.Errorf("%[3]s callback signature invalid")
	}
	return &CallbackResult{
		OutTradeNo: gatewayFirst(values, "out_trade_no", "order_no", "merchant_order_no"),
		TradeNo:    gatewayFirst(values, "transaction_id", "trade_no", "payment_id"),
		Status:     gatewayPaymentStatus(gatewayFirst(values, "status", "trade_status", "result")),
		Amount:     gatewayParseAmount(gatewayFirst(values, "amount", "total_amount", "total_fee")),
		RawData:    rawData,
	}, nil
}

func (c *%[1]s) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	endpoint, err := gatewayEndpoint(c.gatewayURL, "/payment/refund")
	if err != nil {
		return nil, err
	}
	payload := map[string]string{
		"mch_id":        c.mchID,
		"out_trade_no":  outTradeNo,
		"out_refund_no": outRefundNo,
		"total_amount":  fmt.Sprintf("%%d", totalAmount),
		"refund_amount": fmt.Sprintf("%%d", refundAmount),
		"reason":        reason,
	}
	payload["sign"] = gatewaySign(c.apiKey, payload)
	var resp %[1]sResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return nil, err
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return nil, fmt.Errorf("%[3]s refund failed: %%s", gatewayFirst(resp.values(), "message", "msg", "code"))
	}
	return &RefundResult{
		RefundNo: gatewayFirst(resp.values(), "refund_no", "refund_id", "out_refund_no"),
		Status:   gatewayRefundStatus(gatewayFirst(resp.values(), "status", "refund_status", "result")),
	}, nil
}

func (c *%[1]s) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	endpoint, err := gatewayEndpoint(c.gatewayURL, "/payment/query")
	if err != nil {
		return nil, err
	}
	payload := map[string]string{
		"mch_id":       c.mchID,
		"out_trade_no": outTradeNo,
	}
	payload["sign"] = gatewaySign(c.apiKey, payload)
	var resp %[1]sResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return nil, err
	}
	values := resp.values()
	return &CallbackResult{
		OutTradeNo: gatewayFirst(values, "out_trade_no", "order_no", "merchant_order_no"),
		TradeNo:    gatewayFirst(values, "transaction_id", "trade_no", "payment_id"),
		Status:     gatewayPaymentStatus(gatewayFirst(values, "status", "trade_status", "result")),
		Amount:     gatewayParseAmount(gatewayFirst(values, "amount", "total_amount", "total_fee")),
	}, nil
}

func (c *%[1]s) ClosePayment(ctx context.Context, outTradeNo string) error {
	endpoint, err := gatewayEndpoint(c.gatewayURL, "/payment/close")
	if err != nil {
		return err
	}
	payload := map[string]string{
		"mch_id":       c.mchID,
		"out_trade_no": outTradeNo,
	}
	payload["sign"] = gatewaySign(c.apiKey, payload)
	var resp %[1]sResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return err
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return fmt.Errorf("%[3]s close payment failed: %%s", gatewayFirst(resp.values(), "message", "msg", "code"))
	}
	return nil
}

func (c *%[1]s) CallbackSuccess() any { return "success" }

func (c *%[1]s) CallbackFail() any { return "fail" }

type %[1]sResponse struct {
	Success bool            ` + "`json:\"success\"`" + `
	Code    string          ` + "`json:\"code\"`" + `
	Message string          ` + "`json:\"message\"`" + `
	Raw     json.RawMessage ` + "`json:\"-\"`" + `

	PayURL        string ` + "`json:\"pay_url\"`" + `
	QRCode        string ` + "`json:\"qr_code\"`" + `
	PrepayID      string ` + "`json:\"prepay_id\"`" + `
	TransactionID string ` + "`json:\"transaction_id\"`" + `
	OutTradeNo    string ` + "`json:\"out_trade_no\"`" + `
	Amount        string ` + "`json:\"amount\"`" + `
	Status        string ` + "`json:\"status\"`" + `
	RefundNo      string ` + "`json:\"refund_no\"`" + `
}

func (r *%[1]sResponse) UnmarshalJSON(data []byte) error {
	type alias %[1]sResponse
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	*out = %[1]sResponse(out)
	r.Raw = append(r.Raw[:0], data...)
	return nil
}

func (r %[1]sResponse) values() map[string]string {
	return map[string]string{
		"success":        fmt.Sprintf("%%t", r.Success),
		"code":           r.Code,
		"message":        r.Message,
		"pay_url":        r.PayURL,
		"qr_code":        r.QRCode,
		"prepay_id":      r.PrepayID,
		"transaction_id": r.TransactionID,
		"out_trade_no":   r.OutTradeNo,
		"amount":         r.Amount,
		"status":         r.Status,
		"refund_no":      r.RefundNo,
	}
}
`

// 中文：oauthProviderTemplate 声明当前包使用的常量。
// English: oauthProviderTemplate declares constants used by this package.
const oauthProviderTemplate = `package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type %[1]s struct {
	clientID     string
	clientSecret string
	authURL      string
	tokenURL     string
	userInfoURL  string
	scopes       []string
	client       *http.Client
}

func New%[2]s(clientID, clientSecret, authURL, tokenURL, userInfoURL string, scopes ...string) *%[1]s {
	return &%[1]s{
		clientID:     clientID,
		clientSecret: clientSecret,
		authURL:      authURL,
		tokenURL:     tokenURL,
		userInfoURL:  userInfoURL,
		scopes:       scopes,
		client:       http.DefaultClient,
	}
}

func (p *%[1]s) Name() string { return "%[3]s" }

func (p *%[1]s) AuthorizeURL(state, redirectURL string) string {
	u, err := url.Parse(p.authURL)
	if err != nil {
		return ""
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", p.clientID)
	q.Set("redirect_uri", redirectURL)
	q.Set("state", state)
	if len(p.scopes) > 0 {
		q.Set("scope", strings.Join(p.scopes, " "))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (p *%[1]s) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	return p.GetUserInfoWithRedirect(ctx, code, "")
}

func (p *%[1]s) GetUserInfoWithRedirect(ctx context.Context, code string, redirectURL string) (*UserInfo, error) {
	token, err := p.exchangeToken(ctx, code, redirectURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%[3]s userinfo failed: http %%d: %%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return &UserInfo{
		Provider:    "%[3]s",
		ProviderUID: firstOAuthString(raw, "id", "openid", "sub", "user_id"),
		OpenID:      firstOAuthString(raw, "openid", "id", "sub", "user_id"),
		UnionID:     firstOAuthString(raw, "unionid", "union_id"),
		Username:    firstOAuthString(raw, "username", "login", "screen_name"),
		Nickname:    firstOAuthString(raw, "nickname", "name", "username", "display_name"),
		Avatar:      firstOAuthString(raw, "avatar", "avatar_url", "picture", "headimgurl"),
		Email:       firstOAuthString(raw, "email"),
		Phone:       firstOAuthString(raw, "phone", "mobile"),
		RawData:     raw,
	}, nil
}

func (p *%[1]s) AuthURL(ctx context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

func (p *%[1]s) GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error) {
	info, err := p.GetUserInfoWithRedirect(ctx, code, redirectURL)
	if err != nil {
		return nil, err
	}
	return &OAuthUser{
		Provider:   info.Provider,
		ProviderID: firstOAuthString(map[string]any{"uid": info.ProviderUID, "openid": info.OpenID, "unionid": info.UnionID}, "uid", "openid", "unionid"),
		Username:   info.Username,
		Nickname:   info.Nickname,
		Avatar:     info.Avatar,
		Email:      info.Email,
		Phone:      info.Phone,
		RawData:    info.RawData,
	}, nil
}

func (p *%[1]s) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("%[3]s refresh token is required")
	}
	_, err := p.refreshAccessToken(ctx, refreshToken)
	return err
}

func (p *%[1]s) exchangeToken(ctx context.Context, code string, redirectURL string) (%[1]sToken, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("code", code)
	if redirectURL != "" {
		form.Set("redirect_uri", redirectURL)
	}
	return p.requestToken(ctx, form, "token exchange")
}

func (p *%[1]s) refreshAccessToken(ctx context.Context, refreshToken string) (%[1]sToken, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("refresh_token", refreshToken)
	return p.requestToken(ctx, form, "refresh token")
}

func (p *%[1]s) requestToken(ctx context.Context, form url.Values, operation string) (%[1]sToken, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return %[1]sToken{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return %[1]sToken{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return %[1]sToken{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return %[1]sToken{}, fmt.Errorf("%[3]s %%s failed: http %%d: %%s", operation, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	token, err := parse%[2]sToken(body)
	if err != nil {
		return %[1]sToken{}, err
	}
	if token.AccessToken == "" {
		return %[1]sToken{}, fmt.Errorf("%[3]s %%s response missing access_token", operation)
	}
	return token, nil
}

type %[1]sToken struct {
	AccessToken string ` + "`json:\"access_token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token\"`" + `
	TokenType   string ` + "`json:\"token_type\"`" + `
	ExpiresIn   int    ` + "`json:\"expires_in\"`" + `
}

func parse%[2]sToken(body []byte) (%[1]sToken, error) {
	var token %[1]sToken
	if err := json.Unmarshal(body, &token); err == nil && token.AccessToken != "" {
		return token, nil
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return %[1]sToken{}, err
	}
	token.AccessToken = values.Get("access_token")
	token.RefreshToken = values.Get("refresh_token")
	token.TokenType = values.Get("token_type")
	if expires := values.Get("expires_in"); expires != "" {
		fmt.Sscanf(expires, "%%d", &token.ExpiresIn)
	}
	if token.AccessToken == "" && values.Get("error") != "" {
		return %[1]sToken{}, fmt.Errorf("%[3]s token error: %%s %%s", values.Get("error"), values.Get("error_description"))
	}
	return token, nil
}

func firstOAuthString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key]; ok {
			if text := strings.TrimSpace(fmt.Sprint(value)); text != "" && text != "<nil>" {
				return text
			}
		}
	}
	return ""
}
`

// 中文：migrationTemplate 声明当前包使用的常量。
// English: migrationTemplate declares constants used by this package.
const migrationTemplate = `package migrations

import (
	"context"

	"github.com/spiringo/spiringo/internal/core/module"
)

func New%[2]sMigration() module.Migration {
	return module.Migration{
		ID: "%[1]s",
		Up: func(ctx context.Context) error {
			// Add schema or data changes here.
			return nil
		},
		Down: func(ctx context.Context) error {
			// Add rollback logic here when the change is reversible.
			return nil
		},
	}
}
`
