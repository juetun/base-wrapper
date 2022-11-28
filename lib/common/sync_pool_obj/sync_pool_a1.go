package sync_pool_obj

import (
	"fmt"
)

var (
	//定义要初始化对象的映射
	SliceNewHandlerSyncPool = map[string]NewHandler{}
	SyncPoolObjects         SyncPoolObjectMap
)

type (
	SyncPoolObjectMap map[string]*SyncPoolObject
	NewHandler        func() (res SyncPoolInterface)
)

//初始化调用
func init() {
	ReInit()
	return
}
func ReInit() {
	if SyncPoolObjects == nil {
		SyncPoolObjects = make(map[string]*SyncPoolObject, len(SliceNewHandlerSyncPool))
	}
	for key, handlerNewObject := range SliceNewHandlerSyncPool {
		if _, ok := SyncPoolObjects[key]; ok {
			continue
		}
		SyncPoolObjects[key] = NewSyncPoolObject(func() (res any) {
			return handlerNewObject()
		}, func(obj SyncPoolInterface) {
			obj.ReSet()
			return
		})
	}
}
func (r SyncPoolObjectMap) GetByKey(key string) (res *SyncPoolObject, err error) {
	var ok bool
	if res, ok = r[key]; !ok {
		err = fmt.Errorf("the sync pool(key:%s) object is not exists", key)
		return
	}
	return
}
