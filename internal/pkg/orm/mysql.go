package orm

import "gorm.io/driver/mysql"

// 中文：newMySQLDialector 执行当前包中的对应流程。
// English: newMySQLDialector executes the corresponding workflow in this package.
func newMySQLDialector(dsn string) Dialector {
	return mysql.Open(dsn)
}
