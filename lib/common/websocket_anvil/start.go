package websocket_anvil

import (
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil/ext"
)

func WebsocketStart() {
	ext.HubMessage = ext.MessageHub{

		// 用户最后活跃时间
		UserLastActive: make(map[string]int64, ext.ClientConnectMax),
		UserIds:        make([]string, 0, ext.ClientConnectMax),
		// 客户端集合(用户id为每个socket key)
		Clients:   make(map[string]*ext.MessageClient, ext.ClientConnectMax),
		Broadcast: ext.NewCh(),
		// 刷新用户消息通道
		RefreshUserMessage: ext.NewCh(),
	}
	ext.HubMessage.Count()
	ext.HubMessage.Run()
}
