package repository

import (
	"context"

	"github.com/spiringo/spiringo/internal/modules/user/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：UserRepository 定义当前包使用的数据结构或接口。
// English: UserRepository defines a data structure or interface used by this package.
// UserRepository 用户数据访问
type UserRepository struct {
	// 中文：tdb 保存当前结构中的配置或数据值。
	// English: tdb stores a configuration or data value for this struct.
	tdb *orm.TenantDB
}

// 中文：NewUserRepository 创建并返回对应组件实例。
// English: NewUserRepository creates and returns the corresponding component instance.
// NewUserRepository 创建用户仓库
func NewUserRepository(tdb *orm.TenantDB) *UserRepository {
	return &UserRepository{tdb: tdb}
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.tdb.Create(ctx, user)
}

// 中文：GetByID 执行当前包中的对应流程。
// English: GetByID executes the corresponding workflow in this package.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.tdb.First(ctx, &user, "id = ?", id); err != nil {
		return nil, err
	}
	return &user, nil
}

// 中文：GetByUsername 执行当前包中的对应流程。
// English: GetByUsername executes the corresponding workflow in this package.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.tdb.First(ctx, &user, "username = ?", username); err != nil {
		return nil, err
	}
	return &user, nil
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	count, err := r.tdb.Paginate(ctx, &users, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.tdb.Update(ctx, user)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.tdb.Delete(ctx, &model.User{}, "id = ?", id)
}

// 中文：GetByEmail 执行当前包中的对应流程。
// English: GetByEmail executes the corresponding workflow in this package.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.tdb.First(ctx, &user, "email = ?", email); err != nil {
		return nil, err
	}
	return &user, nil
}
