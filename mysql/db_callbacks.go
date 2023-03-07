package mysql

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/syklinux/golib/log"

	"gorm.io/gorm"
)

const (
	BeforeCreate = iota // gorm:create 之前
	BeforeUpdate        // gorm:update 之前

	AfterCreate // gorm:create 之后
	AfterQuery  // gorm:query 之后
	AfterDelete // gorm:delete 之后

	BeforeCreateAfter // 位于 gorm:before_create 之后 gorm:create 之前

	BeforeAllCreate // 所有其它 callback 之前
	AfterAllCreate  // 所有其它 callback 之后
)

const (
	TimeFormat = "2006-01-02 15:04:05.000"
)

type fn func(db *gorm.DB)

// RegisterCallbacks 注册回调
func (db *DB) RegisterCallbacks(cb ...int) {
	if len(cb) == 0 {
		return
	}

	var err error
	for _, v := range cb {
		_v := v
		switch _v {
		case BeforeCreate:
			err = Db.
				Callback().
				Create().
				Before("gorm:create").
				Register("update_created_at", updateTimeStampForCreateCallback())
		case BeforeUpdate:
			err = Db.
				Callback().
				Update().
				Before("gorm:update").
				Register("my_plugin:before_update", updateTimeStampForUpdateCallback())
		}
	}

	if err != nil {
		panic(fmt.Errorf("registration callback failed: %w", err))
	}
}

func updateTimeStampForCreateCallback() fn {
	return func(db *gorm.DB) {
		if db.Statement.Schema == nil {
			return
		}

		for _, field := range db.Statement.Schema.Fields {
			name := field.Name
			value := db.Statement.ReflectValue
			switch value.Kind() {
			case reflect.Struct:
				now := time.Now().Format(TimeFormat)
				switch name {
				case "CreatedAt":
					if err := field.Set(context.TODO(), value, now); err != nil {
						log.Errorf("[BeforeCreate] created_at: %v", err)
						goto EXIT
					}
				case "UpdatedAt":
					if err := field.Set(context.TODO(), value, now); err != nil {
						log.Errorf("[BeforeCreate] updated_at: %v", err)
						goto EXIT
					}
				}
			case reflect.Slice:
				for i := 0; i < value.Len(); i++ {
					v := value.Index(i)
					now := time.Now().Format(TimeFormat)
					switch name {
					case "CreatedAt":
						if err := field.Set(context.TODO(), v, now); err != nil {
							log.Errorf("[BeforeCreate] created_at: %v", err)
							goto EXIT
						}
					case "UpdatedAt":
						if err := field.Set(context.TODO(), v, now); err != nil {
							log.Errorf("[BeforeCreate] updated_at: %v", err)
							goto EXIT
						}
					}
				}
			}
		}

	EXIT:
		return
	}
}

func updateTimeStampForUpdateCallback() fn {
	return func(db *gorm.DB) {
		if db.Statement.Schema == nil {
			return
		}

		field := db.Statement.Schema.LookUpField("UpdatedAt")
		value := db.Statement.ReflectValue
		switch value.Kind() {
		case reflect.Struct:
			now := time.Now().Format(TimeFormat)
			if err := field.Set(context.TODO(), value, now); err != nil {
				log.Errorf("[BeforeUpdate] updated_at: %v", err)
				return
			}
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				v := value.Index(i)
				now := time.Now().Format(TimeFormat)
				if err := field.Set(context.TODO(), v, now); err != nil {
					log.Errorf("[BeforeUpdate] updated_at: %v", err)
					break
				}
			}
		}
	}
}
