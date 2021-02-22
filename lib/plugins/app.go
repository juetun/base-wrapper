// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package plugins

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

func PluginsApp(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	app_obj.App = common.NewApplication()
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintln("Load app config start")
	defer func() {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("app config is: '%v' \n", app_obj.App.ToString())
		io.SystemOutPrintf("load app config finished \n")
	}()
	dir := common.GetConfigFileDirectory()
	io.SystemOutPrintf("config directory is : '%s' ", dir)

	type app struct {
		App *app_obj.Application `json:"app" yaml:"app"`
	}
	var data = app{App: app_obj.App}
	data.App.Default()
	data.App.AppTemplateDirectory = common.DefaultAppTemplateDirectory(io)
	var yamlFile []byte
	filePath := common.GetConfigFilePath("app.yaml", true)
	if yamlFile, err = ioutil.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &data); err != nil {
		io.SystemOutFatalf("load app config err(%#v) \n", err)
	}

	io.SystemOutPrintf("App Config is : '%#v' ", app_obj.App)
	return
}
