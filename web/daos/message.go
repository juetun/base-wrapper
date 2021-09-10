package daos

import (
	"github.com/juetun/base-wrapper/lib/common/websocket_anvil/ext"
)

type DaoWebSocket interface {
	ext.MysqlInterface
}
