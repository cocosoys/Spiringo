package migrate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	"gorm.io/gorm"
)

// 中文：Record 定义当前包使用的数据结构或接口。
// English: Record defines a data structure or interface used by this package.
type Record struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `gorm:"primaryKey;size:191" json:"id"`
	// 中文：Module 保存当前结构中的配置或数据值。
	// English: Module stores a configuration or data value for this struct.
	Module string `gorm:"size:64;index;not null" json:"module"`
	// 中文：AppliedAt 保存当前结构中的配置或数据值。
	// English: AppliedAt stores a configuration or data value for this struct.
	AppliedAt time.Time `gorm:"autoCreateTime" json:"applied_at"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (Record) TableName() string { return "schema_migrations" }

// 中文：Store 定义当前包使用的数据结构或接口。
// English: Store defines a data structure or interface used by this package.
type Store struct {
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
	// 中文：ensured 保存当前结构中的配置或数据值。
	// English: ensured stores a configuration or data value for this struct.
	ensured bool
}

// 中文：NewStore 创建并返回对应组件实例。
// English: NewStore creates and returns the corresponding component instance.
func NewStore(db *orm.DB) *Store {
	if db == nil {
		return nil
	}
	return &Store{db: db}
}

// 中文：RunMigrations 执行当前包中的对应流程。
// English: RunMigrations executes the corresponding workflow in this package.
func (s *Store) RunMigrations(ctx context.Context, moduleName string, migrations []module.Migration) error {
	if s == nil || s.db == nil {
		return RunDirect(ctx, migrations)
	}
	if err := s.ensureTable(); err != nil {
		return err
	}
	for _, migration := range migrations {
		if migration.ID == "" {
			return fmt.Errorf("migration id is required for module %s", moduleName)
		}
		applied, err := s.Applied(ctx, migration.ID)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if migration.Up != nil {
			if err := migration.Up(ctx); err != nil {
				return fmt.Errorf("run migration %s: %w", migration.ID, err)
			}
		}
		if err := s.record(ctx, moduleName, migration.ID); err != nil {
			return err
		}
	}
	return nil
}

// 中文：Rollback 执行当前包中的对应流程。
// English: Rollback executes the corresponding workflow in this package.
func (s *Store) Rollback(ctx context.Context, moduleName string, migrations []module.Migration, steps int) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("migration store is not configured")
	}
	if steps <= 0 {
		return fmt.Errorf("rollback steps must be greater than zero")
	}
	if err := s.ensureTable(); err != nil {
		return err
	}

	byID := make(map[string]module.Migration, len(migrations))
	for _, migration := range migrations {
		if migration.ID != "" {
			byID[migration.ID] = migration
		}
	}

	records, err := s.appliedRecords(ctx, moduleName, steps)
	if err != nil {
		return err
	}
	for _, record := range records {
		migration, ok := byID[record.ID]
		if !ok {
			return fmt.Errorf("migration %s is applied but not registered", record.ID)
		}
		if migration.Down != nil {
			if err := migration.Down(ctx); err != nil {
				return fmt.Errorf("rollback migration %s: %w", record.ID, err)
			}
		}
		if err := s.deleteRecord(ctx, record.ID); err != nil {
			return err
		}
	}
	return nil
}

// 中文：Applied 执行当前包中的对应流程。
// English: Applied executes the corresponding workflow in this package.
func (s *Store) Applied(ctx context.Context, id string) (bool, error) {
	if err := s.ensureTable(); err != nil {
		return false, err
	}
	var record Record
	err := s.db.DB().WithContext(ctx).First(&record, "id = ?", id).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("check migration %s: %w", id, err)
}

// 中文：Records 执行当前包中的对应流程。
// English: Records executes the corresponding workflow in this package.
func (s *Store) Records(ctx context.Context, moduleName string) ([]Record, error) {
	if err := s.ensureTable(); err != nil {
		return nil, err
	}
	db := s.db.DB().WithContext(ctx).Order("applied_at asc")
	if moduleName != "" {
		db = db.Where("module = ?", moduleName)
	}
	var records []Record
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list migrations: %w", err)
	}
	return records, nil
}

// 中文：RunDirect 执行当前包中的对应流程。
// English: RunDirect executes the corresponding workflow in this package.
func RunDirect(ctx context.Context, migrations []module.Migration) error {
	for _, migration := range migrations {
		if migration.Up != nil {
			if err := migration.Up(ctx); err != nil {
				return fmt.Errorf("run migration %s: %w", migration.ID, err)
			}
		}
	}
	return nil
}

// 中文：ensureTable 执行当前包中的对应流程。
// English: ensureTable executes the corresponding workflow in this package.
func (s *Store) ensureTable() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("migration store is not configured")
	}
	if s.ensured {
		return nil
	}
	if err := s.db.AutoMigrate(&Record{}); err != nil {
		return fmt.Errorf("migrate schema_migrations: %w", err)
	}
	s.ensured = true
	return nil
}

// 中文：record 执行当前包中的对应流程。
// English: record executes the corresponding workflow in this package.
func (s *Store) record(ctx context.Context, moduleName, id string) error {
	record := Record{ID: id, Module: moduleName}
	if err := s.db.DB().WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("record migration %s: %w", id, err)
	}
	return nil
}

// 中文：appliedRecords 执行当前包中的对应流程。
// English: appliedRecords executes the corresponding workflow in this package.
func (s *Store) appliedRecords(ctx context.Context, moduleName string, limit int) ([]Record, error) {
	db := s.db.DB().WithContext(ctx).Order("applied_at desc").Limit(limit)
	if moduleName != "" {
		db = db.Where("module = ?", moduleName)
	}
	var records []Record
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list applied migrations: %w", err)
	}
	return records, nil
}

// 中文：deleteRecord 执行当前包中的对应流程。
// English: deleteRecord executes the corresponding workflow in this package.
func (s *Store) deleteRecord(ctx context.Context, id string) error {
	if err := s.db.DB().WithContext(ctx).Delete(&Record{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete migration record %s: %w", id, err)
	}
	return nil
}
