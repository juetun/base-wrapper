package plugins

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)


//定时任务调度器
func PluginTimerTask(arg *app_start.PluginsOperate) (err error) {
	if !app_obj.App.AppRunTimerTask {
		io.SystemOutPrintln("当前服务将不会执行定时任务")
		return
	}

	defer func() {
		io.SystemOutPrintln("Start timer task running")
	}()
	for _, handler := range app_start.TimerTaskHandler {
		handler(arg)
	}
	io.SystemOutPrintln("Load timer task")

	return
}
