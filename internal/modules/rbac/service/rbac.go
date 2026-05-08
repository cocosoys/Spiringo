package service

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/modules/rbac/dto"
	"github.com/spiringo/spiringo/internal/modules/rbac/model"
	"github.com/spiringo/spiringo/internal/modules/rbac/repository"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：RBACService 定义当前包使用的数据结构或接口。
// English: RBACService defines a data structure or interface used by this package.
type RBACService struct {
	// 中文：repo 保存当前结构中的配置或数据值。
	// English: repo stores a configuration or data value for this struct.
	repo *repository.RBACRepository
	// 中文：abac 保存当前结构中的配置或数据值。
	// English: abac stores a configuration or data value for this struct.
	abac *ABACEngine
}

// 中文：NewRBACService 创建并返回对应组件实例。
// English: NewRBACService creates and returns the corresponding component instance.
func NewRBACService(repo *repository.RBACRepository) *RBACService {
	return &RBACService{repo: repo, abac: NewABACEngine(nil)}
}

// 中文：CreateRole 执行当前包中的对应流程。
// English: CreateRole executes the corresponding workflow in this package.
func (s *RBACService) CreateRole(ctx context.Context, req dto.CreateRoleReq) (*model.Role, error) {
	existing, _ := s.repo.GetRoleByCode(ctx, req.Code)
	if existing != nil {
		return nil, types.ErrRoleExists.WithMessagef("role code %s already exists", req.Code)
	}

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      "active",
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}
	return role, nil
}

// 中文：GetRole 执行当前包中的对应流程。
// English: GetRole executes the corresponding workflow in this package.
func (s *RBACService) GetRole(ctx context.Context, id string) (*model.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, types.ErrRoleNotFound
	}
	return role, nil
}

// 中文：ListRoles 执行当前包中的对应流程。
// English: ListRoles executes the corresponding workflow in this package.
func (s *RBACService) ListRoles(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error) {
	return s.repo.ListRoles(ctx, page, pageSize)
}

// 中文：UpdateRole 执行当前包中的对应流程。
// English: UpdateRole executes the corresponding workflow in this package.
func (s *RBACService) UpdateRole(ctx context.Context, id string, req dto.UpdateRoleReq) (*model.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, types.ErrRoleNotFound
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != "" {
		role.Status = req.Status
	}

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}
	return role, nil
}

// 中文：DeleteRole 执行当前包中的对应流程。
// English: DeleteRole executes the corresponding workflow in this package.
func (s *RBACService) DeleteRole(ctx context.Context, id string) error {
	if _, err := s.repo.GetRoleByID(ctx, id); err != nil {
		return types.ErrRoleNotFound
	}
	return s.repo.DeleteRole(ctx, id)
}

// 中文：ListPermissions 执行当前包中的对应流程。
// English: ListPermissions executes the corresponding workflow in this package.
func (s *RBACService) ListPermissions(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error) {
	return s.repo.ListPermissions(ctx, page, pageSize)
}

// 中文：AssignPermissions 执行当前包中的对应流程。
// English: AssignPermissions executes the corresponding workflow in this package.
func (s *RBACService) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	if _, err := s.repo.GetRoleByID(ctx, roleID); err != nil {
		return types.ErrRoleNotFound
	}

	for _, pid := range permissionIDs {
		if _, err := s.repo.GetPermissionByID(ctx, pid); err != nil {
			return types.ErrPermNotFound.WithMessagef("permission %s not found", pid)
		}
	}

	return s.repo.AssignPermissions(ctx, roleID, permissionIDs)
}

// 中文：GetRolePermissions 执行当前包中的对应流程。
// English: GetRolePermissions executes the corresponding workflow in this package.
func (s *RBACService) GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error) {
	return s.repo.GetPermissionsByRoleID(ctx, roleID)
}

// 中文：CheckPermission 执行当前包中的对应流程。
// English: CheckPermission executes the corresponding workflow in this package.
func (s *RBACService) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	return s.repo.CheckPermission(ctx, userID, resource, action)
}

