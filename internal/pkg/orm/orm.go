package orm

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 中文：Dialector 定义当前包使用的数据结构或接口。
// English: Dialector defines a data structure or interface used by this package.
type Dialector = gorm.Dialector

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string
	// 中文：DSN 保存当前结构中的配置或数据值。
	// English: DSN stores a configuration or data value for this struct.
	DSN string
	// 中文：MaxIdle 保存当前结构中的配置或数据值。
	// English: MaxIdle stores a configuration or data value for this struct.
	MaxIdle int
	// 中文：MaxOpen 保存当前结构中的配置或数据值。
	// English: MaxOpen stores a configuration or data value for this struct.
	MaxOpen int
	// 中文：ConnMaxLifetime 保存当前结构中的配置或数据值。
	// English: ConnMaxLifetime stores a configuration or data value for this struct.
	ConnMaxLifetime time.Duration
	// 中文：LogLevel 保存当前结构中的配置或数据值。
	// English: LogLevel stores a configuration or data value for this struct.
	LogLevel string
	// 中文：ReadReplicas 保存当前结构中的配置或数据值。
	// English: ReadReplicas stores a configuration or data value for this struct.
	ReadReplicas []EndpointConfig
}

// 中文：EndpointConfig 定义当前包使用的数据结构或接口。
// English: EndpointConfig defines a data structure or interface used by this package.
type EndpointConfig struct {
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string
	// 中文：DSN 保存当前结构中的配置或数据值。
	// English: DSN stores a configuration or data value for this struct.
	DSN string
}

