// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/web/cons/con_impl"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet,
		func(r *gin.Engine, urlPrefix string) {
			page := con_impl.NewConPage()
			p := r.Group(urlPrefix, middlewares.AuthParse())
			p.GET("/page", page.Main)
			p.GET("/page/test", page.Tsst)
			p.GET("/page/ws", common.GinWebsocketHandler(page.Websocket))
		})
}
