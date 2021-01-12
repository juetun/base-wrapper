package app_start

import (
	stytemLog "log"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
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

		if err = handle.FuncHandler();err != nil {
			panic(err)
		}
	}
	io.SystemOutPrintf("Start load plugins finished \n")
	stytemLog.Printf("")
	stytemLog.Printf("----插件加载完成----")
	stytemLog.Printf("")
	return
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
