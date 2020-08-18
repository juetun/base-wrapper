/**
* @Author:changjiang
* @Description:
* @File:log
* @Version: 1.0.0
* @Date 2020/8/18 2:46 下午
 */
package middlewares

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/sirupsen/logrus"
)

// gin框架日志收集
func GinLogCollect() gin.HandlerFunc {
	logger := app_log.GetLog()
	return func(c *gin.Context) {
		// if c.Request.Method != "POST" && c.Request.Method != "GET" || c.Request.URL.String() == "/metrics" {
		// 	c.Next()
		// 	return
		// }
		c.Request.URL.RawQuery, _ = url.QueryUnescape(c.Request.URL.RawQuery)
		path := c.Request.URL
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			c.Set("body", string(bodyBytes))
		}
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		start := time.Now()
		defer func() {
			fields := logrus.Fields{
				"status":   c.Writer.Status(),
				"method":   c.Request.Method,
				"path":     path.String(),
				"ip":       c.ClientIP(),
				"latency":  time.Now().Sub(start).String(),
				"request":  string(bodyBytes),
				"response": blw.body.String(),
				"header":   c.Request.Header,
				// "header_version": c.GetHeader("version"),
				// "header_build":   c.GetHeader("build"),
				// "header_user_id": c.GetHeader("user_id"),
				// "header_userid":  c.GetHeader("userid"),
				// "header_token":   c.GetHeader("token"),
			}
			if len(c.Request.Form) > 0 {
				fields["request"] = c.Request.Form.Encode()
			}
			if len(c.Errors) > 0 {
				logger.Logger.Error(fields, c.Errors.String())
				// Append error field if this is an erroneous request.
			} else {
				logger.Logger.Info(fields, c.Errors.String())
			}
		}()
		c.Next()

	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
