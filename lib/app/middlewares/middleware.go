package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
)

var MiddleWareComponent = []gin.HandlerFunc{

	CrossOriginResourceSharing(), // 配置跨域逻辑
	GinLogCollect(),              // 日志操作逻辑
}
var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

// LoadMiddleWare 加载权限验证Gin中间件
func LoadMiddleWare(privateMiddleWares ...gin.HandlerFunc) {

	io.SystemOutPrintln("Load GIN middleWare start")
	switch common.GetEnv() {
	case gin.ReleaseMode: // 线上环境增加异常捕获
		// 捕捉程序异常操作
		MiddleWareComponent = append([]gin.HandlerFunc{ErrorHandler()}, MiddleWareComponent...)
	default:
	}

	MiddleWareComponent = append(MiddleWareComponent, privateMiddleWares...)

	io.SystemOutPrintln("Load GIN middleWare finished")
}
