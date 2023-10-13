// Package middlewares
/**
* @Author:ChangJiang
* @Description:
* @File:sign
* @Version: 1.0.0
* @Date 2021/2/23 9:53 下午
 */
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
)

// SignHttp 接口签名验证
func SignHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodOptions:
			c.Next()
			return
		case http.MethodHead:
			c.Next()
			return
		}
		//websocket的参数验证跳过
		if c.Request.Header.Get("Upgrade") == "websocket" {
			c.Next()
			return
		}
		var res bool
		var err error
		if res, _, err = base.NewSign().
			SignGinRequest(c); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusForbidden,
				Msg:  "sign err",
			})
			return
		}
		if !res {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusForbidden,
				Msg:  "sign validate failure",
			})
		}
		c.Next()
	}
}
