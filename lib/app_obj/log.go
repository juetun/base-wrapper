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
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var logConfig *LogConfig

type LogConfig struct {
	Outputs         []string     `json:"outputs"`           // []string{"stdout","file"}
	LogFileConfig                                           // 配置文件信息. 当Outputs值中有输入出到文件标记"file"时有效.
	Format          string       `json:"format"`            // 日志格式 "json"
	LogCollectLevel logrus.Level `json:"log_collect_level"` // 日志收集等级
}

// 如果是
type LogFileConfig struct {
	LogFilePath string `json:"log_file_path"` // 日志文件输出路径，空 不输出
	LogFileName string `json:"log_file_name"` // 日志文件名(或文件名前缀)，空 不输出
	LogIsCut    bool   `json:"log_is_cut"`    // 日志文件是否切割
}

func GetLogConfig() (config *LogConfig) {
	return logConfig
}
func InitConfig() {
	var once sync.Once
	// 只初始化一次
	once.Do(func() {
		dir, _ := os.Getwd()
		logConfig = &LogConfig{
			// Outputs: []string{"stdout"},
			Outputs: []string{"file", "stdout"},
			LogFileConfig: LogFileConfig{
				LogFilePath: dir,
				LogFileName: "log.log",
				LogIsCut:    false,
			},
			Format:          "json",
			// LogCollectLevel: logrus.WarnLevel,
			LogCollectLevel: logrus.InfoLevel,
		}
		systemLog.Printf("【INFO】log config: %#v", logConfig)
	})
	return
}

// 拼接日志文件字符串
func (r *LogConfig) GetFileName(suffix ...string) (res string) {
	res = fmt.Sprintf("%s/%s%s.log", r.LogFilePath, strings.TrimSuffix(r.LogFileName, ".log"), strings.Join(suffix, ""))
	systemLog.Printf("【INFO】log file Name is '%s' ", res)
	return
}

func (r *LogConfig) GetFileWriter() (file io.Writer, err error) {
	logFile := r.GetFileName()
	if logFile == "" {
		systemLog.Printf("【WARN】log file Name is empty! ")
		return
	}
	if r.LogIsCut { // 如果日志文件需要切割
		file, err = rotatelogs.New(r.GetFileName("_%Y%m%d"),
			rotatelogs.WithLinkName(r.LogFilePath),
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithMaxAge(24*100*time.Hour),
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
