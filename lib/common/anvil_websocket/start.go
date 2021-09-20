package anvil_websocket

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/anvil_websocket/ext"
	"github.com/juetun/base-wrapper/web/daos/dao_impl"
)

func WebsocketStart() {
	context := base.NewContext()
	daoMessage := dao_impl.NewDaoWebSocketImpl(context)
	ext.HubMessage = ext.MessageHub{
		// 用户最后活跃时间
		UserLastActive: make(map[string]int64, ext.ClientConnectMax),
		UserIds:        make([]string, 0, ext.ClientConnectMax),
		// 客户端集合(用户id为每个socket key)
		Clients:   make(map[string]*ext.MessageClient, ext.ClientConnectMax),
		Broadcast: ext.NewCh(),
		// 刷新用户消息通道
		RefreshUserMessage: ext.NewCh(),
		Service: ext.NewMessageService(
			ext.MessageServiceContext(context),
			ext.MessageServiceMysql(daoMessage),
		),
	}
	ext.HubMessage.Count().Run()
}
