package plugins

import (
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/juetun/base-wrapper/lib/common"
)

func PluginLog() (err error) {
	var io = common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	io.SystemOutPrintln("init log system")
	app_log.InitAppLog()
	io.SystemOutPrintln("Init log finished")

	return
}
