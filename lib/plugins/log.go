package plugins

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"os"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginLog(arg  *app_start.PluginsOperate) (err error) {
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

	io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("Load log start")

	var yamlFile []byte
	if yamlFile, err = os.ReadFile(common.GetConfigFilePath("log.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &mysqlConfig); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load log config config err(%#v) \n", err)
	}

	return
}
