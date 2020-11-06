// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
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
