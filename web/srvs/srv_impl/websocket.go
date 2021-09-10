package srv_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil"
	"github.com/juetun/base-wrapper/web/daos/dao_impl"
	"github.com/juetun/base-wrapper/web/srvs"
	"golang.org/x/net/websocket"
)

type SrvWebSocketImpl struct {
	base.ServiceBase
}

func (r *SrvWebSocketImpl) WebsocketSrv(conn *websocket.Conn) {
	daoMessage := dao_impl.NewDaoWebSocketImpl(r.Context)
	websocket_anvil.NewMessageService(
		websocket_anvil.MessageServiceContext(r.Context),
		websocket_anvil.MessageServiceMysql(daoMessage),
	)
}

func NewSrvWebSocketImpl(ctx ...*base.Context) srvs.SrvWebSocket {

	p := &SrvWebSocketImpl{}
	p.SetContext(ctx...)

	return p
}
