// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common"
)

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

// 用户登录逻辑处理
func Auth(c *gin.Context) (exit bool) {

	token := c.Request.Header.Get(app_obj.HttpUserToken)


	if token == "" {
		msg := "token is null"
		app_obj.GetLog().Error(c, map[string]interface{}{
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
		app_obj.GetLog().Error(c, map[string]interface{}{
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
