package common

import (
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/speps/go-hashids"
	"sync"
)

var (
	ZHashId *hashids.HashID
)

func PluginsHashId() (err error) {
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
	ZHashId = zHashId
	return
}
