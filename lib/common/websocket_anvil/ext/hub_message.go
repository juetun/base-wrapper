package ext

import (
	"fmt"
	"sync"
	"time"
)

// MessageHub 消息仓库, 用于维护整个消息中心连接
type MessageHub struct {
	WebsocketBaseHandler

	lock sync.RWMutex

	// redis连接
	Service MessageService

	// 客户端用户id集合
	UserIds []string

	// 用户最后活跃时间
	UserLastActive map[string]int64

	// 客户端集合(用户id为每个socket key)
	Clients map[string]*MessageClient

	// 广播通道
	Broadcast *Chan

	// 刷新用户消息通道
	RefreshUserMessage *Chan

	// 幂等性token校验方法
	CheckIdempotenceTokenFunc func(token string) bool
}

// Contains 判断uint数组是否包含item元素
func (r *MessageHub) Contains(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

// 广播(全部用户均可接收)
func (r *MessageHub) broadcast(broadcast *MessageBroadcast) {

	for _, client := range r.GetClients() {
		userHid, _, _ := client.GetUserHid()
		// 通知指定用户
		if r.Contains(broadcast.UserIds, userHid) {
			client.SendChan.SafeSend(broadcast)
		}
	}
	return
}

func (r *MessageHub) refreshMsg(userIds ...string) (err error) {

	// 同步用户消息
	if err = HubMessage.Service.Dao.SyncMessageByUserIds(userIds...); err != nil {
		panic(err)
	}
	for _, client := range r.GetClients() {
		for _, id := range userIds {
			userHid, _, _ := client.GetUserHid()
			if userHid == id {
				// 获取未读消息条数
				total, _ := HubMessage.Service.Dao.GetUnReadMessageCount(userHid)
				// 将当前消息条数发送给用户
				msg := MessageWsResponseStruct{
					Type: MessageRespUnRead,
					Detail: r.GetSuccessWithData(map[string]int64{
						"unReadCount": total,
					}),
				}
				client.SendChan.SafeSend(msg)
			}
		}
	}
	return
}

// Run 运行仓库
func (r *MessageHub) Run() {
	for {
		select {
		case data := <-r.Broadcast.C: // 广播(全部用户均可接收)
			broadcast := data.(MessageBroadcast)
			r.broadcast(&broadcast)
		case data := <-r.RefreshUserMessage.C: // 刷新客户端消息
			userIds := data.([]string)
			_ = r.refreshMsg(userIds...)
		}
	}
}

// Count 活跃连接检查
func (r *MessageHub) Count() {
	// 创建定时器, 超出指定时间间隔
	ticker := time.NewTicker(HeartBeatPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		// 到心跳检测时间
		case <-ticker.C:
			infos := make([]string, 0)
			for _, client := range r.GetClients() {
				uid, _, _ := client.GetUserHid()
				infos = append(infos, fmt.Sprintf("%s-%s", uid, client.Ip))
			}
			r.Service.Context.Debug(map[string]interface{}{
				"infos": infos,
				"desc":  fmt.Sprintf("[消息中心]当前活跃连接"),
			}, "MessageHubCount")
		}
	}
}

// GetClients 获取client列表
func (r *MessageHub) GetClients() map[string]*MessageClient {
	HubMessage.lock.RLock()
	defer HubMessage.lock.RUnlock()
	return HubMessage.Clients
}
