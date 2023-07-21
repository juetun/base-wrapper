// Package middlewares
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package middlewares

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"net/http"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
)

// GetRUri 获取重写URL参数
func GetRUri(c *base.Context) string {
	uri := strings.TrimLeft(c.GinContext.Request.RequestURI, common.GetAppConfig().AppName+"/"+common.GetAppConfig().AppApiVersion)
	if uri == "" { // 如果是默认页 ，则直接让过
		return "default"
	}
	s1 := strings.Split(uri, "?")
	s2 := strings.TrimRight(s1[0], "/")
	// fmt.Printf("Uri is :'%v'", s2)
	return s2
}

// Auth 用户登录逻辑处理
func Auth(c *base.Context) (exit bool) {

	token := c.GinContext.Request.Header.Get(app_obj.HttpUserToken)
	var (
		err     error
		jwtUser base.JwtUser
	)
	logContent := make(map[string]interface{})
	defer func() {
		if err != nil {
			c.Error(logContent, "baseWrapperMiddleWaresAuth")
			return
		}

	}()

	if token == "" {
		logContent["desc"] = "zgh.ginmiddleware.Auth"
		msg := "token is null"
		err = fmt.Errorf(msg)
		c.GinContext.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusUnauthorized).SetMessage(msg))
		c.GinContext.Abort()
		exit = true
		return
	}

	if err = base.ParseJwtKey(token, c, &jwtUser); err != nil {
		logContent["desc"] = "ParseToken"
		logContent["token"] = token

		c.GinContext.JSON(http.StatusOK, common.NewHttpResult().SetCode(http.StatusForbidden).SetMessage(err.Error()))
		c.GinContext.Abort()
		exit = true
		return
	}
	c.GinContext.Set(base.ContextUserObjectKey, jwtUser)
	c.GinContext.Set(base.ContextUserTokenKey, token)
	return
}
