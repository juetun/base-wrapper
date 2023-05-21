package common

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"net/http"
	"strconv"
)

// 用户登录逻辑处理
// param  notStrictValue    	true:当token=""时跳过
// return bool 					true:用户信息获取失败，false:正常操作
func TokenValidate(c *base.Context, notStrictValue bool) (jwtUser base.JwtUser, exit bool) {
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
		c.GinContext.JSON(http.StatusOK, NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
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
		c.GinContext.JSON(http.StatusOK, NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
		exit = true
		return
	}
	if jwtUser.UserId != userHid {
		err = fmt.Errorf("用户信息(token uid)不匹配")
		c.GinContext.JSON(http.StatusOK, NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
		exit = true
		return
	}
	// 解析token成功 将用户信息放进gin 上下文对象context中
	c.GinContext.Set(base.ContextUserObjectKey, jwtUser)
	c.GinContext.Set(base.ContextUserTokenKey, token)
	return
}
