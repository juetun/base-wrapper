package plugins

import (
	"fmt"
	"sync"

	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/spf13/viper"
)

func PluginLog() (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	io.SystemOutPrintln("init log system")

	// 读取配置文件初始化日志配置数据
	configLogFile, err := loadLogConfig()

	if err != nil {
		return
	}

	app_obj.InitConfig(&configLogFile)

	// 初始化日志操作对象
	app_obj.InitAppLog()
	io.SystemOutPrintln("Init log finished")
	return
}
func loadLogConfig() (mysqlConfig app_obj.OptionLog, err error) {

	io.SystemOutPrintln("Load database start")

	configSource := viper.New()
	configSource.SetConfigName("log")                           // name of config file (without extension)
	configSource.SetConfigType("yaml")                          // REQUIRED if the config file does not have the extension in the name
	configSource.AddConfigPath(common.GetConfigFileDirectory()) // path to look for the config file in

	err = configSource.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error database file: %v \n", err))
		return
	}

	if err = configSource.Unmarshal(&mysqlConfig); err != nil {
		io.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("Load database config failure  '%v' ", mysqlConfig)
		panic(err)
	}
	return
}
