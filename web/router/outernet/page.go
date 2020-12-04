// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/web/controllers/con_impl"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet,
		func(r *gin.Engine, urlPrefix string) {
			page := con_impl.NewControllerPage()
			p := r.Group(urlPrefix, middlewares.AuthParse())
			p.GET("/page", page.Main)
			p.GET("/test", page.Tsst)
			p.GET("/ws", common.GinWebsocketHandler(page.Websocket))
		})
}
