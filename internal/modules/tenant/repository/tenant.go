package repository

import (
	"context"

	"github.com/spiringo/spiringo/internal/modules/tenant/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：TenantRepository 定义当前包使用的数据结构或接口。
// English: TenantRepository defines a data structure or interface used by this package.
// TenantRepository 租户数据访问
type TenantRepository struct {
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
}

// 中文：NewTenantRepository 创建并返回对应组件实例。
// English: NewTenantRepository creates and returns the corresponding component instance.
// NewTenantRepository 创建租户仓库
func NewTenantRepository(db *orm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (r *TenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	return r.db.Create(ctx, tenant)
}

// 中文：GetByID 执行当前包中的对应流程。
// English: GetByID executes the corresponding workflow in this package.
func (r *TenantRepository) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.First(ctx, &tenant, "id = ?", id); err != nil {
		return nil, err
	}
	return &tenant, nil
}

// 中文：GetByCode 执行当前包中的对应流程。
// English: GetByCode executes the corresponding workflow in this package.
func (r *TenantRepository) GetByCode(ctx context.Context, code string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.First(ctx, &tenant, "code = ?", code); err != nil {
		return nil, err
	}
	return &tenant, nil
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (r *TenantRepository) List(ctx context.Context, page, pageSize int) ([]*model.Tenant, int64, error) {
	var tenants []*model.Tenant
	count, err := r.db.Paginate(ctx, &tenants, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return tenants, count, nil
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (r *TenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	return r.db.Update(ctx, tenant)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	return r.db.Delete(ctx, &model.Tenant{}, "id = ?", id)
}
