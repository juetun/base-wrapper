package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
)

// ErrorHandler 错误请求处理
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				switch e.(type) { // 系统级错误屏蔽
				case *base.ErrorRuntimeStruct:
					structContent := e.(*base.ErrorRuntimeStruct)
					result := base.NewResult().
						SetErrorMsg(structContent)
					if structContent.Code > 0 {
						result.SetCode(structContent.Code)
					}
					c.AbortWithStatusJSON(http.StatusOK, result)
					return
				case base.ErrorRuntimeStruct:
					structContent := e.(base.ErrorRuntimeStruct)
					result := base.NewResult().
						SetErrorMsg(&structContent)
					if structContent.Code > 0 {
						result.SetCode(structContent.Code)
					}
					c.AbortWithStatusJSON(http.StatusOK, result)
				default:
					err := e.(error)
					result := base.NewResult().
						SetErrorMsg(err)
					result.SetCode(base.ErrorSystem)
					c.AbortWithStatusJSON(http.StatusOK, result)
				}
			}
		}()
		c.Next()
	}
}
