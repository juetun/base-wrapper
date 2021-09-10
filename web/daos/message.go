package daos

import (
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil"
)

type DaoWebSocket interface {
	websocket_anvil.MysqlInterface
}
