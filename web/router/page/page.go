// Package page
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	"github.com/juetun/base-wrapper/web/cons/page/con_impl"
)

func init() {
	app_start.HandleFuncPage = append(app_start.HandleFuncPage,
		func(r *gin.Engine, urlPrefix string) {
			page := con_impl.NewConPage()
			p := r.Group(urlPrefix, middlewares.AuthParse())
			p.POST("/page", page.Main)
			p.POST("/page_sign", page.MainSign)
			p.GET("/page_sign", page.MainSign)
			p.PUT("/page_sign", page.MainSign)
			p.DELETE("/page_sign/:id", page.MainSign)
			p.GET("/page/test", page.Tsst)

			//websocket 操作
			p.GET("/page/ws", page.Websocket)

		})
}
