package middlewares

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
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
		jwtUser, _ = tokenValidate(ctx, true)

		if callBacks == nil && len(callBacks) == 0 { // 如果没配置回调 则直接结束
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
	}
}

// Authentication 判断用户是否登录如果未登录则退出
func Authentication(callBacks ...AuthenticationCallBack) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtUser base.JwtUser
		var exitStatus bool
		var err error
		// 验证token是否合法
		if jwtUser, exitStatus = tokenValidate(base.CreateContext(&base.ControllerBase{Log: app_obj.GetLog()}, c), false); exitStatus {
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

// 用户登录逻辑处理
// param  notStrictValue    	true:当token=""时跳过
// return bool 					true:用户信息获取失败，false:正常操作
func tokenValidate(c *base.Context, notStrictValue bool) (jwtUser base.JwtUser, exit bool) {
	jwtUser = base.JwtUser{}
	var (
		token      string
		userHid    int64
		err        error
		logContent = make(map[string]interface{}, 10)
	)
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			c.Error(logContent, "baseWrapperTokenValidate")
		}
	}()
	c.GinContext.Set(app_obj.TraceId, c.GinContext.GetHeader(app_obj.HttpTraceId))
	if token = c.GinContext.Request.Header.Get(app_obj.HttpUserToken); token == "" { // 如果token为空

		// 如果token为空且设置了空跳过，则直接退出
		if notStrictValue == true {
			return
		}
		logContent["desc"] = "getHeader"
		msg := "token is null"
		err = fmt.Errorf(msg)
		c.GinContext.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
		exit = true
		return
	}

	userHidString := c.GinContext.Request.Header.Get(app_obj.HttpUserHid)
	if userHid, err = strconv.ParseInt(userHidString, 10, 64); err != nil {
		logContent["desc"] = "ParseInt"
		return
	}
	if err = base.ParseJwtKey(token, c, &jwtUser); err != nil { // 如果解析token失败
		logContent["desc"] = "ParseToken"
		logContent["token"] = token
		c.GinContext.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
		exit = true
		return
	}
	if jwtUser.UserId != userHid {
		err = fmt.Errorf("用户信息(token uid)不匹配")
		c.GinContext.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
		exit = true
		return
	}
	// 解析token成功 将用户信息放进gin 上下文对象context中
	c.GinContext.Set(base.ContextUserObjectKey, jwtUser)
	c.GinContext.Set(base.ContextUserTokenKey, token)
	return
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