// 中文：CreateDefaultRoles 执行当前包中的对应流程。
// English: CreateDefaultRoles executes the corresponding workflow in this package.
func (s *RBACService) CreateDefaultRoles(ctx context.Context, tenantID string) error {
	ctx = types.WithTenantID(ctx, tenantID)
	permissionIDs, err := s.ensureDefaultPermissions(ctx)
	if err != nil {
		return err
	}

	defaultRoles := []struct {
		Name        string
		Code        string
		Description string
		Permissions []string
	}{
		{"Administrator", "admin", "Tenant administrator with all default permissions", permissionIDs},
		{"Viewer", "viewer", "Read-only tenant user", nil},
	}

	for _, r := range defaultRoles {
		existing, _ := s.repo.GetRoleByCode(ctx, r.Code)
		role := existing
		if role == nil {
			role = &model.Role{
				Name:        r.Name,
				Code:        r.Code,
				Description: r.Description,
				Status:      "active",
			}
			role.TenantID = tenantID
			if err := s.repo.CreateRole(ctx, role); err != nil {
				return fmt.Errorf("create default role %s: %w", r.Code, err)
			}
		}
		if len(r.Permissions) > 0 {
			if err := s.repo.AssignPermissions(ctx, role.ID, r.Permissions); err != nil {
				return fmt.Errorf("assign default permissions to role %s: %w", r.Code, err)
			}
		}
	}
	return nil
}

// 中文：AssignDefaultRole 执行当前包中的对应流程。
// English: AssignDefaultRole executes the corresponding workflow in this package.
func (s *RBACService) AssignDefaultRole(ctx context.Context, tenantID, userID string) error {
	return s.AssignDefaultRoleByCode(ctx, tenantID, userID, "viewer")
}

// 中文：AssignDefaultRoleByCode 执行当前包中的对应流程。
// English: AssignDefaultRoleByCode executes the corresponding workflow in this package.
func (s *RBACService) AssignDefaultRoleByCode(ctx context.Context, tenantID, userID, roleCode string) error {
	ctx = types.WithTenantID(ctx, tenantID)
	role, err := s.repo.GetRoleByCode(ctx, roleCode)
	if err != nil {
		if createErr := s.CreateDefaultRoles(ctx, tenantID); createErr != nil {
			return createErr
		}
		role, err = s.repo.GetRoleByCode(ctx, roleCode)
	}
	if err != nil {
		return fmt.Errorf("default %s role not found: %w", roleCode, err)
	}
	return s.repo.AssignRoleToUser(ctx, userID, role.ID)
}

// 中文：ensureDefaultPermissions 执行当前包中的对应流程。
// English: ensureDefaultPermissions executes the corresponding workflow in this package.
func (s *RBACService) ensureDefaultPermissions(ctx context.Context) ([]string, error) {
	defaultPermissions := []struct {
		Name     string
		Code     string
		Resource string
		Action   string
	}{
		{"All Permissions", "system.all", "*", "*"},
		{"Read RBAC Roles", "rbac.roles.read", "rbac.roles", "read"},
		{"Create RBAC Roles", "rbac.roles.create", "rbac.roles", "create"},
		{"Update RBAC Roles", "rbac.roles.update", "rbac.roles", "update"},
		{"Delete RBAC Roles", "rbac.roles.delete", "rbac.roles", "delete"},
		{"Read RBAC Permissions", "rbac.permissions.read", "rbac.permissions", "read"},
		{"Assign RBAC Permissions", "rbac.permissions.assign", "rbac.permissions", "assign"},
	}

	ids := make([]string, 0, len(defaultPermissions))
	for _, p := range defaultPermissions {
		perm, _ := s.repo.GetPermissionByCode(ctx, p.Code)
		if perm == nil {
			perm = &model.Permission{
				Name:     p.Name,
				Code:     p.Code,
				Resource: p.Resource,
				Action:   p.Action,
			}
			if err := s.repo.CreatePermission(ctx, perm); err != nil {
				return nil, fmt.Errorf("create default permission %s: %w", p.Code, err)
			}
		}
		ids = append(ids, perm.ID)
	}
	return ids, nil
}
