package plugins

import (
	"io/ioutil"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
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


 	var yamlFile []byte

	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("log.yaml")); err != nil {
		io.SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &mysqlConfig); err != nil {
		io.SystemOutFatalf("Load log config config err(%#v) \n", err)
	}
 	io.SystemOutPrintf("Load log config is : '%#v' ", mysqlConfig)

	return
}
