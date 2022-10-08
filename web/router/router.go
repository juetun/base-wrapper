// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	_ "github.com/juetun/base-wrapper/web/router/admin"    // 加载超管访问路由
	_ "github.com/juetun/base-wrapper/web/router/intranet" // 加载内网访问路由
	_ "github.com/juetun/base-wrapper/web/router/outernet" // 加载外网访问路由
	_ "github.com/juetun/base-wrapper/web/router/page"     // 加载网页访问路由
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	app_start.HandleFuncPage = append(app_start.HandleFuncPage, func(r *gin.Engine, urlPrefix string) {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.LoadHTMLGlob("web/views/**/*.htm")
		r.Static("/static/home", "./static/home")
		r.Static("/static/car", "./static/car")
		r.StaticFile("/jd_root.txt", "./static/jd_root.txt")
		r.StaticFile("/favicon.ico", "./static/favicon.ico")
	})
}

