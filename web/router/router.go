// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package router

import (
	_ "github.com/juetun/base-wrapper/web/router/intranet" // 加载内网访问路由
	_ "github.com/juetun/base-wrapper/web/router/outernet" // 加载外网访问路由
	_ "github.com/juetun/base-wrapper/web/router/admin" // 加载超管访问路由
	_ "github.com/juetun/base-wrapper/web/router/page" // 加载网页访问路由
)
func init()  {

}
