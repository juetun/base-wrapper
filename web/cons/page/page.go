/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:04 下午
 */

// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type ConPage interface {
	//web socket操作
	Websocket(conn *websocket.Conn)
	Tsst(c *gin.Context)
	Main(c *gin.Context)
	MainSign(c *gin.Context)
}
