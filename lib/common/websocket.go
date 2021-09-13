// Package common
// @Deprecated 即将弃用方法
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
// websocket gin集成方法
package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler func(conn *websocket.Conn, arg interface{})

// GinWebsocketHandler websocket.Handler 转 gin HandlerFunc
func GinWebsocketHandler(wsConnHandle WebsocketHandler, argObject interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("new websocket request: %v", c.Request.Header)
		if c.IsWebsocket() {
			var err error
			if err = c.ShouldBind(argObject); err != nil {
				log.Printf("new websocket request err : %v", err.Error())
				return
			}
			var conn *websocket.Conn
			if conn, err = upGrader.Upgrade(c.Writer, c.Request, nil); err != nil {
				log.Printf("new websocket request err : %v", err.Error())
				return
			}
			wsConnHandle(conn, argObject)
			// wsConnHandle.ServeHTTP(c.Writer, &req)
		} else {
			_, _ = c.Writer.WriteString("===not websocket request===")
		}
	}
}

//
// // WsConnHandle websocket连接处理
// func WsConnHandle(conn *websocket.Conn) {
// 	for {
// 		var msg string
// 		if err := websocket.Message.Receive(conn, &msg); err != nil {
// 			log.Println(err)
// 			return
// 		}
//
// 		log.Printf("recv: %v", msg)
//
// 		data := []byte(time.Now().Format(time.RFC3339))
// 		if _, err := conn.Write(data); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }
