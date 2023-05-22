package base

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"strings"
)

//判断是否为内网访问
func InterPath(c *gin.Context) (ok bool) {
	paths := strings.Split(c.Request.URL.Path, "/")
	if len(paths) > 2 && paths[2] == app_obj.RouteTypeDefaultIntranet {
		ok = true
		return
	}
	return
}
