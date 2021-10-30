// Package plugins /**
package plugins

import (
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

func PluginJwt(arg *app_start.PluginsOperate) (err error) {

	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	io.SystemOutPrintln("init JWT config")
	defer io.SystemOutPrintln("Init JWT finished")
	err = app_obj.NewJwtParam().JwtInit()
	return
}
