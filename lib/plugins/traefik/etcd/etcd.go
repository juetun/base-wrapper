// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//安装、学习文档:http://blueskykong.com/2020/06/06/etcd-3/
//https://pkg.go.dev/go.etcd.io/etcd/clientv3#pkg-overview
type TraefikEtcd struct {
	Client  *clientv3.Client
	Err     error
	Timeout time.Duration `json:"timeout"`
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTraefikEtcd() (res *TraefikEtcd) {
	res = &TraefikEtcd{}
	res.Client, res.Err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	ctx := context.Background()
	res.ctx, res.cancel = context.WithTimeout(ctx, res.Timeout)
	return
}
func (r *TraefikEtcd) PutByTxt(mapValue map[string]string) (res *clientv3.TxnResponse, err error) {
	_, cancel := context.WithTimeout(r.ctx, r.Timeout)
	kvc := clientv3.NewKV(r.Client)
	dt := kvc.Txn(r.ctx)
	var params = make([]clientv3.Op, 0, len(mapValue))
	for s, s2 := range mapValue {
		params = append(params, clientv3.OpPut(s, s2))
	}
	res, err = dt.Then(params...).
		Commit()
	cancel()
	return
}
func (r *TraefikEtcd) Put(Key, val string) (res bool, err error) {
	var resp *clientv3.PutResponse

	if resp, r.Err = r.Client.Put(r.ctx, Key, val); r.Err != nil {
		r.cancel()
		return
	}
	r.cancel()
	if resp.OpResponse().Get().Count > 0 {
		res = true
	}
	return
}

func (r *TraefikEtcd) Close() {
	r.Client.Close()
}
