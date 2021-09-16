package srv_impl

import (
	"github.com/gorilla/websocket"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil/ext"
	"github.com/juetun/base-wrapper/web/models"
	"github.com/juetun/base-wrapper/web/srvs"
	"github.com/juetun/base-wrapper/web/wrapper"
)

type SrvWebSocketImpl struct {
	base.ServiceBase
}

func (r *SrvWebSocketImpl) WebsocketSrv(conn *websocket.Conn, arg *wrapper.ArgWebSocket) {

	ext.NewWebSocketAnvil(
		ext.WebSocketAnvilOptionCommonParams(&arg.ArgWebSocketBase),
		ext.WebSocketAnvilOptionContext(r.Context),
		ext.WebSocketAnvilOptionUser(func() (user ext.UserInterface, err error) {
			if arg.UserHid == "" {
				return
			}
			user, err = r.getCurrentUserByUid(arg.UserHid)
			return
		}),
		ext.WebSocketAnvilOptionConn(conn),
		ext.WebSocketAnvilOptionMessageHandler(r.messageLogicHandler),
	).Start()

}

// GetCurrentUser 获取当前请求用户信息
func (r *SrvWebSocketImpl) getCurrentUserByUid(userHid string) (res *models.User, err error) {
	res = &models.User{UserHid: userHid}
	return
}


// 消息接收处理
func (r *SrvWebSocketImpl) messageLogicHandler(data interface{}) (res interface{}, err error) {

	return
}

func NewSrvWebSocketImpl(ctx ...*base.Context) srvs.SrvWebSocket {

	p := &SrvWebSocketImpl{}
	p.SetContext(ctx...)

	return p
}
