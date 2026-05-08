package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/modules/user/dto"
	"github.com/spiringo/spiringo/internal/modules/user/model"
	"github.com/spiringo/spiringo/internal/modules/user/repository"
	"github.com/spiringo/spiringo/pkg/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 中文：DefaultAdminConfig 定义当前包使用的数据结构或接口。
// English: DefaultAdminConfig defines a data structure or interface used by this package.
type DefaultAdminConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `yaml:"username" mapstructure:"username"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `yaml:"password" mapstructure:"password"`
	// 中文：EmailTemplate 保存当前结构中的配置或数据值。
	// English: EmailTemplate stores a configuration or data value for this struct.
	EmailTemplate string `yaml:"email_template" mapstructure:"email_template"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `yaml:"nickname" mapstructure:"nickname"`
}

// 中文：defaultAdminSeed 定义当前包使用的数据结构或接口。
// English: defaultAdminSeed defines a data structure or interface used by this package.
type defaultAdminSeed struct {
	// 中文：username 保存当前结构中的配置或数据值。
	// English: username stores a configuration or data value for this struct.
	username string
	// 中文：password 保存当前结构中的配置或数据值。
	// English: password stores a configuration or data value for this struct.
	password string
	// 中文：email 保存当前结构中的配置或数据值。
	// English: email stores a configuration or data value for this struct.
	email string
	// 中文：nickname 保存当前结构中的配置或数据值。
	// English: nickname stores a configuration or data value for this struct.
	nickname string
}

// 中文：UserService 定义当前包使用的数据结构或接口。
// English: UserService defines a data structure or interface used by this package.
type UserService struct {
	// 中文：repo 保存当前结构中的配置或数据值。
	// English: repo stores a configuration or data value for this struct.
	repo *repository.UserRepository
	// 中文：eventBus 保存当前结构中的配置或数据值。
	// English: eventBus stores a configuration or data value for this struct.
	eventBus *event.Bus
	// 中文：defaultAdmin 保存当前结构中的配置或数据值。
	// English: defaultAdmin stores a configuration or data value for this struct.
	defaultAdmin DefaultAdminConfig
}

// 中文：NewUserService 创建并返回对应组件实例。
// English: NewUserService creates and returns the corresponding component instance.
func NewUserService(repo *repository.UserRepository, eventBus *event.Bus) *UserService {
	return &UserService{repo: repo, eventBus: eventBus, defaultAdmin: defaultDefaultAdminConfig()}
}

// 中文：SetDefaultAdminConfig 执行当前包中的对应流程。
// English: SetDefaultAdminConfig executes the corresponding workflow in this package.
func (s *UserService) SetDefaultAdminConfig(cfg DefaultAdminConfig) {
	s.defaultAdmin = normalizeDefaultAdminConfig(cfg)
}

// 中文：defaultDefaultAdminConfig 执行当前包中的对应流程。
// English: defaultDefaultAdminConfig executes the corresponding workflow in this package.
func defaultDefaultAdminConfig() DefaultAdminConfig {
	return DefaultAdminConfig{
		Enabled:       true,
		Username:      "admin_%s",
		Password:      "changeme",
		EmailTemplate: "admin@%s.spiringo",
		Nickname:      "System Admin",
	}
}

// 中文：normalizeDefaultAdminConfig 执行当前包中的对应流程。
// English: normalizeDefaultAdminConfig executes the corresponding workflow in this package.
func normalizeDefaultAdminConfig(cfg DefaultAdminConfig) DefaultAdminConfig {
	defaults := defaultDefaultAdminConfig()
	defaults.Enabled = cfg.Enabled
	if strings.TrimSpace(cfg.Username) != "" {
		defaults.Username = strings.TrimSpace(cfg.Username)
	}
	if cfg.Password != "" {
		defaults.Password = cfg.Password
	}
	if strings.TrimSpace(cfg.EmailTemplate) != "" {
		defaults.EmailTemplate = strings.TrimSpace(cfg.EmailTemplate)
	}
	if strings.TrimSpace(cfg.Nickname) != "" {
		defaults.Nickname = strings.TrimSpace(cfg.Nickname)
	}
	return defaults
}

// 中文：seed 执行当前包中的对应流程。
// English: seed executes the corresponding workflow in this package.
func (cfg DefaultAdminConfig) seed(tenantID string) defaultAdminSeed {
	return defaultAdminSeed{
		username: formatTenantValue(cfg.Username, tenantID),
		password: cfg.Password,
		email:    formatTenantValue(cfg.EmailTemplate, tenantID),
		nickname: cfg.Nickname,
	}
}

