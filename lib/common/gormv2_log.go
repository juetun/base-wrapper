// Package common
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package common

import (
	"context"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger *logrus.Logger
}

func NewWithLogger(logger *logrus.Logger) logger.Interface {
	return &GormLogger{
		Logger: logger,
	}
}
func convertLevel(level logger.LogLevel) logrus.Level {
	switch level {
	case logger.Silent:
		return logrus.PanicLevel // No silent equivalent in logrus
	case logger.Error:
		return logrus.ErrorLevel
	case logger.Warn:
		return logrus.WarnLevel
	case logger.Info:
		return logrus.InfoLevel
	default:
		return logrus.InfoLevel
	}
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.Logger.SetLevel(convertLevel(level))
	return l
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithField("data", data).WithField(app_obj.TraceId, ctx.Value(app_obj.TraceId)).Info(msg)
}

// Warn print warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithField("data", data).WithField(app_obj.TraceId, ctx.Value(app_obj.TraceId)).Warn(msg)
}

// Error print error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithField("data", data).WithField(app_obj.TraceId, ctx.Value(app_obj.TraceId)).Error(msg)
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	f := logrus.Fields{
		"sql":     sql,
		"rows":    rows,
		"elapsed": elapsed,
	}
	if tId := ctx.Value(app_obj.TraceId); tId != nil {
		f[app_obj.TraceId] = tId
	}
	
	if err != nil {
		f["err"] = err.Error()
		l.Logger.WithFields(f).Error("GormLogger")
	} else {
		l.Logger.WithFields(f).Info("GormLogger")
	}
}
