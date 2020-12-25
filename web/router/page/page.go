// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/web/cons/page/con_impl"
)

func init() {
	app_start.HandleFuncPage = append(app_start.HandleFuncPage,
		func(r *gin.Engine, urlPrefix string) {
			page := con_impl.NewConPage()
			p := r.Group(urlPrefix, middlewares.AuthParse())
			p.POST("/page", page.Main)
			p.GET("/page/test", page.Tsst)
			p.GET("/page/ws", common.GinWebsocketHandler(page.Websocket))
		})
}
