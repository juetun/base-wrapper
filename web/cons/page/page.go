// Package page
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page

import (
	"github.com/gin-gonic/gin"
)

type ConPage interface {

	// Websocket web socket操作
	Websocket(c *gin.Context)

	Tsst(c *gin.Context)

	Main(c *gin.Context)

	MainSign(c *gin.Context)
}
