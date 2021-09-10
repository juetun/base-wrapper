package srvs

import (
	"golang.org/x/net/websocket"
)

type SrvWebSocket interface {

	WebsocketSrv(conn *websocket.Conn)

}
