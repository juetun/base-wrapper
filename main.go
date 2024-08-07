/**
* @Author:changjiang
* @Description:
* @File:main
* @Version: 1.0.0
* @Date 2020/4/19 10:19 下午
 */
// //go:generate statik -src=./web/views

// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package main

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/app_start/micro_register"
	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
	"github.com/juetun/base-wrapper/lib/authorization/model"
	. "github.com/juetun/base-wrapper/lib/plugins" // 组件目录
	"github.com/juetun/base-wrapper/lib/plugins/short_message_impl"
	_ "github.com/juetun/base-wrapper/web/router" // 加载路由

	_ "github.com/juetun/base-wrapper/docs"
)

type Authorization struct {
}

func (a *Authorization) Load() (res map[string][]model.AdminAuthorization, err error) {
	res = map[string][]model.AdminAuthorization{}
	return
}

var authorization Authorization

// https://github.com/izghua/go-blog
func main() {
	app_start.NewPlugins(app_start.Authorization(&authorization)).Use(
		PluginRegistry,
		PluginClickHouse,
		PluginOss,
		PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		PluginElasticSearchV8,
		short_message_impl.PluginShortMessage,
		PluginAppMap,
		PluginAuthorization,
		// func(arg *app_start.PluginsOperate) (err error) {
		// 	// 启动websocket
		// 	go anvil_websocket.WebsocketStart()
		// 	return
		// },
		// plugins.PluginOss,
	).LoadPlugins() // 加载插件动作


	// 启动GIN服务
	_ = app_start.NewWebApplication().
		SetFlagMicro(micro_register.NewETCDRegister()).
 		Run()

}
