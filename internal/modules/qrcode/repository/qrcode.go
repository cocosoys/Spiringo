package repository

import (
	"context"

	"github.com/spiringo/spiringo/internal/modules/qrcode/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	"github.com/spiringo/spiringo/pkg/types"
	"gorm.io/gorm"
)

// 中文：QRCodeRepository 定义当前包使用的数据结构或接口。
// English: QRCodeRepository defines a data structure or interface used by this package.
// QRCodeRepository 二维码数据访问
type QRCodeRepository struct {
	// 中文：tdb 保存当前结构中的配置或数据值。
	// English: tdb stores a configuration or data value for this struct.
	tdb *orm.TenantDB
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
}

// 中文：NewQRCodeRepository 创建并返回对应组件实例。
// English: NewQRCodeRepository creates and returns the corresponding component instance.
// NewQRCodeRepository 创建二维码仓库
func NewQRCodeRepository(tdb *orm.TenantDB, db *orm.DB) *QRCodeRepository {
	return &QRCodeRepository{tdb: tdb, db: db}
}

// 中文：Create 执行当前包中的对应流程。
// English: Create executes the corresponding workflow in this package.
func (r *QRCodeRepository) Create(ctx context.Context, record *model.QRCodeRecord) error {
	return r.tdb.Create(ctx, record)
}

// 中文：GetByID 执行当前包中的对应流程。
// English: GetByID executes the corresponding workflow in this package.
func (r *QRCodeRepository) GetByID(ctx context.Context, id string) (*model.QRCodeRecord, error) {
	var record model.QRCodeRecord
	if err := r.tdb.First(ctx, &record, "id = ?", id); err != nil {
		return nil, err
	}
	return &record, nil
}

// 中文：GetByShortCode 执行当前包中的对应流程。
// English: GetByShortCode executes the corresponding workflow in this package.
func (r *QRCodeRepository) GetByShortCode(ctx context.Context, code string) (*model.QRCodeRecord, error) {
	var record model.QRCodeRecord
	if err := r.tdb.First(ctx, &record, "short_code = ?", code); err != nil {
		return nil, err
	}
	return &record, nil
}

// 中文：Update 执行当前包中的对应流程。
// English: Update executes the corresponding workflow in this package.
func (r *QRCodeRepository) Update(ctx context.Context, record *model.QRCodeRecord) error {
	return r.tdb.Update(ctx, record)
}

// 中文：IncrScanCount 执行当前包中的对应流程。
// English: IncrScanCount executes the corresponding workflow in this package.
func (r *QRCodeRepository) IncrScanCount(ctx context.Context, shortCode string) error {
	db := r.db.DB().WithContext(ctx).Model(&model.QRCodeRecord{}).Where("short_code = ?", shortCode)
	if tenantID := types.GetTenantID(ctx); tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	return db.UpdateColumn("scan_count", gorm.Expr("scan_count + ?", 1)).Error
}

// 中文：CreateScanLog 执行当前包中的对应流程。
// English: CreateScanLog executes the corresponding workflow in this package.
func (r *QRCodeRepository) CreateScanLog(ctx context.Context, log *model.ScanLog) error {
	return r.tdb.Create(ctx, log)
}
