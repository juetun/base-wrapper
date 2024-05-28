// Package plugins
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package plugins

import (
	"fmt"
	"os"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/micro_service"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginRegistry(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	loadRegistryConfig()
	return
}

func loadRegistryConfig() (err error) {

	io.SystemOutPrintln("开始注册服务")

	// 数据库配置信息存储对象
	var yamlFile []byte
	if yamlFile, err = os.ReadFile(common.GetConfigFilePath("registry.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
		return
	}

	if err = yaml.Unmarshal(yamlFile, &micro_service.ServiceConfig); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load micro server registry err(%#v) \n", err)
		return
	}

	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("registry server (%#v) \n", micro_service.ServiceConfig)
	//

	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("load micro server registry finished \n"))
	return
}
