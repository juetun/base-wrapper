// Package app_start
package app_start

import (
	"github.com/juetun/base-wrapper/lib/plugins"
	stytemLog "log"
	"strings"

	"github.com/juetun/base-wrapper/lib/authorization/model"

	"github.com/juetun/base-wrapper/lib/base"
)

var (
	// PluginsHandleStruct 插件结构体装载对象 ,采用指针共享数据
	PluginsHandleStruct = make([]PluginHandleStruct, 0, 15)

	//定时任务
	TimerTaskHandler []PluginHandleFunction
)

type (
	PluginHandleFunction func(arg *PluginsOperate) (err error)

	// PluginHandleStruct 加载插件结构体
	PluginHandleStruct struct {
		FuncHandler PluginHandleFunction
		Name        string
	}

	AuthorizationStruct interface {
		Load() (map[string][]model.AdminAuthorization, error)
	}

	PluginsOperate struct {
		Author AuthorizationStruct
	}
	PluginsOperateOptionHandler func(arg *PluginsOperate)
)

func NewPlugins(option ...PluginsOperateOptionHandler) (res *PluginsOperate) {
	res = &PluginsOperate{}
	for _, handler := range option {
		handler(res)
	}
	return
}

func Authorization(authorization AuthorizationStruct) (handler PluginsOperateOptionHandler) {
	return func(arg *PluginsOperate) {
		arg.Author = authorization
	}
}

// LoadPlugins 执行加载插件过程
func (r *PluginsOperate) LoadPlugins() (res *PluginsOperate) {
	res = r
	if len(PluginsHandleStruct) == 0 {
		return
	}
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	stytemLog.Printf("")
	stytemLog.Printf("----开始加载插件 ----")
	stytemLog.Printf("")
	var err error
	for _, handle := range PluginsHandleStruct {
		io.SystemOutPrintf(strings.Repeat("-", 30) + "\n")

		if err = handle.FuncHandler(r); err != nil {
			panic(err)
		}
	}

	//加载执行定时任务
	if err = plugins.PluginTimerTask(res); err != nil {
		stytemLog.Printf("----加载定时任务失败 %v----", err.Error())
		return
	}

	io.SystemOutPrintf("Start load plugins finished \n")
	stytemLog.Printf("")
	stytemLog.Printf("----插件加载完成----")
	stytemLog.Printf("")
	return
}

//加载普通插件动作
func (r *PluginsOperate) Use(pluginFunc ...PluginHandleFunction) (res *PluginsOperate) {
	res = r
	Use(pluginFunc...)
	return
}

//加载定时任务操作
func (r *PluginsOperate) UseTimerTask(pluginFunc ...PluginHandleFunction) (res *PluginsOperate) {
	res = r
	TimerTaskHandler = append(TimerTaskHandler, pluginFunc...)
	return
}

// Use 注册系统插件
func Use(pluginFunc ...PluginHandleFunction) {
	for _, value := range pluginFunc {
		PluginsHandleStruct = append(PluginsHandleStruct, PluginHandleStruct{
			FuncHandler: value,
		})
	}
}
