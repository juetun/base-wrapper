package init

import (
	"github.com/juetun/base-wrapper/lib/app_start"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins"
)

// 初始化加载内容
func init() {
	app_start.NewPluginsOperate().Use(
		common.PluginsApp, // 加载系统配置信息
		plugins.PluginLog, // 加载日志插件
		common.PluginsHashId,
		plugins.PluginMysql, // 加载数据库插件
		plugins.PluginRedis, // 加载Redis插件
	)

}