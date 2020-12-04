/**
* @Author:changjiang
* @Description:
* @File:route
* @Version: 1.0.0
* @Date 2020/4/19 10:22 下午
 */

// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package outernet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)

func init() {
	app_start.HandleFuncOuterNet = append(app_start.HandleFuncOuterNet,
		func(r *gin.Engine, urlPrefix string) {
			//c := controllers.NewControllerDefault()
			//p := r.Group(urlPrefix, middlewares.Authentication(func(user *app_obj.JwtUserMessage, c *gin.Context) (err error) {
			//	return
			//}))
			//p.GET("/test", c.Index)
			//p.GET("/test_es", c.TestEs)

		}, )
}
