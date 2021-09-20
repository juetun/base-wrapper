// Package websocket_anvil
// @Deprecated 即将弃用方法
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
// websocket gin集成方法
package anvil_websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler func(conn *websocket.Conn, arg interface{})
