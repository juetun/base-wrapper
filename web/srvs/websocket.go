package srvs

import (
	"github.com/gorilla/websocket"
	"github.com/juetun/base-wrapper/web/wrapper"
)

type SrvWebSocket interface {
	WebsocketSrv(conn *websocket.Conn, arg *wrapper.ArgWebSocket)
}
