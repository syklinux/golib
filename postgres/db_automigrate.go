package postgres

import (
	"fmt"

	"gorm.io/plugin/dbresolver"
)

// AutoMigrate 迁移
func (db *DB) AutoMigrate(isExecute bool, dst ...interface{}) {
	if isExecute {
		err := Db.
			Clauses(dbresolver.Write).
			AutoMigrate(dst...)
		if err != nil {
			panic(fmt.Errorf("tables: [%s] autoMigrate failed: %w", dst, err))
		}
	}
}
