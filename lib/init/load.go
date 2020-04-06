package init

import (
	"github.com/juetun/app-web/lib/app_start"
	"github.com/juetun/app-web/lib/common"
	"github.com/juetun/app-web/lib/plugins"
)

// 初始化加载内容
func init() {
	app_start.NewPluginsOperate().Use(
		common.PluginsApp, // 加载系统配置信息
		common.PluginsHashId,
		plugins.PluginLog,   // 加载日志插件
		plugins.PluginMysql, // 加载数据库插件
		plugins.PluginRedis, // 加载Redis插件
	)

}
