package mysql

import (
	"fmt"

	"gorm.io/plugin/dbresolver"
)

// AutoMigrate 迁移
func (db *DB) AutoMigrate(isExecute bool, dst ...interface{}) {
	if isExecute {
		err := Db.
			Clauses(dbresolver.Write).
			Set("gorm:table_options", "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").
			AutoMigrate(dst...)
		if err != nil {
			panic(fmt.Errorf("tables: [%s] autoMigrate failed: %w", dst, err))
		}
	}
}
