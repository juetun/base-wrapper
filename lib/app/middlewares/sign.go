/**
* @Author:changjiang
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
	"github.com/juetun/base-wrapper/lib/common/signencrypt"
)

// 接口签名验证
func SignHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var res bool
		var err error
		if res, _, err = signencrypt.NewSign().
			SignGinRequest(c, func(appName string) (secret string, err error) {
				secret = "signxxx"
				// TODO 通过appName获取签名值
				return
			}); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  "sign err",
			})
			return
		}
		if !res {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  "sign validate failure",
			})
		}
		c.Next()
	}
}
