/**
* @Author:changjiang
* @Description:
* @File:jwt
* @Version: 1.0.0
* @Date 2020/3/28 6:40 下午
 */
package plugins

import (
	"github.com/juetun/app-web/lib/app_obj"
)

func PluginJwt() (err error) {
	io.SystemOutPrintln("init JWT config")
	jwtParam := app_obj.NewJwtParam()
	err = jwtParam.JwtInit()
	io.SystemOutPrintln("Init JWT finished")
	return
}
