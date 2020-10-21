/**
* @Author:changjiang
* @Description:
* @File:main
* @Version: 1.0.0
* @Date 2020/4/19 10:19 下午
 */

package main

import (
	"github.com/juetun/base-wrapper/lib/app_start"
	_ "github.com/juetun/base-wrapper/lib/init"    // 加载公共插件项
	. "github.com/juetun/base-wrapper/lib/plugins" // 加载路由信息
	_ "github.com/juetun/base-wrapper/web/router"  // 加载路由信息
)

// https://github.com/izghua/go-blog
func main() {
	app_start.NewPluginsOperate().Use(
		// PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		// PluginElasticSearchV7,
		PluginShortMessage,
		PluginAppMap,
		// plugins.PluginOss,
		// plugins.PluginUser, // 用户登录,jwt等用户信息逻辑处理
	).LoadPlugins() // 加载插件动作

	// 启动GIN服务
	app_start.NewWebApplication().
		LoadRouter(). // 记载gin 路由配置
		Run()

}
