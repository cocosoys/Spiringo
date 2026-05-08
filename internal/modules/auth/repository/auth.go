package repository

import (
	"context"

	"github.com/spiringo/spiringo/internal/modules/auth/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：AuthRepository 定义当前包使用的数据结构或接口。
// English: AuthRepository defines a data structure or interface used by this package.
// AuthRepository 认证数据访问
type AuthRepository struct {
	// 中文：tdb 保存当前结构中的配置或数据值。
	// English: tdb stores a configuration or data value for this struct.
	tdb *orm.TenantDB
}

// 中文：NewAuthRepository 创建并返回对应组件实例。
// English: NewAuthRepository creates and returns the corresponding component instance.
// NewAuthRepository 创建认证仓库
func NewAuthRepository(tdb *orm.TenantDB) *AuthRepository {
	return &AuthRepository{tdb: tdb}
}

// 中文：CreateOAuthBinding 执行当前包中的对应流程。
// English: CreateOAuthBinding executes the corresponding workflow in this package.
// CreateOAuthBinding 创建OAuth绑定
func (r *AuthRepository) CreateOAuthBinding(ctx context.Context, binding *model.OAuthBinding) error {
	return r.tdb.Create(ctx, binding)
}

// 中文：GetOAuthBinding 执行当前包中的对应流程。
// English: GetOAuthBinding executes the corresponding workflow in this package.
// GetOAuthBinding 按提供商和UID查询绑定
func (r *AuthRepository) GetOAuthBinding(ctx context.Context, provider, providerUID string) (*model.OAuthBinding, error) {
	var binding model.OAuthBinding
	if err := r.tdb.First(ctx, &binding, "provider = ? AND provider_uid = ?", provider, providerUID); err != nil {
		return nil, err
	}
	return &binding, nil
}

// 中文：GetOAuthBindingsByUserID 执行当前包中的对应流程。
// English: GetOAuthBindingsByUserID executes the corresponding workflow in this package.
// GetOAuthBindingsByUserID 按用户ID查询所有绑定
func (r *AuthRepository) GetOAuthBindingsByUserID(ctx context.Context, userID string) ([]*model.OAuthBinding, error) {
	var bindings []*model.OAuthBinding
	if err := r.tdb.Find(ctx, &bindings, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return bindings, nil
}

// 中文：DeleteOAuthBinding 执行当前包中的对应流程。
// English: DeleteOAuthBinding executes the corresponding workflow in this package.
// DeleteOAuthBinding 删除OAuth绑定
func (r *AuthRepository) DeleteOAuthBinding(ctx context.Context, id string) error {
	return r.tdb.Delete(ctx, &model.OAuthBinding{}, "id = ?", id)
}

// 中文：DeleteByUserID 执行当前包中的对应流程。
// English: DeleteByUserID executes the corresponding workflow in this package.
// DeleteByUserID 删除用户的所有OAuth绑定
func (r *AuthRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.tdb.Delete(ctx, &model.OAuthBinding{}, "user_id = ?", userID)
}
