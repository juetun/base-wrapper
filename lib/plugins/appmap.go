// Package plugins /**
package plugins

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginAppMap(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	io.SystemOutPrintln("Load AppMap start")
	defer io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf(fmt.Sprintf("Load appMap config finished \n"))

	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("appmap.yaml")); err != nil {
		io.SystemOutFatalf("yamlFile.Get err #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &app_obj.AppMap); err != nil {
		io.SystemOutFatalf("Load  appMap config failure(%#v) \n", err)
	}
	for key, value := range app_obj.AppMap {
		io.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("【%s】 %#v", key, value)
	}
	return

}
