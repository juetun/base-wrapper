// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_start"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/middlewares"
	"github.com/juetun/base-wrapper/web/controllers"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc,
		func(r *gin.Engine, urlPrefix string) {
			page := controllers.NewControllerPage()
			p := r.Group(urlPrefix, middlewares.AuthParse())
			p.GET("/page", page.Main)
			p.GET("/test", page.Tsst)
			p.GET("/ws", common.GinWebsocketHandler(page.Websocket))
		})
}
