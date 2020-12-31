/**
* @Author:changjiang
* @Description:
* @File:main
* @Version: 1.0.0
* @Date 2020/4/19 10:19 下午
 */
////go:generate statik -src=./web/views
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	_ "github.com/juetun/base-wrapper/lib/init"    // 加载公共插件项
	. "github.com/juetun/base-wrapper/lib/plugins" // 组件目录
	_ "github.com/juetun/base-wrapper/web/router"  // 加载路由

	_ "github.com/juetun/base-wrapper/docs"
)

// https://github.com/izghua/go-blog
func main() {
	app_start.NewPluginsOperate().Use(
		PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		// PluginElasticSearchV7,
		PluginShortMessage,
		PluginAppMap,
		// plugins.PluginOss,
	).LoadPlugins() // 加载插件动作

	// 启动GIN服务
	app_start.NewWebApplication().
		LoadRouter(func(r *gin.Engine) (err error) {
			r.LoadHTMLGlob("web/views/**/*.htm")
			r.Static("/static/home", "./static/home")
			r.Static("/static/car", "./static/car")
			r.StaticFile("/jd_root.txt", "./static/jd_root.txt")
			r.StaticFile("/favicon.ico", "./static/favicon.ico")
			return
		}). // 记载gin 路由配置
		Run()

}
