package sync_pool_obj

import "sync"

type SyncPoolInterface interface {
	ReSet() //重置属性值
}

type (
	SyncPoolObject struct {
		Pool  sync.Pool
		Reset HandlerReset //重置对象属性动作,用于清空结构体对象的所有属性的值
	}
	HandlerReset     func(syncPoolInterface SyncPoolInterface)
	HandlerNewObject func() (res any)
)

func (r *SyncPoolObject) Get() (res any) {
	return r.Pool.Get()
}

func (r *SyncPoolObject) Put(object SyncPoolInterface) {
	r.Reset(object) //还回对象时，清空对象的属性值
	r.Pool.Put(object)
	return
}

// handlerNewObject 初始化结构体对象
// resetHandler 重置结构体对象
func NewSyncPoolObject(handlerNewObject HandlerNewObject, resetHandler HandlerReset) (res *SyncPoolObject) {
	res = &SyncPoolObject{
		Pool: sync.Pool{
			New: func() any {
				return handlerNewObject()
			},
		},
		Reset: resetHandler,
	}
	return
}
