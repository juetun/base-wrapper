// Package middlewares
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
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/sirupsen/logrus"
)

var (
	GinLogHeaderNotCollect = []string{
		"Accept",
		"User-Agent",
		"Cache-Control",
		"Sec-Ch-Ua-Mobile",
		"Sec-Ch-Ua",
		"Pragma",
		"Accept-Encoding",
		"Connection",
		"Accept-Language",
	} // GIN框架不用收集的header字段
)

// GinLogCollect gin框架日志收集
func GinLogCollect() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set(app_obj.TraceId, c.GetHeader(app_obj.HttpTraceId)) // 日志对象获取,最先执行的中间件
		start := time.Now()
		defer func() { // 异步操作写日志
			logger := app_obj.GetLog()
			go delayExecGinLogCollect(start, c, c.Request.URL, logger)
		}()
		c.Next()
	}
}
func inExceptHeaderSlice(key string) (res bool) {
	for _, s := range GinLogHeaderNotCollect {
		if s == key {
			res = true
			return
		}
	}
	return
}
func getUseHeader(header *http.Header) (res http.Header) {
	res = http.Header{}
	for s := range *header {
		if inExceptHeaderSlice(s) {
			continue
		}
		res.Set(s, header.Get(s))
	}
	return
}

// 流量日志收集
func delayExecGinLogCollect(start time.Time, c *gin.Context, path *url.URL, logger *app_obj.AppLog) {
	c.Request.URL.RawQuery, _ = url.QueryUnescape(c.Request.URL.RawQuery)
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		// c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// c.Set("body", string(bodyBytes))
	}
	fields := logrus.Fields{
		app_obj.AppFieldKey: "GIN",
		"status":            c.Writer.Status(),
		"method":            c.Request.Method,
		"path":              path.String(),
		"ip":                c.ClientIP(),
		"duration":          float64(time.Now().Sub(start) / 1e3), // 时长单位微秒
		"request":           string(bodyBytes),
		"header":            getUseHeader(&c.Request.Header),
	}
	// 只收集 http code>400的错误日志
	if c.Writer.Status() >= 400 {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		fields["response"] = blw.body.String()
	}

	if len(c.Request.Form) > 0 {
		fields["request"] = c.Request.Form.Encode()
	}
	if len(c.Errors) > 0 {
		logger.Error(c, fields, c.Errors.String())
		return
	}
	logger.Info(c, fields, c.Errors.String())
	return
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
