// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package plugins

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory/traefik/etcd"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

func PluginRegistry(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	loadRegistryConfig()
	return
}

func loadRegistryConfig() (err error) {

	io.SystemOutPrintln("Load micro server registry start")

	// 数据库配置信息存储对象
	var config service_discory.ServerRegistry
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("registry.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load micro server registry err(%#v) \n", err)
	}
	if err = etcd.NewTraefikEtcd(&config).Action(); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load micro server registry err(%#v) \n", err)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("load micro server registry finished \n"))
	return
}
