package service

import (
	"context"
	"fmt"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/modules/tenant/model"
	"github.com/spiringo/spiringo/internal/modules/tenant/repository"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TenantService 定义当前包使用的数据结构或接口。
// English: TenantService defines a data structure or interface used by this package.
// TenantService 租户业务逻辑
type TenantService struct {
	// 中文：repo 保存当前结构中的配置或数据值。
	// English: repo stores a configuration or data value for this struct.
	repo *repository.TenantRepository
	// 中文：eventBus 保存当前结构中的配置或数据值。
	// English: eventBus stores a configuration or data value for this struct.
	eventBus *event.Bus
}

// 中文：NewTenantService 创建并返回对应组件实例。
// English: NewTenantService creates and returns the corresponding component instance.
// NewTenantService 创建租户服务
func NewTenantService(repo *repository.TenantRepository, eventBus *event.Bus) *TenantService {
	return &TenantService{repo: repo, eventBus: eventBus}
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (s *TenantService) Create(ctx context.Context, m *model.Tenant) error {
	// 校验Code唯一性
	existing, err := s.repo.GetByCode(ctx, m.Code)
	if err == nil && existing != nil {
		return types.ErrConflict.WithMessagef("tenant code %s already exists", m.Code)
	}

	if m.Strategy == "" {
		m.Strategy = "shared_db"
	}
	if m.Status == "" {
		m.Status = "active"
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return fmt.Errorf("create tenant: %w", err)
	}

	s.publishAsync(ctx, event.NewEvent(event.EventTenantCreated, &event.TenantEventPayload{
		TenantID:   m.ID,
		TenantName: m.Name,
		Strategy:   m.Strategy,
	}))
	return nil
}

// 中文：GetByID 执行当前包中的对应流程。
// English: GetByID executes the corresponding workflow in this package.
func (s *TenantService) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, types.ErrTenantNotFound
	}
	return tenant, nil
}

// 中文：GetByCode 执行当前包中的对应流程。
// English: GetByCode executes the corresponding workflow in this package.
func (s *TenantService) GetByCode(ctx context.Context, code string) (*model.Tenant, error) {
	tenant, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, types.ErrTenantNotFound
	}
	return tenant, nil
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (s *TenantService) List(ctx context.Context, page, pageSize int) ([]*model.Tenant, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (s *TenantService) Update(ctx context.Context, m *model.Tenant) error {
	if err := s.repo.Update(ctx, m); err != nil {
		return fmt.Errorf("update tenant: %w", err)
	}

	if m.Status == "suspended" {
		s.publishAsync(ctx, event.NewEvent(event.EventTenantSuspended, &event.TenantEventPayload{
			TenantID:   m.ID,
			TenantName: m.Name,
		}))
	} else if m.Status == "active" {
		s.publishAsync(ctx, event.NewEvent(event.EventTenantActivated, &event.TenantEventPayload{
			TenantID:   m.ID,
			TenantName: m.Name,
		}))
	}

	return nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (s *TenantService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// 中文：publishAsync 执行当前包中的对应流程。
// English: publishAsync executes the corresponding workflow in this package.
func (s *TenantService) publishAsync(ctx context.Context, e *event.Event) {
	if s.eventBus != nil {
		_ = s.eventBus.PublishAsync(ctx, e)
	}
}
