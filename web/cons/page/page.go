// Package page
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type ConPage interface {

	// Websocket web socket操作
	Websocket(c *websocket.Conn)

	Tsst(c *gin.Context)

	Main(c *gin.Context)

	MainSign(c *gin.Context)
}
