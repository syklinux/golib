package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

/*事务*/

type Operator func(tx *gorm.DB) error

func Transaction(o Operator, ctx ...context.Context) error {
	var err error
	// 开启事务
	var tx *gorm.DB
	if len(ctx) > 0 {
		tx = Db.WithContext(ctx[0]).Begin()
	} else {
		tx = Db.Begin()
	}

	if err = tx.Error; err != nil {
		return fmt.Errorf("transaction begin error! :%w", err)
	}

	// 执行事务
	if Err := o(tx); Err != nil {
		// 回滚
		if err := tx.Rollback().Error; err != nil {
			return fmt.Errorf("transaction rollback error! :%w", err)
		}

		return Err
	}

	// 事务提交
	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit error! :%w", err)
	}

	return nil
}
