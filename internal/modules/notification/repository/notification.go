package repository

import (
	"context"
	"time"

	"github.com/spiringo/spiringo/internal/modules/notification/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	"github.com/spiringo/spiringo/pkg/types"
	"gorm.io/gorm"
)

// 中文：NotificationFilter 定义当前包使用的数据结构或接口。
// English: NotificationFilter defines a data structure or interface used by this package.
// NotificationFilter controls inbox queries.
type NotificationFilter struct {
	// 中文：Page 保存当前结构中的配置或数据值。
	// English: Page stores a configuration or data value for this struct.
	Page int
	// 中文：PageSize 保存当前结构中的配置或数据值。
	// English: PageSize stores a configuration or data value for this struct.
	PageSize int
	// 中文：Event 保存当前结构中的配置或数据值。
	// English: Event stores a configuration or data value for this struct.
	Event string
	// 中文：RecipientID 保存当前结构中的配置或数据值。
	// English: RecipientID stores a configuration or data value for this struct.
	RecipientID string
	// 中文：UnreadOnly 保存当前结构中的配置或数据值。
	// English: UnreadOnly stores a configuration or data value for this struct.
	UnreadOnly bool
}

// 中文：NotificationRepository 定义当前包使用的数据结构或接口。
// English: NotificationRepository defines a data structure or interface used by this package.
// NotificationRepository persists and queries inbox notifications.
type NotificationRepository struct {
	// 中文：tdb 保存当前结构中的配置或数据值。
	// English: tdb stores a configuration or data value for this struct.
	tdb *orm.TenantDB
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
}

// 中文：NewNotificationRepository 创建并返回对应组件实例。
// English: NewNotificationRepository creates and returns the corresponding component instance.
func NewNotificationRepository(tdb *orm.TenantDB, db *orm.DB) *NotificationRepository {
	return &NotificationRepository{tdb: tdb, db: db}
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (r *NotificationRepository) Create(ctx context.Context, msg *model.Notification) error {
	return r.tdb.Create(ctx, msg)
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (r *NotificationRepository) List(ctx context.Context, filter NotificationFilter) ([]*model.Notification, int64, error) {
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := r.scopedQuery(ctx, filter).Model(&model.Notification{})
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*model.Notification, 0, pageSize)
	err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// 中文：MarkRead 执行当前包中的对应流程。
// English: MarkRead executes the corresponding workflow in this package.
func (r *NotificationRepository) MarkRead(ctx context.Context, id, recipientID string) error {
	if id == "" {
		return gorm.ErrRecordNotFound
	}
	now := time.Now()
	db := r.db.DB().WithContext(ctx).Model(&model.Notification{}).Where("id = ?", id)
	if tenantID := types.GetTenantID(ctx); tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	if recipientID != "" {
		db = db.Where("recipient_id = ?", recipientID)
	}
	result := db.Update("read_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 中文：scopedQuery 执行当前包中的对应流程。
// English: scopedQuery executes the corresponding workflow in this package.
func (r *NotificationRepository) scopedQuery(ctx context.Context, filter NotificationFilter) *gorm.DB {
	db := r.db.DB().WithContext(ctx)
	if tenantID := types.GetTenantID(ctx); tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	if filter.Event != "" {
		db = db.Where("event = ?", filter.Event)
	}
	if filter.RecipientID != "" {
		db = db.Where("recipient_id = ?", filter.RecipientID)
	}
	if filter.UnreadOnly {
		db = db.Where("read_at IS NULL")
	}
	return db
}
