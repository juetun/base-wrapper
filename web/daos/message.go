package daos

import (
	"github.com/juetun/base-wrapper/lib/common/anvil_websocket/ext"
)

type DaoWebSocket interface {
	ext.MysqlInterface
}
