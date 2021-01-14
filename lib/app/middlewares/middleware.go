package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
)

var MiddleWareComponent = []gin.HandlerFunc{
	CrossOriginResourceSharing(), // 配置跨域逻辑
}
var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

// 加载权限验证Gin中间件
func LoadMiddleWare(privateMiddleWares ...gin.HandlerFunc) {

	io.SystemOutPrintln("Load GIN middleWare start")

	MiddleWareComponent = append(MiddleWareComponent, privateMiddleWares...)

	io.SystemOutPrintln("Load GIN middleWare finished")
}
