package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

type RequestIDOptions struct {
	AllowSetting bool
}

// RequestID 获取请求ID 中间件
func RequestID(options RequestIDOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string
		beginTime := strconv.FormatInt(time.Now().UnixNano(), 10)

		if options.AllowSetting {
			// If Set-Request-Id header is set on request, use that for
			// Request-Id response header. Otherwise, generate a new one.
			requestID = c.Request.Header.Get("Set-Request-Id")
		}

		if requestID == "" {
			s := uuid.NewV4()
			// if err != nil {
			//	zgh.ZLog().Error("message","uuid create  error","error",err.Error())
			// }
			requestID = s.String()
		}

		c.Writer.Header().Set("X-Begin-Time", beginTime)
		c.Writer.Header().Set("X-Request-Id", requestID)
		// zgh.ZLog().Info("Message", "API Request", "header", c.Request.Header)
		c.Next()
	}
}
