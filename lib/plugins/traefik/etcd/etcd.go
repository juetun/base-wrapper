// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"time"
)

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
	res.ctx, res.cancel = context.WithTimeout(context.Background(), res.Timeout)
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
