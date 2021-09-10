package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil"
	"github.com/juetun/base-wrapper/web/daos"
)

type DaoWebSocketImpl struct {
	base.ServiceDao
}

func (r *DaoWebSocketImpl) SyncMessageByUserIds(userIds ...string) (err error) {
	panic("implement me")
}

func (r *DaoWebSocketImpl) GetUnReadMessageCount(userHIds string) (total int64, err error) {
	panic("implement me")
}

func (r *DaoWebSocketImpl) CreateMessage(arg *websocket_anvil.PushMessageRequestStruct) (err error) {
	panic("implement me")
}

func (r *DaoWebSocketImpl) BatchUpdateMessageRead(ids []string) error {
	panic("implement me")
}

func (r *DaoWebSocketImpl) UpdateAllMessageDeleted(hid string) error {
	panic("implement me")
}

func (r *DaoWebSocketImpl) UpdateAllMessageRead(hid string) error {
	panic("implement me")
}

func (r *DaoWebSocketImpl) BatchUpdateMessageDeleted(ids []string) error {
	panic("implement me")
}

func NewDaoWebSocketImpl(ctx ...*base.Context) daos.DaoWebSocket {
	p := &DaoWebSocketImpl{}
	p.SetContext(ctx...)

	return p
}
