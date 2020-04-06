package app_log

import (
	systemLog "log"
	"os"

	"github.com/sirupsen/logrus"
)

var logApp *AppLog

type AppLog struct {
	*logrus.Logger
}

func newAppLog() *AppLog {
	return &AppLog{}
}

// 获取日志操作对象
func GetLog() *AppLog {
	if logApp == nil {
		systemLog.Printf("【%s】Log object is nil", "ERROR")
	}
	return logApp
}
func (r *AppLog) SetLog(log *logrus.Logger) *AppLog {
	r.Logger = log
	return r
}

func (r *AppLog) Error(data map[string]string, message ...string) {
	fields := logrus.Fields{}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.WithFields(fields).Error(message)
}
func (r *AppLog) Info(data map[string]string, message ...string) {
	fields := logrus.Fields{}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.WithFields(fields).Error(message)
}

// 初始化日志操作对象
func InitAppLog() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.WarnLevel)
	logApp = newAppLog()
	logApp.SetLog(log)
}
