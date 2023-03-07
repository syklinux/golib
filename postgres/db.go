package postgres

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/syklinux/golib/log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var Db *gorm.DB
var once sync.Once

var (
	ErrorNeedInvalidRead = errors.New("provide at least one valid read database address")
)

func InitDb(conf *Conf, isAutoMigrate bool, datas ...interface{}) {
	db := newDB()
	db.InitDataBase(conf, datas...)
	db.RegisterCallbacks(BeforeCreate,
		BeforeUpdate,
		AfterCreate,
		AfterQuery,
		AfterDelete,
		BeforeCreateAfter,
		BeforeAllCreate,
		AfterAllCreate)
	db.AutoMigrate(isAutoMigrate, datas...)
}

type DB struct{}

func newDB() *DB {
	return new(DB)
}

func (db *DB) InitDataBase(conf *Conf, datas ...interface{}) {
	var (
		err error
	)
	cfg := &gorm.Config{}
	level := zap.DebugLevel
	if conf.LogLevel == "info" {
		level = zap.InfoLevel
	}
	if conf.LogLevel == "error" {
		level = zap.ErrorLevel
	}
	if conf.LogLevel == "waring" {
		level = zap.WarnLevel
	}
	cfg.Logger = NewGorm(level, level, time.Second)
	cfg.PrepareStmt = false
	cfg.NowFunc = func() time.Time {
		sh, _ := time.LoadLocation("Asia/Shanghai")
		return time.Now().In(sh)
	}
	cfg.DryRun = false
	cfg.DisableAutomaticPing = false
	Db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: db.getWriteDsn(conf.DataBase, conf.WriteServer, conf.User, conf.Password),
	}), cfg)
	if err != nil {
		panic(err)
	}

	readHosts := conf.ReadServer
	if len(readHosts) == 0 {
		panic(ErrorNeedInvalidRead)
	}
	dials := make([]gorm.Dialector, len(readHosts))
	for index, host := range readHosts {
		dials[index] = postgres.Open(db.getReadDsn(conf.DataBase, host, conf.User, conf.Password))
	}

	// 连接池配置
	maxOpenConn := conf.MaxOpenConn
	maxIdleConn := conf.MaxIdleConn
	maxIdleTime := time.Duration(conf.MaxIdleTime) * time.Second
	maxLiftTime := time.Duration(conf.MaxLiftTime) * time.Second

	// 初始化
	if err := Db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(db.getWriteDsn(conf.DataBase, conf.WriteServer, conf.User, conf.Password))},
		Replicas: dials,
		Policy:   dbresolver.RandomPolicy{},
	}, datas...).
		SetConnMaxIdleTime(maxIdleTime).
		SetConnMaxLifetime(maxLiftTime).
		SetMaxIdleConns(maxIdleConn).
		SetMaxOpenConns(maxOpenConn)); err != nil {
		panic(fmt.Errorf("database initialization failed: %w", err))
	}
}

func (db *DB) getReadDsn(dbname, host, user string, passwd string) string {
	hostPost := strings.Split(host, ":")
	port, _ := strconv.Atoi(hostPost[1])
	a := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d TimeZone=Asia/Shanghai",
		user, passwd, dbname, hostPost[0], port)
	return a
}

func (db *DB) getWriteDsn(dbname, host, user string, passwd string) string {
	hostPost := strings.Split(host, ":")
	port, _ := strconv.Atoi(hostPost[1])
	a := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d TimeZone=Asia/Shanghai",
		user, passwd, dbname, hostPost[0], port)
	return a
}

func Close() error {
	once.Do(func() {
		if Db != nil {
			Db = nil
		}
	})

	log.Infof("[DATABASE] connection closed successfully")
	return nil
}
