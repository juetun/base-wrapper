// Package plugins
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package plugins

import (
	"os"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginsApp(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	app_obj.App = common.NewApplication()
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("Load app config start")
	defer func() {
		base.Io.SetInfoType(base.LogLevelInfo).SetInfoType(base.LogLevelInfo).SystemOutPrintf("app config \n")
		for key, value := range app_obj.App.ToMap() {
			base.Io.SetInfoType(base.LogLevelInfo).SetInfoType(base.LogLevelInfo).SystemOutPrintf("\t【%s】'%#v' \n", key, value)
		}
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("load app config finished \n")
	}()
	dir := common.GetConfigFileDirectory()
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("config directory is : '%s' ", dir)

	type app struct {
		App *app_obj.Application `json:"app" yaml:"app"`
	}
	var data = app{App: app_obj.App}
	data.App.Default()
	data.App.AppTemplateDirectory = common.DefaultAppTemplateDirectory( )
	var yamlFile []byte
	filePath := common.GetConfigFilePath("app.yaml", true)
	if yamlFile, err = os.ReadFile(filePath); err != nil {
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &data); err != nil {
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutFatalf("load app config err(%#v) \n", err)
	}
	return
}
