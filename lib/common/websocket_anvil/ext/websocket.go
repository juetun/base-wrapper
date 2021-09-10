package ext

import (
	"fmt"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"golang.org/x/net/websocket"
)

// WebSocketAnvil websocket操作基本结构体
type WebSocketAnvil struct {
	Context  *base.Context   `json:"-"`
	Conn     *websocket.Conn `json:"-"`
	UserFunc UserHandler     `json:"-"`

	Ip  string `json:"ip"`
	Key string `json:"key"`
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

func WebSocketAnvilOptionIp(ip string) WebSocketAnvilOption {
	return func(arg *WebSocketAnvil) {
		arg.Ip = ip
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
	if err = r.initWebSocketKey(); err != nil {
		return
	}
	logContent["key"] = r.Key

	if err = r.initClientIp(); err != nil {
		return
	}

	logContent["ip"] = r.Ip

	// 注册到消息仓库
	client := &MessageClient{
		Context:  r.Context,
		Key:      r.Key,
		Conn:     r.Conn,
		UserFunc: r.UserFunc,
		Ip:       r.Ip,
		SendChan: NewCh(),
	}

	go client.Register()

	// 监听数据的接收/发送/心跳
	go client.Receive()

	go client.Send()
	// go client.heartBeat()

	return
}

func (r *WebSocketAnvil) initWebSocketKey() (err error) {

	if key, ok := r.Conn.Request().Header[app_obj.WebSocketKey]; ok {
		r.Key = strings.Join(key, "")
	}
	if r.Key == "" {
		err = fmt.Errorf("没获取到(%s)值", app_obj.WebSocketKey)
	}
	return
}

func (r *WebSocketAnvil) initClientIp() (err error) {
	header := r.Conn.Request().Header
	if key, ok := header[app_obj.WebSocketHeaderIp]; ok {
		r.Ip = strings.Join(key, "")
	}
	if r.Ip == "" {
		r.Ip = r.Conn.Request().RemoteAddr
	}

	return
}
