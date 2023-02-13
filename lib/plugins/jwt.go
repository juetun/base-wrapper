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

	io.SystemOutPrintln("init JWT config")
	defer io.SystemOutPrintln("Init JWT finished")
	err = base.NewJwtParam().JwtInit()
	return
}
