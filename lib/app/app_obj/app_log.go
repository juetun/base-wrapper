// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package app_obj

import (
	"fmt"
	"io"
	systemLog "log"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logApp *AppLog

type AppLog struct {
	Logger *logrus.Logger
	GoPath string `json:"go_path"`
}

func newAppLog() (res *AppLog) {
	logApp = &AppLog{
		Logger: logrus.New(),
		GoPath: os.Getenv("GOPATH"),
	}
	return logApp
}

// 获取日志操作对象
func GetLog() *AppLog {
	if logApp == nil {
		systemLog.Printf("【%s】Log object is nil", "ERROR")
	}
	return logApp
}
func (r *AppLog) getFields() (res logrus.Fields) {
	var file = "-" // 获取当前日志写入时的代码位置 （文件名称，函数名称）
	// 获取上层调用者PC，文件名，所在行	// 拼接文件名与所在行
	if _, codePath, codeLine, ok := runtime.Caller(2); ok {
		file = fmt.Sprintf("%s(line:%d)",
			codePath,
			// runtime.FuncForPC(pc).Name(),
			codeLine)
	}
	if r.GoPath != "" {
		res = logrus.Fields{
			APP_LOG_KEY: App.AppName,
			APP_LOG_LOC: "$GOPATH/" + strings.TrimLeft(file, r.GoPath),
		}
		return
	}
	res = logrus.Fields{
		APP_LOG_KEY: App.AppName,
		APP_LOG_LOC: file,
	}
	return
}

func (r *AppLog) Error(context *gin.Context, data map[string]interface{}, message ...interface{}) {
	fields := r.getFields()
	if context != nil {
		fields[TRACE_ID] = context.GetHeader(HTTP_TRACE_ID)
	}

	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.Logger.WithFields(fields).Error(message)
}
func (r *AppLog) Info(context *gin.Context, data map[string]interface{}, message ...interface{}) {
	fields := r.getFields()
	if context != nil {
		fields[TRACE_ID] = context.GetHeader(HTTP_TRACE_ID)
	}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.Logger.WithFields(fields).Info(message)
}
func (r *AppLog) Debug(context *gin.Context, data map[string]interface{}, message ...interface{}) {
	fields := r.getFields()
	if context != nil {
		fields[TRACE_ID] = context.GetHeader(HTTP_TRACE_ID)
	}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.Logger.WithFields(fields).Debug(message)
}
func (r *AppLog) Fatal(context *gin.Context, data map[string]interface{}, message ...interface{}) {
	fields := r.getFields()
	if context != nil {
		fields[TRACE_ID] = context.GetHeader(HTTP_TRACE_ID)
	}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.Logger.WithFields(fields).Fatal(message)
}
func (r *AppLog) Warn(context *gin.Context, data map[string]interface{}, message ...interface{}) {
	fields := r.getFields()
	if context != nil {
		fields[TRACE_ID] = context.GetHeader(HTTP_TRACE_ID)
	}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.Logger.WithFields(fields).Warn(message)
}

// 初始化日志操作对象
func InitAppLog() {
	systemLog.Printf("【INFO】初始化日志操作对象")
	defer systemLog.Printf("【INFO】初始化日志操作对象操作完成 对象内容为:%#v \n", logApp)
	logApp = newAppLog()

	// 获取日志配置
	logConfig := GetLogConfig()

	// 设置日志输出格式
	logFormatter(logConfig, logApp.Logger)

	// 设置日志输出位置
	outputWriter(logConfig, logApp.Logger)

	// 日志收集等级
	logApp.Logger.SetLevel(logrus.Level(logConfig.LogCollectLevel))

}
func logFormatter(logConfig *OptionLog, log *logrus.Logger) {
	switch strings.ToLower(logConfig.Format) { // 日志格式
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{})
	}

}
func outputWriter(config *OptionLog, log *logrus.Logger) {
	var ioWriter []io.Writer
	for _, value := range config.Outputs {
		switch strings.ToLower(value) {
		case "stdout":
			ioWriter = append(ioWriter, os.Stdout)
		case "file":
			if file, err := config.GetFileWriter(); err != nil {
				return
			} else {
				ioWriter = append(ioWriter, file)
			}
		default:
			panic(fmt.Sprintf("当前不支持您配置的日志文件格式(%s)输出", value))
		}
	}

	if len(ioWriter) == 0 { // 默认输出到控制台
		log.SetOutput(os.Stdout)
		return
	}
	// 多个端输出
	log.SetOutput(io.MultiWriter(ioWriter...))
	return
}
