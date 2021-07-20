// Package common
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package common

import (
	"context"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger *app_obj.AppLog
}

func NewWithLogger(logger *app_obj.AppLog) logger.Interface {
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
func (r *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	r.Logger.Logger.SetLevel(convertLevel(level))
	return r
}

// Info print info
func (r *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	dta := map[string]interface{}{"data": data}
	if dt, ok := r.GetContextParameter(ctx); ok {
		dta[app_obj.TraceId] = dt.TraceId
		dta[app_obj.DbNameKey] = dt.DbName
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Info(msg)
	} else {
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Info(msg)
	}
}

// Warn print warn messages
func (r *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {

	dta := map[string]interface{}{"data": data}
	if dt, ok := r.GetContextParameter(ctx); ok {
		dta[app_obj.TraceId] = dt.TraceId
		dta[app_obj.DbNameKey] = dt.DbName
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Warn(msg)
	} else {
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Warn(msg)
	}

}

func (r *GormLogger) GetContextParameter(ctx context.Context) (res base.DbContextValue, ok bool) {
	dt := ctx.Value(app_obj.DbContextValueKey)
	res, ok = dt.(base.DbContextValue)
	return
}

// Error print error messages
func (r *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	dta := map[string]interface{}{"data": data}
	if dt, ok := r.GetContextParameter(ctx); ok {
		dta[app_obj.TraceId] = dt.TraceId
		dta[app_obj.DbNameKey] = dt.DbName
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Error(msg)
	} else {
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Error(msg)
	}
}

// Trace print sql message
func (r *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)

	dta := map[string]interface{}{
		"sql":      sql,
		"rows":     rows,
		"elapsed":  elapsed,
		"duration": float64((time.Now().UnixNano() - begin.UnixNano()) / 1e6), // 时长单位微秒
	}
	if dt, ok := r.GetContextParameter(ctx); ok {
		dta[app_obj.TraceId] = dt.TraceId
		dta[app_obj.DbNameKey] = dt.DbName
	}

	if err != nil {
		dta["err"] = err.Error()
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Error("Gorm")
	} else {
		r.Logger.Logger.WithFields(r.Logger.GetFields(dta)).Info("Gorm")
	}
}
