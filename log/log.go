package log

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

// LogConfig LogConfig
type LogConfig struct {
	Level string `json:"level"` // INFO、DEBUG、ERROR、WARN
	Type  string `json:"type"`  // 输出类型 json、text、default
}

type Log struct {
	logger *slog.Logger
	config LogConfig
}

var l Log

// init 忧于main函数之前执行
func init() {
	conf := LogConfig{
		Type:  "json",
		Level: "INFO",
	}
	SetConfig(conf)
	SetLogger(newLogger(conf))
	slog.SetDefault(l.logger)
}

func ExportAttr(args []any) []slog.Attr {
	var (
		attrs []slog.Attr
		other []any
	)
	for _, v := range args {
		switch v.(type) {
		case error:
			attrs = append(attrs, slog.Any("err", v))
		default:
			other = append(other, v)
		}
	}
	attrs = append(attrs, slog.Any("other", other))
	return attrs
}

// InitByConf 通过conf 初始化log
func InitByConf(conf LogConfig) error {
	SetConfig(conf)
	SetLogger(newLogger(conf))
	slog.SetDefault(l.logger)
	return nil
}

func stringToLevel(level string) slog.Level {
	switch level {
	case "ERROR":
		return slog.LevelError
	case "WARNING":
		return slog.LevelWarn
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	}
	return slog.LevelDebug
}
func SetConfig(conf LogConfig) {
	l.config = conf
}

func SetLogger(logger *slog.Logger) {
	l.logger = logger
}

func newLogger(conf LogConfig) (logger *slog.Logger) {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     stringToLevel(conf.Level),
	}
	switch conf.Type {
	case "text":
		logger = slog.New(opts.NewTextHandler(os.Stderr))
	case "json":
		logger = slog.New(opts.NewJSONHandler(os.Stderr))
	default:
		logger = slog.New(opts.NewJSONHandler(os.Stderr))
	}
	return logger
}

func Debug(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg,
			slog.Any("other", args))
	} else {
		l.logger.Debug(msg)
	}
}

func Debugf(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg,
			slog.Any("other", args))
	} else {
		l.logger.Debug(msg)
	}
}

func Info(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg,
			slog.Any("other", args))
	} else {
		l.logger.Info(msg)
	}
}

func Infof(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg,
			slog.Any("other", args))
	} else {
		l.logger.Info(msg)
	}
}

func Warning(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg,
			slog.Any("other", args))
	} else {
		l.logger.Warn(msg)
	}
}

func Warningf(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg,
			slog.Any("other", args))
	} else {
		l.logger.Warn(msg)
	}
}

func Warn(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg,
			slog.Any("other", args))
	} else {
		l.logger.Warn(msg)
	}
}

func Warnf(msg string, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg,
			slog.Any("other", args))
	} else {
		l.logger.Warn(msg)
	}
}

func Error(msg string, err error, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelError, msg, ExportAttr(args)...)
	} else {
		l.logger.Error(msg, err)
	}
}

func Errorf(msg string, args ...any) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, ExportAttr(args)...)
}

func Fatal(msg string, err error, args ...any) {
	if len(args) > 0 {
		l.logger.LogAttrs(context.Background(), slog.LevelError, msg, ExportAttr(args)...)
	} else {
		l.logger.Error(msg, err)
	}
	os.Exit(1)
}

func Fatalf(msg string, args ...any) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, ExportAttr(args)...)
	os.Exit(1)
}

func Panicf(msg string, args ...any) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, ExportAttr(args)...)
	os.Exit(1)
}
