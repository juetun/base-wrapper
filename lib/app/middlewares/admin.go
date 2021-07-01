// Package middlewares
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
)

// AdminMiddlewares 客服后台接口中间件
func AdminMiddlewares() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get(app_obj.HttpHeaderAdminToken) != app_obj.App.AppAdminToken {
			msg := "HTTP_HEADER_ADMIN_TOKEN value is null"
			c.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
			c.Abort()
			return
		}
		c.Next()
	}
}
