// Package plugins /**
package plugins

import (
	"github.com/juetun/base-wrapper/lib/base"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"
)

func PluginJwt(arg *app_start.PluginsOperate) (err error) {

	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("init JWT config")
	defer  base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("Init JWT finished")
	err = base.NewJwtParam().JwtInit()
	return
}
