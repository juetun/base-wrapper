/**
* @Author:changjiang
* @Description:
* @File:route
* @Version: 1.0.0
* @Date 2020/4/19 10:22 下午
 */
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_start"
	"github.com/juetun/base-wrapper/web/controllers"
)

func init() {
	app_start.HandleFunc = append(app_start.HandleFunc,
		func(r *gin.Engine, urlPrefix string) {
			c := controllers.NewControllerDefault()
			page := controllers.NewControllerPage()
			p := r.Group(urlPrefix)
			p.GET("/test", c.Index)
			p.GET("/test_es", c.TestEs)
			p.GET("/page", page.Main)
		}, )
}
