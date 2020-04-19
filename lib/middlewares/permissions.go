/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-13
 * Time: 22:36
 */
package middlewares

import (
	"fmt"
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
		status := cors(c)
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
			c.Next()
			return
		}
		// 跨域配置
		if status {
			return
		}

		// var res bool
		// apiG := common.NewGin(c)
		// s := getRUri(c)
		// // 如果是白名单的链接，则直接让过(用户不需要登录就让访问的URL)
		// res = web.CheckWhite(c, s)
		// if res {
		// 	c.Next()
		// 	return
		// }
		//
		// // 用户登录信息验证
		// if exitStatus := auth(c); exitStatus {
		// 	return
		// }
		// // 用户登录了的验证权限
		// res = web.CheckPermissions(c, s)
		//
		// // 如果不在白名单范围内，则让过
		// if !res {
		// 	app_log.GetLog().Error(map[string]string{
		// 		"method":      "middleware.Permission",
		// 		"info":        "router permission",
		// 		"router name": c.Request.RequestURI,
		// 		"httpMethod":  c.Request.Method,
		// 	})
		// 	obj := base.NewResult()
		// 	obj.Code = http.StatusForbidden
		// 	obj.Msg = fmt.Sprintf("no auth(%s)",s)
		// 	c.JSON(http.StatusOK, obj)
		// 	c.Abort()
		// 	return
		// }
		//
		// // 获取当前登录用户信息
		// code, rd := userMessageSet(c, c.Request.RequestURI)
		// if code > 0 {
		// 	apiG.Response(code, rd)
		// 	return
		// }
		//
		c.Next()
	}
}

// 用户登录逻辑处理
func Auth(c *gin.Context) (exit bool) {
	token := c.Request.Header.Get("x-auth-token")
	if token == "" {
		msg := "token is null"
		app_log.GetLog().Error(map[string]string{
			"method": "zgh.ginmiddleware.auth",
			"error":  msg,
		})
		c.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
		c.Abort()
		exit = true
		return
	}
	jwtUser, err := common.ParseToken(token)
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"method": "zgh.ginmiddleware.auth",
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

func cors(c *gin.Context) (exitStatus bool) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE,PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token, X-Auth-UUID, X-Auth-Openid, referrer, Authorization, x-client-id, x-client-version, x-client-type")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	return
}

// 用户信息获取
// func UserMessageSet(c *gin.Context, routerAsName string) (code int, res interface{}) {
// 	token := c.GetHeader("x-auth-token")
// 	if routerAsName == "console.post.imgUpload" { // 如果是上传图片，则用的POST获取用户信息
// 		token = c.PostForm("upload-token")
// 	}
//
// 	if token == "" {
// 		app_log.GetLog().Errorln("method", "middleware.Permission", "info", "token null")
// 		return 400001005, nil
// 	}
//
// 	jwtUser, err := common.ParseToken(token)
// 	if err != nil {
// 		app_log.GetLog().Errorln("method", "middleware.Permission", "info", "parse token error")
// 		return 400001005, nil
// 	}
// 	c.Set(app_obj.ContextUserObjectKey, jwtUser)
// 	c.Set("token", token)
// 	return
// }


func GetRUri(c *gin.Context) string {
	uri := strings.TrimLeft(c.Request.RequestURI, common.GetAppConfig().AppName+"/"+common.GetAppConfig().AppApiVersion)
	if uri == "" { // 如果是默认页 ，则直接让过
		return "default"
	}
	s1 := strings.Split(uri, "?")
	s2 := strings.TrimRight(s1[0], "/")
	fmt.Printf("Uri is :'%v'", s2)
	return s2
}
