package websocket_anvil

import (
	"github.com/juetun/base-wrapper/lib/base"
)

var HubMessage MessageHub

// NewMessageService 初始化服务
// NewMessageService([]MessageServiceOption{
// MessageServiceContext(),
// MessageServiceMysql(),
// }...)
func NewMessageService(opts ...MessageServiceOption) (s *MessageService) {
	s = &MessageService{}
	for _, opt := range opts {
		opt(s)
	}
	if s.Context == nil {
		panic("必须设置Context的值")
	}
	return
}

// MessageService 所有的查询可以走redis, 但数据的更新还是走mysql
type MessageService struct {
	Context *base.Context  // 上下文
	Dao     MysqlInterface // 保留mysql, 如果没开启redis可以走mysql
}

type MysqlInterface interface {

	// SyncMessageByUserIds 同步用户消息
	SyncMessageByUserIds(userIds ...string) (err error)

	// GetUnReadMessageCount 获取未读消息条数
	GetUnReadMessageCount(userHIds string) (total int64, err error)

	CreateMessage(r *PushMessageRequestStruct) (err error)

	BatchUpdateMessageRead(ids []string) error

	UpdateAllMessageDeleted(hid string) error

	UpdateAllMessageRead(hid string) error

	BatchUpdateMessageDeleted(ids []string) error
}

type MessageServiceOption func(messageService *MessageService)

func MessageServiceContext(Context *base.Context) MessageServiceOption {
	return func(messageService *MessageService) {
		messageService.Context = Context
	}
}

func MessageServiceMysql(Mysql MysqlInterface) MessageServiceOption {
	return func(messageService *MessageService) {
		messageService.Dao = Mysql
	}
}

// // Find 查询, model需使用指针, 否则可能无法绑定数据
// func (r MessageService) Find(query *redis.QueryRedis, page *response.PageInfo, model interface{}) (err error) {
// 	// 获取model值
// 	rv := reflect.ValueOf(model)
// 	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Slice) {
// 		return fmt.Errorf("model必须是非空指针数组类型")
// 	}
//
// 	if !page.NoPagination {
// 		// 查询条数
// 		err = query.Count(&page.Total).Error
// 		if err == nil && page.Total > 0 {
// 			// 获取分页参数
// 			limit, offset := page.GetLimit()
// 			err = query.Limit(limit).Offset(offset).Find(model).Error
// 		}
// 	} else {
// 		// 不使用分页
// 		err = query.Find(model).Error
// 		if err == nil {
// 			page.Total = int64(rv.Elem().Len())
// 			// 获取分页参数
// 			page.GetLimit()
// 		}
// 	}
// 	return
// }