// 中文：formatTenantValue 执行当前包中的对应流程。
// English: formatTenantValue executes the corresponding workflow in this package.
func formatTenantValue(pattern, tenantID string) string {
	if strings.Contains(pattern, "%s") {
		return fmt.Sprintf(pattern, tenantID)
	}
	return pattern
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (s *UserService) Create(ctx context.Context, req dto.CreateUserReq) (*model.User, error) {
	existing, err := s.repo.GetByUsername(ctx, req.Username)
	if err == nil && existing != nil {
		return nil, types.ErrUserExists.WithMessagef("username %s already exists", req.Username)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check username: %w", err)
	}

	if req.Email != "" {
		existingEmail, err := s.repo.GetByEmail(ctx, req.Email)
		if err == nil && existingEmail != nil {
			return nil, types.ErrUserExists.WithMessagef("email %s already exists", req.Email)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("check email: %w", err)
		}
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPwd),
		Nickname: req.Nickname,
		Status:   "active",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	s.publishAsync(ctx, event.NewEvent(event.EventUserCreated, &event.UserEventPayload{
		UserID:   user.ID,
		Username: user.Username,
		TenantID: user.TenantID,
	}))

	return user, nil
}

// 中文：GetByID 执行当前包中的对应流程。
// English: GetByID executes the corresponding workflow in this package.
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (s *UserService) Update(ctx context.Context, id string, req dto.UpdateUserReq) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, types.ErrUserNotFound
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	s.publishAsync(ctx, event.NewEvent(event.EventUserUpdated, &event.UserEventPayload{
		UserID:   user.ID,
		Username: user.Username,
		TenantID: user.TenantID,
	}))

	return user, nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return types.ErrUserNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	s.publishAsync(ctx, event.NewEvent(event.EventUserDeleted, &event.UserEventPayload{
		UserID: id,
	}))

	return nil
}

// 中文：VerifyPassword 执行当前包中的对应流程。
// English: VerifyPassword executes the corresponding workflow in this package.
func (s *UserService) VerifyPassword(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, types.ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, types.ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, types.ErrPasswordWrong
	}

	return user, nil
}

// 中文：VerifyPasswordForAuth 执行当前包中的对应流程。
// English: VerifyPasswordForAuth executes the corresponding workflow in this package.
func (s *UserService) VerifyPasswordForAuth(ctx context.Context, username, password string) (userID, tenantID string, err error) {
	user, err := s.VerifyPassword(ctx, username, password)
	if err != nil {
		return "", "", err
	}
	return user.ID, user.TenantID, nil
}

// 中文：CreateUserForAuth 执行当前包中的对应流程。
// English: CreateUserForAuth executes the corresponding workflow in this package.
func (s *UserService) CreateUserForAuth(ctx context.Context, username, email, phone, password, nickname string) (userID, tenantID string, err error) {
	user, err := s.Create(ctx, dto.CreateUserReq{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: password,
		Nickname: nickname,
	})
	if err != nil {
		return "", "", err
	}
	return user.ID, user.TenantID, nil
}

// 中文：CreateAdminForTenant 执行当前包中的对应流程。
// English: CreateAdminForTenant executes the corresponding workflow in this package.
func (s *UserService) CreateAdminForTenant(ctx context.Context, tenantID string) error {
	if !s.defaultAdmin.Enabled {
		return nil
	}

	tenantCtx := types.WithTenantID(ctx, tenantID)
	seed := s.defaultAdmin.seed(tenantID)
	existing, err := s.repo.GetByUsername(tenantCtx, seed.username)
	if err == nil && existing != nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check default admin: %w", err)
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(seed.password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	admin := &model.User{
		Username: seed.username,
		Email:    seed.email,
		Nickname: seed.nickname,
		Password: string(hashedPwd),
		Status:   "active",
	}
	admin.TenantID = tenantID

	if err := s.repo.Create(tenantCtx, admin); err != nil {
		return err
	}

	s.publishAsync(tenantCtx, event.NewEvent(event.EventUserCreated, &event.UserEventPayload{
		UserID:   admin.ID,
		Username: admin.Username,
		TenantID: tenantID,
	}))
	return nil
}

// 中文：publishAsync 执行当前包中的对应流程。
// English: publishAsync executes the corresponding workflow in this package.
func (s *UserService) publishAsync(ctx context.Context, e *event.Event) {
	if s.eventBus != nil {
		_ = s.eventBus.PublishAsync(ctx, e)
	}
}
