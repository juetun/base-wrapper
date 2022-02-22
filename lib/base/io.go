package base

import (
	"fmt"
	systemLog "log"
	"sync"
)

type SystemOut struct {
	InfoType string `json:"info_type"`
}

const (
	LogLevelError = "ERROR"
	LogLevelInfo  = "INFO"
	LogLevelDebug = "DEBUG"
	LogLevelFatal = "FATAL"
)

var systemLogInit sync.Once
var logObj *SystemOut

func NewSystemOut() *SystemOut {
	systemLogInit.Do(func() {
		logObj = &SystemOut{
			InfoType: LogLevelInfo,
		}
	})
	return logObj
}
func (r *SystemOut) SetInfoType(infoType string) *SystemOut {
	r.InfoType = fmt.Sprintf("【%s】", infoType)
	return r
}
func (r *SystemOut) SystemOutPrintln(s string) *SystemOut {
	systemLog.Println(r.InfoType + s)
	return r
}

func (r *SystemOut) SystemOutFatalf(format string, v ...interface{}) *SystemOut {
	systemLog.Fatalf(r.InfoType+format, v...)
	return r
}

func (r *SystemOut) SystemOutPrintf(format string, v ...interface{}) *SystemOut {
	systemLog.Printf(r.InfoType+format, v...)
	return r
}
