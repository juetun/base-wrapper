package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/common"
)

var MiddleWareComponent = []gin.HandlerFunc{
	Permission(),
}
var io = common.NewSystemOut().SetInfoType(common.LogLevelInfo)

// 加载权限验证Gin中间件
func LoadMiddleWare(privateMiddleWares ...gin.HandlerFunc) {

	io.SystemOutPrintln("Load GIN middleWare start")

	MiddleWareComponent = append(MiddleWareComponent, privateMiddleWares...)

	io.SystemOutPrintln("Load GIN middleWare finished")
}
