package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/anvil_websocket/ext"
	"github.com/juetun/base-wrapper/web/daos"
)

type DaoWebSocketImpl struct {
	base.ServiceDao
}

func (r *DaoWebSocketImpl) SyncMessageByUserIds(userIds ...string) (err error) {
	panic("implement me SyncMessageByUserIds")
}

func (r *DaoWebSocketImpl) GetUnReadMessageCount(userHIds string) (total int64, err error) {
	panic("implement me GetUnReadMessageCount")
}

func (r *DaoWebSocketImpl) CreateMessage(arg *ext.PushMessageRequestStruct) (err error) {
	panic("implement me CreateMessage")
}

func (r *DaoWebSocketImpl) BatchUpdateMessageRead(ids []string) error {
	panic("implement me BatchUpdateMessageRead")
}

func (r *DaoWebSocketImpl) UpdateAllMessageDeleted(hid string) error {
	panic("implement me UpdateAllMessageDeleted")
}

func (r *DaoWebSocketImpl) UpdateAllMessageRead(hid string) error {
	panic("implement me UpdateAllMessageRead")
}

func (r *DaoWebSocketImpl) BatchUpdateMessageDeleted(ids []string) error {
	panic("implement me BatchUpdateMessageDeleted")
}

func NewDaoWebSocketImpl(ctx ...*base.Context) daos.DaoWebSocket {
	p := &DaoWebSocketImpl{}
	p.SetContext(ctx...)

	return p
}
