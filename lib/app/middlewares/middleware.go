package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
)

var MiddleWareComponent = []gin.HandlerFunc{
	ErrorHandler(),                  // 捕捉程序异常操作
	CrossOriginResourceSharing(),    // 配置跨域逻辑
	GinLogCollect(), // 日志操作逻辑
}
var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

// LoadMiddleWare 加载权限验证Gin中间件
func LoadMiddleWare(privateMiddleWares ...gin.HandlerFunc) {

	io.SystemOutPrintln("Load GIN middleWare start")

	MiddleWareComponent = append(MiddleWareComponent, privateMiddleWares...)

	io.SystemOutPrintln("Load GIN middleWare finished")
}
