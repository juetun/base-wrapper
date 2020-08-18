package plugins

import (
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
)

func PluginLog() (err error) {
	var io = common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	io.SystemOutPrintln("init log system")

	// 初始化日志配置数据
	app_obj.InitConfig()

	// 初始化日志操作对象
	app_log.InitAppLog()
	io.SystemOutPrintln("Init log finished")
	return
}
