// Package init
package init

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	. "github.com/juetun/base-wrapper/lib/plugins"
)

// 初始化加载内容，注册必须使用的组件
func init() {
	app_start.Use(
		PluginsApp,    // 加载系统配置信息
		PluginLog,     // 加载日志插件
		PluginsHashId, // 用于生成数据表唯一数据ID
		PluginMysql,   // 加载数据库插件
		PluginRedis,   // 加载Redis插件
	)

}
