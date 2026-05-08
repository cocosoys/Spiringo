package orm

import "gorm.io/driver/postgres"

// 中文：newPostgresDialector 执行当前包中的对应流程。
// English: newPostgresDialector executes the corresponding workflow in this package.
func newPostgresDialector(dsn string) Dialector {
	return postgres.Open(dsn)
}
