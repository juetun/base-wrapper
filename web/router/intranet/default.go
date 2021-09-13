// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package intranet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	con_impl2 "github.com/juetun/base-wrapper/web/cons/outernet/con_impl"
	"github.com/juetun/base-wrapper/web/cons/page/con_impl"
)

func init() {
	app_start.HandleFuncIntranet = append(app_start.HandleFuncIntranet,
		func(r *gin.Engine, urlPrefix string) {
			page := con_impl.NewConPage()
			p := r.Group(urlPrefix, middlewares.AuthParse())
			p.GET("/page", page.Main)
			p.GET("/test", page.Tsst)

			// websocket操作
			p.GET("/ws", page.Websocket)

			con := con_impl2.NewConDefault()
			p.GET("/index", con.Index)
		})
}
