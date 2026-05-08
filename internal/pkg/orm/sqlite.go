package orm

import "gorm.io/driver/sqlite"

// 中文：newSQLiteDialector 执行当前包中的对应流程。
// English: newSQLiteDialector executes the corresponding workflow in this package.
func newSQLiteDialector(dsn string) Dialector {
	return sqlite.Open(dsn)
}
