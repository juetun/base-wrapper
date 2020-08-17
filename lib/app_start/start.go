package app_start

import (
	stytemLog "log"

	"github.com/juetun/base-wrapper/lib/common"
)

type PluginHandleFunction func() (err error)

// 加载插件结构体
type PluginHandleStruct struct {
	FuncHandler PluginHandleFunction
	Name        string
}

// 插件结构体装载对象 ,采用指针共享数据
var PluginsHandleStruct = &[]PluginHandleStruct{}

type PluginsOperate struct {
}

func NewPluginsOperate() *PluginsOperate {
	return &PluginsOperate{}
}

// 执行加载插件过程
func (r *PluginsOperate) LoadPlugins() *PluginsOperate {
	var io = common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	stytemLog.Printf("")
	stytemLog.Printf("------------ Start load plugins-----------------------")
	stytemLog.Printf("")
	var err error
	for _, handle := range *PluginsHandleStruct {
		err = handle.FuncHandler()
		if err != nil {
			panic(err)
		}
	}
	io.SystemOutPrintf("Start load plugins finished \n")
	stytemLog.Printf("")
	stytemLog.Printf("-------------Load plugins finished----------------------")
	stytemLog.Printf("")
	return r
}

// 注册系统插件
func (r *PluginsOperate) Use(pluginFunc ...PluginHandleFunction) *PluginsOperate {
	for _, value := range pluginFunc {
		*PluginsHandleStruct = append(*PluginsHandleStruct, PluginHandleStruct{
			FuncHandler: value,
		})
	}
	return r
}
