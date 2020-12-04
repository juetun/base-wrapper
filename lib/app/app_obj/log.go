/**
* @Author:changjiang
* @Description:
* @File:log.
* @Version: 1.0.0
* @Date 2020/8/17 11:47 下午
 */
package app_obj

import (
	"fmt"
	"io"
	systemLog "log"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var logConfig *OptionLog

type OptionLog struct {
	Outputs         []string `json:"outputs" yml:"outputs"`                       // []string{"stdout","file"}
	LogFileConfig                                                                 // 配置文件信息. 当Outputs值中有输入出到文件标记"file"时有效.
	Format          string       `json:"format" yml:"format"`                     // 日志格式 "json"
	LogCollectLevel logrus.Level `json:"log_collect_level" yml:"logcollectlevel"` // 日志收集等级
}

// 如果是
type LogFileConfig struct {
	LogFilePath string `json:"log_file_path" yml:"logfilepath"` // 日志文件输出路径，空 不输出
	LogFileName string `json:"log_file_name" yml:"logfilename"` // 日志文件名(或文件名前缀)，空 不输出
	LogIsCut    bool   `json:"log_is_cut" yml:"logiscut"`       // 日志文件是否切割
}

type SetOption func(opt *OptionLog)

func NewOption(setOption ...SetOption) {
	opt := &OptionLog{}
	for _, setOpt := range setOption {
		setOpt(opt)
	}
}
func Outputs(outPuts []string) SetOption {
	return func(opt *OptionLog) {
		opt.Outputs = outPuts
	}
}
func Format(format string) SetOption {
	return func(opt *OptionLog) {
		opt.Format = format
	}
}

func LogCollectLevel(arg logrus.Level) SetOption {
	return func(opt *OptionLog) {
		opt.LogCollectLevel = arg
	}
}
func LogFilePath(arg string) SetOption {
	return func(opt *OptionLog) {
		opt.LogFilePath = arg
	}
}
func LogFileName(arg string) SetOption {
	return func(opt *OptionLog) {
		opt.LogFileName = arg
	}
}
func LogIsCut(arg bool) SetOption {
	return func(opt *OptionLog) {
		opt.LogIsCut = arg
	}
}

func GetLogConfig() (config *OptionLog) {
	return logConfig
}

func InitConfig(config *OptionLog) {
	logConfig = config
	if len(logConfig.Outputs) == 0 {
		logConfig.Outputs = []string{"file", "stdout"}
	}
	if logConfig.LogFilePath == "" {
		dir, _ := os.Getwd()
		logConfig.LogFilePath = dir
	}
	if logConfig.LogFileName == "" {
		logConfig.LogFileName = "log.log"
	}
	if logConfig.Format == "" {
		logConfig.Format = "json"
	}

	//
	systemLog.Printf("【INFO】log config: %#v", logConfig)

	return
}

// 拼接日志文件字符串
func (r *OptionLog) GetFileName(suffix ...string) (res string) {
	res = fmt.Sprintf("%s/%s%s.log", r.LogFilePath, strings.TrimSuffix(r.LogFileName, ".log"), strings.Join(suffix, ""))
	systemLog.Printf("【INFO】log file Name is '%s' ", res)
	return
}

func (r *OptionLog) GetFileWriter() (file io.Writer, err error) {
	logFile := r.GetFileName()
	if logFile == "" {
		systemLog.Printf("【WARN】log file Name is empty! ")
		return
	}
	if r.LogIsCut { // 如果日志文件需要切割
		file, err = rotatelogs.New(r.GetFileName("_%Y%m%d"),
			rotatelogs.WithLinkName(r.LogFilePath),
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithMaxAge(14*24*time.Hour), //日志文件保留14天
		)
		if err != nil {
			systemLog.Printf("failed to create rotatelogs: %s", err)
			return
		}
		return
	}
	// 日志文件不需要切割
	file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		systemLog.Printf("【ERROR】Could Not Open Log File(%s) : %s ", logFile, err.Error())
		return
	}
	return
}