// 中文：DB 定义当前包使用的数据结构或接口。
// English: DB defines a data structure or interface used by this package.
// DB wraps GORM with primary-write and replica-read routing.
type DB struct {
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *gorm.DB
	// 中文：read 保存当前结构中的配置或数据值。
	// English: read stores a configuration or data value for this struct.
	read *gorm.DB
	// 中文：replicas 保存当前结构中的配置或数据值。
	// English: replicas stores a configuration or data value for this struct.
	replicas []*gorm.DB
	// 中文：nextRead 保存当前结构中的配置或数据值。
	// English: nextRead stores a configuration or data value for this struct.
	nextRead atomic.Uint64
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
func New(cfg Config) (*DB, error) {
	primary, err := openGormDB(cfg.Driver, cfg.DSN, cfg)
	if err != nil {
		return nil, err
	}

	replicas := make([]*gorm.DB, 0, len(cfg.ReadReplicas))
	for i, replica := range cfg.ReadReplicas {
		driver := replica.Driver
		if driver == "" {
			driver = cfg.Driver
		}
		db, err := openGormDB(driver, replica.DSN, cfg)
		if err != nil {
			_ = closeGormDB(primary)
			for _, opened := range replicas {
				_ = closeGormDB(opened)
			}
			return nil, fmt.Errorf("open read replica %d: %w", i, err)
		}
		replicas = append(replicas, db)
	}

	return &DB{db: primary, replicas: replicas, config: cfg}, nil
}

// 中文：openGormDB 执行当前包中的对应流程。
// English: openGormDB executes the corresponding workflow in this package.
func openGormDB(driver, dsn string, cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch driver {
	case "mysql":
		dialector = newMySQLDialector(dsn)
	case "postgres":
		dialector = newPostgresDialector(dsn)
	case "sqlite":
		dialector = newSQLiteDialector(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	logLevel := logger.Info
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	return db, nil
}

// 中文：DB 执行当前包中的对应流程。
// English: DB executes the corresponding workflow in this package.
// DB returns the primary GORM instance for advanced write-side operations.
func (d *DB) DB() *gorm.DB {
	return d.db
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (d *DB) Close() error {
	errs := []error{closeGormDB(d.db)}
	for _, replica := range d.replicas {
		errs = append(errs, closeGormDB(replica))
	}
	return errors.Join(errs...)
}

// 中文：Ping 执行当前包中的对应流程。
// English: Ping executes the corresponding workflow in this package.
func (d *DB) Ping() error {
	errs := []error{pingGormDB(d.db)}
	for i, replica := range d.replicas {
		if err := pingGormDB(replica); err != nil {
			errs = append(errs, fmt.Errorf("read replica %d: %w", i, err))
		}
	}
	return errors.Join(errs...)
}

// 中文：closeGormDB 执行当前包中的对应流程。
// English: closeGormDB executes the corresponding workflow in this package.
func closeGormDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// 中文：pingGormDB 执行当前包中的对应流程。
// English: pingGormDB executes the corresponding workflow in this package.
func pingGormDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// 中文：reader 执行当前包中的对应流程。
// English: reader executes the corresponding workflow in this package.
func (d *DB) reader() *gorm.DB {
	if d.read != nil {
		return d.read
	}
	if len(d.replicas) == 0 {
		return d.db
	}
	idx := d.nextRead.Add(1) - 1
	return d.replicas[int(idx%uint64(len(d.replicas)))]
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (d *DB) Create(_ context.Context, model any) error {
	return d.db.Create(model).Error
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (d *DB) Update(_ context.Context, model any) error {
	return d.db.Save(model).Error
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (d *DB) Delete(_ context.Context, model any, conds ...any) error {
	return d.db.Delete(model, conds...).Error
}

// 中文：First 执行当前包中的对应流程。
// English: First executes the corresponding workflow in this package.
func (d *DB) First(_ context.Context, model any, conds ...any) error {
	return d.reader().First(model, conds...).Error
}

// 中文：Find 执行当前包中的对应流程。
// English: Find executes the corresponding workflow in this package.
func (d *DB) Find(_ context.Context, models any, conds ...any) error {
	return d.reader().Find(models, conds...).Error
}

// 中文：Count 执行当前包中的对应流程。
// English: Count executes the corresponding workflow in this package.
func (d *DB) Count(_ context.Context, model any, conds ...any) (int64, error) {
	var count int64
	db := d.reader().Model(model)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}
	err := db.Count(&count).Error
	return count, err
}

// 中文：Where 执行当前包中的对应流程。
// English: Where executes the corresponding workflow in this package.
func (d *DB) Where(query any, args ...any) *DB {
	return d.clone(d.db.Where(query, args...), d.reader().Where(query, args...))
}

// 中文：Order 执行当前包中的对应流程。
// English: Order executes the corresponding workflow in this package.
func (d *DB) Order(value string) *DB {
	return d.clone(d.db.Order(value), d.reader().Order(value))
}

// 中文：Limit 执行当前包中的对应流程。
// English: Limit executes the corresponding workflow in this package.
func (d *DB) Limit(limit int) *DB {
	return d.clone(d.db.Limit(limit), d.reader().Limit(limit))
}

// 中文：Offset 执行当前包中的对应流程。
// English: Offset executes the corresponding workflow in this package.
func (d *DB) Offset(offset int) *DB {
	return d.clone(d.db.Offset(offset), d.reader().Offset(offset))
}

// 中文：Transaction 执行当前包中的对应流程。
// English: Transaction executes the corresponding workflow in this package.
func (d *DB) Transaction(_ context.Context, fn func(tx *DB) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return fn(&DB{db: tx, read: tx, config: d.config})
	})
}

// 中文：Raw 执行当前包中的对应流程。
// English: Raw executes the corresponding workflow in this package.
func (d *DB) Raw(sql string, values ...any) *DB {
	return d.clone(d.db.Raw(sql, values...), d.reader().Raw(sql, values...))
}

// 中文：Master 执行当前包中的对应流程。
// English: Master executes the corresponding workflow in this package.
func (d *DB) Master() *DB {
	return d.clone(d.db, d.db)
}

// 中文：Exec 执行当前包中的对应流程。
// English: Exec executes the corresponding workflow in this package.
func (d *DB) Exec(sql string, values ...any) error {
	return d.db.Exec(sql, values...).Error
}

// 中文：Paginate 执行当前包中的对应流程。
// English: Paginate executes the corresponding workflow in this package.
func (d *DB) Paginate(_ context.Context, models any, page, pageSize int, conds ...any) (int64, error) {
	page, pageSize = normalizePagination(page, pageSize)
	var count int64
	db := d.reader()
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	if err := db.Model(models).Count(&count).Error; err != nil {
		return 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(models).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// 中文：normalizePagination 执行当前包中的对应流程。
// English: normalizePagination executes the corresponding workflow in this package.
func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// 中文：AutoMigrate 执行当前包中的对应流程。
// English: AutoMigrate executes the corresponding workflow in this package.
func (d *DB) AutoMigrate(models ...any) error {
	return d.db.AutoMigrate(models...)
}

// 中文：clone 执行当前包中的对应流程。
// English: clone executes the corresponding workflow in this package.
func (d *DB) clone(write, read *gorm.DB) *DB {
	child := &DB{
		db:       write,
		read:     read,
		replicas: d.replicas,
		config:   d.config,
	}
	child.nextRead.Store(d.nextRead.Load())
	return child
}
