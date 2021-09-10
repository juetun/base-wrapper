package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil"
	"github.com/juetun/base-wrapper/web/models"
	"github.com/juetun/base-wrapper/web/srvs"
	"golang.org/x/net/websocket"
)

type SrvWebSocketImpl struct {
	base.ServiceBase
}

func (r *SrvWebSocketImpl) WebsocketSrv(conn *websocket.Conn) {

	websocket_anvil.NewWebSocketAnvil(
		websocket_anvil.WebSocketAnvilOptionContext(r.Context),
		websocket_anvil.WebSocketAnvilOptionUser(func() (user websocket_anvil.UserInterface, err error) {
			if userId := conn.Request().Header.Get(app_obj.HttpUserHid); userId == "" {
				return
			} else {
				user, err = r.getCurrentUserByUid(userId)
				return
			}
		}),
		websocket_anvil.WebSocketAnvilOptionConn(conn),
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

	return
}

func NewSrvWebSocketImpl(ctx ...*base.Context) srvs.SrvWebSocket {

	p := &SrvWebSocketImpl{}
	p.SetContext(ctx...)

	return p
}
