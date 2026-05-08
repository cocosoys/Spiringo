package orm

import (
	"context"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/spiringo/spiringo/pkg/types"
	"gorm.io/gorm"
)

// 中文：TenantDB 定义当前包使用的数据结构或接口。
// English: TenantDB defines a data structure or interface used by this package.
// TenantDB 多租户数据库包装
type TenantDB struct {
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *DB
}

// 中文：NewTenantDB 创建并返回对应组件实例。
// English: NewTenantDB creates and returns the corresponding component instance.
// NewTenantDB 创建多租户数据库
func NewTenantDB(db *DB) *TenantDB {
	return &TenantDB{db: db}
}

// 中文：DB 执行当前包中的对应流程。
// English: DB executes the corresponding workflow in this package.
// DB 获取原始DB
func (tdb *TenantDB) DB() *DB {
	return tdb.db
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
// Create 创建记录（自动注入tenant_id）
func (tdb *TenantDB) Create(ctx context.Context, model any) error {
	tenantID := types.GetTenantID(ctx)
	if tenantID != "" {
		setTenantID(model, tenantID)
	}
	return tdb.db.Create(ctx, model)
}

// 中文：Find 执行当前包中的对应流程。
// English: Find executes the corresponding workflow in this package.
// Find 查询记录（自动过滤tenant_id）
func (tdb *TenantDB) Find(ctx context.Context, models any, conds ...any) error {
	tenantID := types.GetTenantID(ctx)
	db := tdb.db
	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	return db.Find(ctx, models, conds...)
}

// 中文：First 执行当前包中的对应流程。
// English: First executes the corresponding workflow in this package.
// First 查询单条记录
func (tdb *TenantDB) First(ctx context.Context, model any, conds ...any) error {
	tenantID := types.GetTenantID(ctx)
	db := tdb.db
	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	return db.First(ctx, model, conds...)
}

// 中文：Count 执行当前包中的对应流程。
// English: Count executes the corresponding workflow in this package.
// Count 计数
func (tdb *TenantDB) Count(ctx context.Context, model any, conds ...any) (int64, error) {
	tenantID := types.GetTenantID(ctx)
	db := tdb.db
	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	return db.Count(ctx, model, conds...)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
// Delete 删除记录
func (tdb *TenantDB) Delete(ctx context.Context, model any, conds ...any) error {
	tenantID := types.GetTenantID(ctx)
	db := tdb.db
	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	return db.Delete(ctx, model, conds...)
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
// Update 更新记录
func (tdb *TenantDB) Update(ctx context.Context, model any) error {
	tenantID := types.GetTenantID(ctx)
	if tenantID != "" {
		setTenantID(model, tenantID)
		tx := tdb.db.db.WithContext(ctx).
			Model(model).
			Where("tenant_id = ?", tenantID).
			Select("*").
			Updates(model)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}
	return tdb.db.Update(ctx, model)
}

// 中文：Transaction 执行当前包中的对应流程。
// English: Transaction executes the corresponding workflow in this package.
// Transaction 事务
func (tdb *TenantDB) Transaction(ctx context.Context, fn func(tx *TenantDB) error) error {
	return tdb.db.Transaction(ctx, func(tx *DB) error {
		return fn(&TenantDB{db: tx})
	})
}

// 中文：Paginate 执行当前包中的对应流程。
// English: Paginate executes the corresponding workflow in this package.
// Paginate 分页查询
func (tdb *TenantDB) Paginate(ctx context.Context, models any, page, pageSize int, conds ...any) (int64, error) {
	tenantID := types.GetTenantID(ctx)
	db := tdb.db
	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	return db.Paginate(ctx, models, page, pageSize, conds...)
}

// 中文：setTenantID 执行当前包中的对应流程。
// English: setTenantID executes the corresponding workflow in this package.
// setTenantID 通过反射设置模型的tenant_id字段
func setTenantID(model any, tenantID string) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.FieldByName("TenantID")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(tenantID)
	}
}

// 中文：BaseModel 定义当前包使用的数据结构或接口。
// English: BaseModel defines a data structure or interface used by this package.
// BaseModel 基础模型
type BaseModel struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `gorm:"primaryKey;size:36" json:"id"`
	// 中文：CreatedAt 保存当前结构中的配置或数据值。
	// English: CreatedAt stores a configuration or data value for this struct.
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	// 中文：UpdatedAt 保存当前结构中的配置或数据值。
	// English: UpdatedAt stores a configuration or data value for this struct.
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// 中文：TenantModel 定义当前包使用的数据结构或接口。
// English: TenantModel defines a data structure or interface used by this package.
// TenantModel 多租户基础模型
type TenantModel struct {
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `gorm:"index;not null;size:36" json:"tenant_id"`
}

// 中文：TenantBaseModel 定义当前包使用的数据结构或接口。
// English: TenantBaseModel defines a data structure or interface used by this package.
// TenantBaseModel 多租户基础模型（组合）
type TenantBaseModel struct {
	// 中文：BaseModel 嵌入复用该类型提供的能力。
	// English: BaseModel embeds reusable behavior from that type.
	BaseModel
	// 中文：TenantModel 嵌入复用该类型提供的能力。
	// English: TenantModel embeds reusable behavior from that type.
	TenantModel
}

// 中文：BeforeCreate 执行当前包中的对应流程。
// English: BeforeCreate executes the corresponding workflow in this package.
// BeforeCreate Gorm钩子（自动生成ID）
func (m *BaseModel) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
