package init

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	. "github.com/juetun/base-wrapper/lib/plugins"
)

// 初始化加载内容
func init() {
	app_start.NewPluginsOperate().Use(
		PluginsApp, // 加载系统配置信息
		PluginLog,  // 加载日志插件
		PluginsHashId,
		PluginMysql, // 加载数据库插件
		PluginRedis, // 加载Redis插件
	)

}
