package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil/ext"
	"github.com/juetun/base-wrapper/web/models"
	"github.com/juetun/base-wrapper/web/srvs"
	"golang.org/x/net/websocket"
)

type SrvWebSocketImpl struct {
	base.ServiceBase
}

func (r *SrvWebSocketImpl) WebsocketSrv(conn *websocket.Conn) {

	ext.NewWebSocketAnvil(
		ext.WebSocketAnvilOptionContext(r.Context),
		ext.WebSocketAnvilOptionUser(func() (user ext.UserInterface, err error) {
			conn.Request().ParseForm()
			if userId := conn.Request().FormValue(app_obj.WebSocketUid); userId == "" {
				return
			} else {
				user, err = r.getCurrentUserByUid(userId)
				return
			}
		}),
		ext.WebSocketAnvilOptionConn(conn),
	).Start()

	// daoMessage := dao_impl.NewDaoWebSocketImpl(r.Context)
	// operate := websocket_anvil.NewMessageService(
	// 	websocket_anvil.MessageServiceContext(r.Context),
	// 	websocket_anvil.MessageServiceMysql(daoMessage),
	// )
	// operate.
}

// GetCurrentUser 获取当前请求用户信息
func (r *SrvWebSocketImpl) getCurrentUserByUid(userHid string) (res *models.User, err error) {
	res = &models.User{UserHid: userHid}
	return
}

func NewSrvWebSocketImpl(ctx ...*base.Context) srvs.SrvWebSocket {

	p := &SrvWebSocketImpl{}
	p.SetContext(ctx...)

	return p
}
