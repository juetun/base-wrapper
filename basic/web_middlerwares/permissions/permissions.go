/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-13
 * Time: 22:39
 */
package permissions

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type HttpPermit struct {
	Method []string `json:"method"`
	Uri    string   `json:"uri"`
}

var PermissionsWhite = []HttpPermit{
	{
		Method: []string{"GET", "POST"},
		Uri:    "console/login",
	},
	{
		Method: []string{"GET", "POST"},
		Uri:    "console/register",
	},
}

// 用户登录后 具备的接口访问权限
var Permissions = []HttpPermit{
	{
		Method: []string{"GET"},
		Uri:    "console/login",
	},
	{
		Method: []string{"GET"},
		Uri:    "console/home",
	},
	{
		Method: []string{"GET"},
		Uri:    "console/post/trash",
	},
	{
		Method: []string{"GET",},
		Uri:    "console/post",
	},
	{
		Method: []string{"DELETE"},
		Uri:    `^console\/post\/[^\/]+$`,
	},
	{
		Method: []string{"GET"},
		Uri:    `^console\/post\/edit\/[^\/]+$`,
	},
	{
		Method: []string{"GET"},
		Uri:    `console/cate`,
	}, {
		Method: []string{"GET"},
		Uri:    `console/system`,
	}, {
		Method: []string{"GET"},
		Uri:    `console/tag`,
	},
	{
		Method: []string{"DELETE"},
		Uri:    `console/cache`,
	},
	{
		Method: []string{"DELETE"},
		Uri:    `console/cache`,
	},
	{
		Method: []string{"DELETE"},
		Uri:    `console/logout`,
	},
	{
		Method: []string{"DELETE"},
		Uri:    `^console\/post\/[^\/]+\/trash$`,
	},
	{
		Method: []string{"GET"},
		Uri:    `console/link`,
	},
	{
		Method: []string{"GET"},
		Uri:    `/base-wrapper/page/test`,
	},
}

// 需要验证权限的配置列表
// 不需要验证权限的配置列表
func CheckPermissions(c *gin.Context, s string) (res bool) {
	app_obj.GetLog().Info(c, map[string]interface{}{
		"request_Uri": s,
		"router name": c.Request.RequestURI,
		"httpMethod":  c.Request.Method,
	})
	for _, v := range Permissions {
		if res = everyValidateTrueOrFalse(&v.Method, c.Request.Method, v.Uri, s); res {
			return
		}
	}
	return false
}

func everyValidateTrueOrFalse(methodArea *[]string, method, uri, s string) bool {
	var validateMethod bool
	if s == "default" { // 默认 default路径直接让过
		return true
	}
	validateMethod = false

	if len(*methodArea) != 0 {

		// 如果请求方法是返回内的值
		for _, value := range *methodArea {
			if value == method {
				validateMethod = true
			}
		}

		// 如果请求方法是返回内的值 并且请求地址对，则认为对
		if validateMethod == true {
			if uri == s {
				return true
			}
			// 写的正则表达式验证通过
			isMatch, _ := regexp.MatchString(uri, s)
			if isMatch {
				return true
			}
		}
	}

	// 否则 如果权限控制没有设置Method的值 就是表示所有的请求方式都有效，此时只验证URi路径是否正确
	if uri == s {
		return true
	}

	// 写的正则表达式验证通过
	isMatch, _ := regexp.MatchString(uri, s)
	if isMatch {
		return true
	}
	return false
}

// 白名单验证。此部分的接口用户不需要登录即可访问
func CheckWhite(c *gin.Context, s string) (res bool) {

	app_obj.GetLog().Info(c, map[string]interface{}{
		"request_Uri": s,
		"info":        "web.permissions.go(CheckWhite)",
		"router name": c.Request.RequestURI,
		"httpMethod":  c.Request.Method,
	})
	for _, v := range PermissionsWhite {
		if res = everyValidateTrueOrFalse(&v.Method, c.Request.Method, v.Uri, s); res {
			return
		}
	}
	return false
}
