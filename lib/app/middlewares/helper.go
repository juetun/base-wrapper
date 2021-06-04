// Package middlewares
package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MyCors 设置跨域让访问的逻辑
func MyCors(c *gin.Context) {

	method := c.Request.Method

	// origin := c.Request.Header.Get("Origin")
	origin := c.Request.Header.Get("Origin")

	// 如果是本地访问,则URL地址为 http://localhost:8081
	if origin == "http://localhost:8081" {
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", headerStr)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Set("content-type", "application/json")
	}

	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "Options Request!")
	}

	c.Next()

}
