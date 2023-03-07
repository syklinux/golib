package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/syklinux/golib/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

var DefaultSlowThreshold = time.Second

type PluginGorm struct {
	logLevel       zapcore.Level
	slowThreshold  time.Duration
	traceWithLevel zapcore.Level
}

func caller(calldepth int, short bool) string {
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	} else if short {
		file = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// gorm

func NewGorm(logLevel zapcore.Level, traceWithLevel zapcore.Level, slowThreshold ...time.Duration) PluginGorm {
	var slow time.Duration
	if len(slowThreshold) > 0 {
		slow = slowThreshold[0]
	} else {
		slow = DefaultSlowThreshold
	}
	return PluginGorm{
		logLevel:       logLevel,
		slowThreshold:  slow,
		traceWithLevel: traceWithLevel,
	}
}

var logLevelMap = map[logger.LogLevel]zapcore.Level{
	logger.Info:  zap.InfoLevel,
	logger.Warn:  zap.WarnLevel,
	logger.Error: zap.ErrorLevel,
}

func (p PluginGorm) LogMode(level logger.LogLevel) logger.Interface {
	zapLevel, exists := logLevelMap[level]
	if !exists {
		zapLevel = zap.DebugLevel
	}

	newLogger := p
	newLogger.logLevel = zapLevel
	newLogger.slowThreshold = time.Second
	return &newLogger
}

func (p PluginGorm) Info(ctx context.Context, msg string, data ...interface{}) {
	if p.logLevel <= zap.InfoLevel {
		traceId := AttachTraceId(ctx)
		msg = traceId + " " + msg
		log.Infof(msg, data...)
	}
}

func (p PluginGorm) Warn(ctx context.Context, msg string, data ...interface{}) {
	if p.logLevel <= zap.WarnLevel {
		traceId := AttachTraceId(ctx)
		msg = traceId + " " + msg
		log.Warningf(msg, data...)
	}
}

func (p PluginGorm) Error(ctx context.Context, msg string, data ...interface{}) {
	if p.logLevel <= zap.ErrorLevel {
		traceId := AttachTraceId(ctx)
		msg = traceId + " " + msg
		log.Errorf(msg, data...)
	}
}

func (p PluginGorm) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	traceId := AttachTraceId(ctx)
	elapsed := time.Since(begin)
	cost := fmt.Sprintf(" [%.2fms] ", float64(elapsed.Nanoseconds()/1e4)/100.0)

	sql, rows := fc()
	if rows < 0 {
		rows = 0
	}
	sql = fmt.Sprintf("%s %s %s [%d rows affected or returned ]", traceId, cost, sql, rows)
	switch {
	case err != nil:
		log.Errorf("%s err:%s", sql, err.Error())
	case p.slowThreshold != 0 && elapsed > p.slowThreshold:
		ts := fmt.Sprintf(" [%.2fms] ", float64(p.slowThreshold.Nanoseconds()/1e4)/100.0)
		log.Warningf("%s threshold:%s", sql, ts)
	default:
		if p.traceWithLevel == zap.DebugLevel {
			log.Debugf("%s", sql)
		}
		//if p.traceWithLevel == zap.InfoLevel {
		//	log.Infof("%s", sql)
		//} else if p.traceWithLevel == zap.WarnLevel {
		//	log.Warningf("%s", sql)
		//} else if p.traceWithLevel == zap.ErrorLevel {
		//	log.Errorf("%s", sql)
		//}
		//log.Debugf("%s", sql)
	}
}
