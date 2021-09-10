package websocket_anvil

import (
	"github.com/juetun/base-wrapper/lib/base"
	"golang.org/x/net/websocket"
)

// WebSocketAnvil websocket操作基本结构体
type WebSocketAnvil struct {
}

// MessageWebSocket 启动消息连接
func (s WebSocketAnvil) MessageWebSocket(ctx *base.Context, conn *websocket.Conn, key string, user UserInterface, ip string) {

	// 注册到消息仓库
	client := &MessageClient{
		Context:  ctx,
		Key:      key,
		Conn:     conn,
		User:     user,
		Ip:       ip,
		SendChan: NewCh(),
	}

	go client.Register()

	// 监听数据的接收/发送/心跳
	go client.Receive()

	go client.Send()
	// go client.heartBeat()
}

type UserInterface interface {
	// GetUserHid 获取用户信息
	GetUserHid() (userHid string, err error)
}
