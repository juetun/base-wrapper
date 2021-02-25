// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package plugins

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)

// 接口权限授信逻辑
func PluginAuthorization(arg *app_start.PluginsOperate) (err error) {
	if arg.Author == nil {
		err = fmt.Errorf("你没有设置接口验证信息")
		return
	}
	app_obj.Authorization, err = arg.Author.Load()
	return
}
