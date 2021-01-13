/**
* @Author:changjiang
* @Description:
* @File:jwt
* @Version: 1.0.0
* @Date 2020/3/28 6:40 下午
 */
package plugins

import (
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

func PluginJwt() (err error) {

	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	io.SystemOutPrintln("init JWT config")
	defer io.SystemOutPrintln("Init JWT finished")
	err = app_obj.NewJwtParam().JwtInit()
	return
}
