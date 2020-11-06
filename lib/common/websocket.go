// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
//websocket gin集成方法
package common

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// websocket.Handler 转 gin HandlerFunc
func GinWebsocketHandler(wsConnHandle websocket.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("new ws request: %v", c.Request.RemoteAddr)
		if c.IsWebsocket() {
			wsConnHandle.ServeHTTP(c.Writer, c.Request)
		} else {
			_, _ = c.Writer.WriteString("===not websocket request===")
		}
	}
}
// websocket连接处理
func WsConnHandle(conn *websocket.Conn) {
	for {
		var msg string
		if err := websocket.Message.Receive(conn, &msg); err != nil {
			log.Println(err)
			return
		}

		log.Printf("recv: %v", msg)

		data := []byte(time.Now().Format(time.RFC3339))
		if _, err := conn.Write(data); err != nil {
			log.Println(err)
			return
		}
	}
}

