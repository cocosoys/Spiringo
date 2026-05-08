package repository

import (
	"context"
	"errors"

	"github.com/spiringo/spiringo/internal/modules/rbac/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	"gorm.io/gorm"
)

// 中文：RBACRepository 定义当前包使用的数据结构或接口。
// English: RBACRepository defines a data structure or interface used by this package.
// RBACRepository 权限管理数据访问
type RBACRepository struct {
	// 中文：tdb 保存当前结构中的配置或数据值。
	// English: tdb stores a configuration or data value for this struct.
	tdb *orm.TenantDB
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
	// 用于非租户隔离的关联表
}

// 中文：NewRBACRepository 创建并返回对应组件实例。
// English: NewRBACRepository creates and returns the corresponding component instance.
// NewRBACRepository 创建权限管理仓库
func NewRBACRepository(tdb *orm.TenantDB, db *orm.DB) *RBACRepository {
	return &RBACRepository{tdb: tdb, db: db}
}

// ---- 角色 ----

// 中文：CreateRole 执行当前包中的对应流程。
// English: CreateRole executes the corresponding workflow in this package.
func (r *RBACRepository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.tdb.Create(ctx, role)
}

// 中文：GetRoleByID 执行当前包中的对应流程。
// English: GetRoleByID executes the corresponding workflow in this package.
func (r *RBACRepository) GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	if err := r.tdb.First(ctx, &role, "id = ?", id); err != nil {
		return nil, err
	}
	return &role, nil
}

// 中文：GetRoleByCode 执行当前包中的对应流程。
// English: GetRoleByCode executes the corresponding workflow in this package.
func (r *RBACRepository) GetRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	if err := r.tdb.First(ctx, &role, "code = ?", code); err != nil {
		return nil, err
	}
	return &role, nil
}

// 中文：ListRoles 执行当前包中的对应流程。
// English: ListRoles executes the corresponding workflow in this package.
func (r *RBACRepository) ListRoles(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	count, err := r.tdb.Paginate(ctx, &roles, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

// 中文：UpdateRole 执行当前包中的对应流程。
// English: UpdateRole executes the corresponding workflow in this package.
func (r *RBACRepository) UpdateRole(ctx context.Context, role *model.Role) error {
	return r.tdb.Update(ctx, role)
}

// 中文：DeleteRole 执行当前包中的对应流程。
// English: DeleteRole executes the corresponding workflow in this package.
func (r *RBACRepository) DeleteRole(ctx context.Context, id string) error {
	return r.tdb.Delete(ctx, &model.Role{}, "id = ?", id)
}

// ---- 权限 ----

// 中文：CreatePermission 执行当前包中的对应流程。
// English: CreatePermission executes the corresponding workflow in this package.
func (r *RBACRepository) CreatePermission(ctx context.Context, perm *model.Permission) error {
	return r.tdb.Create(ctx, perm)
}

// 中文：ListPermissions 执行当前包中的对应流程。
// English: ListPermissions executes the corresponding workflow in this package.
func (r *RBACRepository) ListPermissions(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error) {
	var perms []*model.Permission
	count, err := r.tdb.Paginate(ctx, &perms, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return perms, count, nil
}

// 中文：GetPermissionByID 执行当前包中的对应流程。
// English: GetPermissionByID executes the corresponding workflow in this package.
func (r *RBACRepository) GetPermissionByID(ctx context.Context, id string) (*model.Permission, error) {
	var perm model.Permission
	if err := r.tdb.First(ctx, &perm, "id = ?", id); err != nil {
		return nil, err
	}
	return &perm, nil
}

// 中文：GetPermissionByCode 执行当前包中的对应流程。
// English: GetPermissionByCode executes the corresponding workflow in this package.
func (r *RBACRepository) GetPermissionByCode(ctx context.Context, code string) (*model.Permission, error) {
	var perm model.Permission
	if err := r.tdb.First(ctx, &perm, "code = ?", code); err != nil {
		return nil, err
	}
	return &perm, nil
}

// ---- 角色-权限关联 ----

// 中文：AssignPermissions 执行当前包中的对应流程。
// English: AssignPermissions executes the corresponding workflow in this package.
func (r *RBACRepository) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	// 先删除旧关联
	if err := r.db.Delete(ctx, &model.RolePermission{}, "role_id = ?", roleID); err != nil {
		return err
	}
	// 批量创建新关联
	for _, pid := range permissionIDs {
		rp := &model.RolePermission{RoleID: roleID, PermissionID: pid}
		if err := r.db.Create(ctx, rp); err != nil {
			return err
		}
	}
	return nil
}

// 中文：GetPermissionsByRoleID 执行当前包中的对应流程。
// English: GetPermissionsByRoleID executes the corresponding workflow in this package.
func (r *RBACRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*model.Permission, error) {
	var rps []*model.RolePermission
	if err := r.db.Find(ctx, &rps, "role_id = ?", roleID); err != nil {
		return nil, err
	}

	var perms []*model.Permission
	for _, rp := range rps {
		var perm model.Permission
		if err := r.tdb.First(ctx, &perm, "id = ?", rp.PermissionID); err == nil {
			perms = append(perms, &perm)
		}
	}
	return perms, nil
}

// ---- 用户-角色关联 ----

// 中文：AssignRoleToUser 执行当前包中的对应流程。
// English: AssignRoleToUser executes the corresponding workflow in this package.
func (r *RBACRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	var existing model.UserRole
	err := r.db.First(ctx, &existing, "user_id = ? AND role_id = ?", userID, roleID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	ur := &model.UserRole{UserID: userID, RoleID: roleID}
	return r.db.Create(ctx, ur)
}

// 中文：GetRolesByUserID 执行当前包中的对应流程。
// English: GetRolesByUserID executes the corresponding workflow in this package.
func (r *RBACRepository) GetRolesByUserID(ctx context.Context, userID string) ([]*model.Role, error) {
	var urs []*model.UserRole
	if err := r.db.Find(ctx, &urs, "user_id = ?", userID); err != nil {
		return nil, err
	}

	var roles []*model.Role
	for _, ur := range urs {
		role, err := r.GetRoleByID(ctx, ur.RoleID)
		if err == nil {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

// 中文：CheckPermission 执行当前包中的对应流程。
// English: CheckPermission executes the corresponding workflow in this package.
// CheckPermission 检查用户是否有指定权限
func (r *RBACRepository) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	roles, err := r.GetRolesByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.Status != "active" {
			continue
		}
		perms, err := r.GetPermissionsByRoleID(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, p := range perms {
			if p.Resource == resource && p.Action == action {
				return true, nil
			}
			if p.Resource == "*" && p.Action == "*" {
				return true, nil
			}
		}
	}
	return false, nil
}
