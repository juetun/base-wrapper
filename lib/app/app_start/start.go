package app_start

import (
	"github.com/juetun/base-wrapper/lib/authorization/model"
	stytemLog "log"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
)

type PluginHandleFunction func(arg *PluginsOperate) (err error)

// 加载插件结构体
type PluginHandleStruct struct {
	FuncHandler PluginHandleFunction
	Name        string
}

// 插件结构体装载对象 ,采用指针共享数据
var PluginsHandleStruct = &[]PluginHandleStruct{}

type AuthorizationStruct interface {
	Load() (map[string][]model.AdminAuthorization, error)
}
type PluginsOperate struct {
	Author AuthorizationStruct
}

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

type PluginsOperateOptionHandler func(arg *PluginsOperate)

// 执行加载插件过程
func (r *PluginsOperate) LoadPlugins() (res *PluginsOperate) {
	res = r
	if len(*PluginsHandleStruct) == 0 {
		return
	}
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	stytemLog.Printf("")
	stytemLog.Printf("----开始加载插 ----")
	stytemLog.Printf("")
	var err error
	for _, handle := range *PluginsHandleStruct {
		io.SystemOutPrintf(strings.Repeat("-", 30) + "\n")

		if err = handle.FuncHandler(r); err != nil {
			panic(err)
		}
	}
	io.SystemOutPrintf("Start load plugins finished \n")
	stytemLog.Printf("")
	stytemLog.Printf("----插件加载完成----")
	stytemLog.Printf("")
	return
}
func (r *PluginsOperate) Use(pluginFunc ...PluginHandleFunction) (res *PluginsOperate) {
	res = r
	Use(pluginFunc...)
	return
}

// 注册系统插件
func Use(pluginFunc ...PluginHandleFunction) {
	for _, value := range pluginFunc {
		*PluginsHandleStruct = append(*PluginsHandleStruct, PluginHandleStruct{
			FuncHandler: value,
		})
	}
}
