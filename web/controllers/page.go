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
package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type ControllerPage interface {
	//web socket操作
	Websocket(conn *websocket.Conn)
	Tsst(c *gin.Context)
	Main(c *gin.Context)
}
