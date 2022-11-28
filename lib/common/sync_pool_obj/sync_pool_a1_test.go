package sync_pool_obj

import (
	"testing"
)

type A1 struct{}

func NewA1() (res SyncPoolInterface) {
	return &A1{}
}
func (r *A1) ReSet() {
	return
}
func TestNewSyncPoolObject(t *testing.T) {
	type args struct {
		handlerNewObject HandlerNewObject
		resetHandler     HandlerReset
	}
	tests := []struct {
		name    string
		args    args
		wantRes *SyncPoolObject
	}{
		{
		},
	}
	SliceNewHandlerSyncPool["A1"] = NewA1
	ReInit()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, _ := SyncPoolObjects.GetByKey("A1")
			data := obj.Get().(*A1)
			defer func() {
				obj.Put(data)
			}()
			return
		})
	}
}
