// Package plugins /**
package plugins

import (
	"fmt"
	"os"
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

	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("Load AppMap start")
	defer base.Io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf(fmt.Sprintf("Load appMap config finished \n"))

	var yamlFile []byte
	if yamlFile, err = os.ReadFile(common.GetCommonConfigFilePath("appmap.yaml", true)); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("yamlFile.Get err #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &app_obj.AppMap); err != nil {
		base.Io.SetInfoType(base.LogLevelError).SystemOutFatalf("Load  appMap config failure(%#v) \n", err)
	}
	for key, value := range app_obj.AppMap {
		base.Io.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("【%s】 %#v", key, value)
	}
	return

}
