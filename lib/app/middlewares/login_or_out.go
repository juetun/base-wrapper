package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"net/http"
)

// AuthenticationCallBack 用户身份验证成功
// err的提示内容会在响应中输出
type AuthenticationCallBack func(user *base.JwtUser, c *gin.Context) (err error)

// AuthParse 不用严格判断登录，如果前端传递了令牌那么解析令牌,否则直接跳过
// notStrictValue=true
// token=""
func AuthParse(callBacks ...AuthenticationCallBack) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtUser base.JwtUser
		var err error
		ctx := base.CreateContext(&base.ControllerBase{Log: app_obj.GetLog()}, c)
		jwtUser, _ = common.TokenValidate(ctx, true)

		if callBacks == nil && len(callBacks) == 0 { // 如果没配置回调 则直接结束
			c.Abort()
			return
		}

		for _, callBack := range callBacks {
			// 调用回调方法
			if err = callBack(&jwtUser, c); err != nil {
				c.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// Authentication 判断用户是否登录如果未登录则退出
func Authentication(callBacks ...AuthenticationCallBack) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtUser base.JwtUser
		var exitStatus bool
		var err error
		// 验证token是否合法
		if jwtUser, exitStatus = common.TokenValidate(base.CreateContext(&base.ControllerBase{Log: app_obj.GetLog()}, c), false); exitStatus {
			c.Abort()
			return
		}

		if len(callBacks) == 0 { // 如果没配置回调 则直接结束
			c.Next()
			return
		}

		for _, callBack := range callBacks {
			// 调用回调方法
			if err = callBack(&jwtUser, c); err != nil {
				c.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
				c.Abort()
				return
			}
		}
		c.Next()
		return
	}
}

// func RequestPathPermit(c *gin.Context, s string) (res bool) {
//	res = true
//	// 用户登录了的验证权限
//	res = permissions.CheckPermissions(c, s)
//	// 如果不在白名单范围内，则让过
//	if !res {
//		app_obj.GetLog().Error(c, map[string]interface{}{
//			"method":      "middleware.Permission",
//			"info":        "router permission",
//			"router name": c.Request.RequestURI,
//			"httpMethod":  c.Request.Method,
//		})
//		obj := base.NewResult()
//		obj.Code = http.StatusForbidden
//		obj.Data = ""
//		obj.Msg = fmt.Sprintf("您没有权限访问本功能(no auth:%s %s)", s, c.Request.Method)
//		c.JSON(http.StatusOK, obj)
//		c.Abort()
//		res = false
//		return
//	}
//	return
// }
