package gormv2

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"pinnacle-primary-be/core/logx"
	"time"
)

type Logger struct {
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func NewLogger() Logger {
	return Logger{
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

func (l Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	logx.WithContext(ctx).Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	logx.Alert(fmt.Sprintf(str, args...))
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	logx.WithContext(ctx).Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	logger := logx.WithContext(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		logger.Errorw("gorm", logx.Field("err", err), logx.Field("elapsed", elapsed), logx.Field("rows", rows), logx.Field("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		logger.Sloww("gorm", logx.Field("elapsed", elapsed), logx.Field("rows", rows), logx.Field("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		logger.Debugw("gorm", logx.Field("elapsed", elapsed), logx.Field("rows", rows), logx.Field("sql", sql))
	}
}
