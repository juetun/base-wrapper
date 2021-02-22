// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package plugins

import (
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/utils"
	"sync"
)

func PluginsHashId(arg  *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	hd := new(utils.HashIdParams)
	salt := hd.SetHashIdSalt("i must add a salt what is only for me")
	hdLength := hd.SetHashIdLength(8)
	zHashId, err := hd.HashIdInit(hdLength, salt)
	if err != nil {
		return
	}
	common.ZHashId = zHashId
	return
}
