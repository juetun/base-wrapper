/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-13
 * Time: 22:36
 */
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
)

// 加载权限验证Gin中间件
func Permission() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
			c.Next()
			return
		}
		{ // 跨域逻辑添加
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Max-Age", "1800")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE,PATCH")
			// c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding,referrer, Authorization, x-*,X-*")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Next()
	}
}

// 用户登录逻辑处理
func Auth(c *gin.Context) (exit bool) {

	token := c.Request.Header.Get(app_obj.HTTP_USER_TOKEN)
	traceId := c.GetHeader(app_obj.HTTP_TRACE_ID)
	c.Set(app_obj.TRACE_ID, traceId)

	if token == "" {
		msg := "token is null"
		app_log.GetLog().Error(c, map[string]interface{}{
			"method": "zgh.ginmiddleware.Auth",
			"error":  msg,
		})
		c.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
		c.Abort()
		exit = true
		return
	}

	jwtUser, err := common.ParseToken(token, c)
	if err != nil {
		app_log.GetLog().Error(c, map[string]interface{}{
			"method": "zgh.ginmiddleware.Auth",
			"token":  token,
			"error":  err.Error(),
		})
		c.JSON(http.StatusOK, common.NewHttpResult().SetCode(403).SetMessage(err.Error()))
		c.Abort()
		exit = true
		return
	}
	c.Set(app_obj.ContextUserObjectKey, jwtUser)
	c.Set(app_obj.ContextUserTokenKey, token)

	return
}

func GetRUri(c *gin.Context) string {
	uri := strings.TrimLeft(c.Request.RequestURI, common.GetAppConfig().AppName+"/"+common.GetAppConfig().AppApiVersion)
	if uri == "" { // 如果是默认页 ，则直接让过
		return "default"
	}
	s1 := strings.Split(uri, "?")
	s2 := strings.TrimRight(s1[0], "/")
	// fmt.Printf("Uri is :'%v'", s2)
	return s2
}
