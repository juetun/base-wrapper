package ext

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/juetun/base-wrapper/lib/base"
)

// WebSocketAnvil websocket操作基本结构体
type WebSocketAnvil struct {
	Context       *base.Context          `json:"-"`
	Conn          *websocket.Conn        `json:"-"`
	UserFunc      UserHandler            `json:"-"`
	MessageAccept MessageHandler         `json:"-"`
	CommonParams  *base.ArgWebSocketBase `json:"common_params"`
}

func NewWebSocketAnvil(options ...WebSocketAnvilOption) (res *WebSocketAnvil) {

	p := &WebSocketAnvil{}

	for _, option := range options {
		option(p)
	}
	return p
}

func WebSocketAnvilOptionUser(user UserHandler) WebSocketAnvilOption {
	return func(arg *WebSocketAnvil) {
		arg.UserFunc = user
	}
}

func WebSocketAnvilOptionCommonParams(commonParams *base.ArgWebSocketBase) WebSocketAnvilOption {
	return func(arg *WebSocketAnvil) {
		arg.CommonParams = commonParams
	}
}

func WebSocketAnvilOptionMessageHandler(messageAccept MessageHandler) WebSocketAnvilOption {
	return func(arg *WebSocketAnvil) {
		arg.MessageAccept = messageAccept
	}
}
func WebSocketAnvilOptionConn(conn *websocket.Conn) WebSocketAnvilOption {
	return func(arg *WebSocketAnvil) {
		arg.Conn = conn
	}
}

func WebSocketAnvilOptionContext(context *base.Context) WebSocketAnvilOption {
	return func(arg *WebSocketAnvil) {
		arg.Context = context
	}
}

// Start 启动消息连接
func (r *WebSocketAnvil) Start() (err error) {
	logContent := map[string]interface{}{}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "WebSocketAnvilStart")
		} else {
			r.Context.Info(logContent, "WebSocketAnvilStart")
		}
	}()
	if r.CommonParams == nil {
		err = fmt.Errorf("common_params must be not null")
		return
	}
	logContent["key"] = r.CommonParams.WebsocketKey

	logContent["ip"] = r.CommonParams.Ip

	// 注册到消息仓库
	client := &MessageClient{
		MessageAction: r.MessageAccept,
		Key:           r.CommonParams.WebsocketKey,
		Conn:          r.Conn,
		UserFunc:      r.UserFunc,
		Ip:            r.CommonParams.Ip,
		SendChan:      NewCh(),
	}
	client.Context = r.Context
	client.RequestId = r.CommonParams.WebsocketKey

	go client.Register() // 注册

	go client.Receive() // 监听数据的接收/发送/心跳

	// 监听发送消息
	go client.Send()

	// go client.heartBeat()

	return
}
