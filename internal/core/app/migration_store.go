package app

import (
	"github.com/spiringo/spiringo/internal/core/migrate"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：newMigrationStore 执行当前包中的对应流程。
// English: newMigrationStore executes the corresponding workflow in this package.
func newMigrationStore(db *orm.DB) module.MigrationStore {
	return migrate.NewStore(db)
}
