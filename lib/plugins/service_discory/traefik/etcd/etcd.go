// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package etcd

import (
	"context"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/micro_service"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory/traefik/discovery"
	"github.com/juetun/base-wrapper/lib/utils"
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

	err    error
	Lease  clientv3.Lease
	syslog *base.SystemOut

	Dir string
}

func NewTraefikEtcd(serverRegistry *service_discory.ServerRegistry, syslog *base.SystemOut) (res *TraefikEtcd) {
	res = &TraefikEtcd{}
	res.Client, res.Err = clientv3.New(clientv3.Config{
		Endpoints:   serverRegistry.Endpoints,
		DialTimeout: 2 * time.Second,
	})
	res.ctx = context.Background()
	res.Lease = clientv3.NewLease(res.Client)
	res.Dir = serverRegistry.Dir
	res.syslog = syslog
	res.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("etcd dir is:'%s'\n", res.Dir)
	return
}

func (r *TraefikEtcd) Action() (err error) {
	mapValue := r.getTraefikConfigToKeyValue()
	err = r.PutByTxt(mapValue)
	return
}

func (r *TraefikEtcd) getTraefikConfigToKeyValue() (res map[string]string) {
	config := discovery.NewTraefikConfig()
	type Data struct{}
	ip, _ := utils.GetLocalIP()
	config.Http = discovery.HttpTraefik{
		Routers: map[string]discovery.HttpTraefikRouters{
			fmt.Sprintf("router-%s", app_obj.App.AppName): {
				EntryPoints: micro_service.ServiceConfig.EtcdEndPoints,
				Rule:        fmt.Sprintf("Host(`%s`) && PathPrefix(`/%s`)", micro_service.ServiceConfig.Host, app_obj.App.AppName),
				Service:     app_obj.App.AppName,
				Middlewares: []string{
					//"my-plugin",
				},
			},
		},
		Services: map[string]discovery.HttpTraefikServiceConfig{
			app_obj.App.AppName: {
				LoadBalancer: &discovery.HttpLoadBalancer{
					Servers: []discovery.HttpLoadBalancerServer{
						{
							Url: fmt.Sprintf("http://%s:%d", ip, app_obj.App.AppPort),
						},
					},
					PassHostHeader: true,
				},
			},
		},
		Middlewares: map[string]discovery.HttpTraefikMiddleware{
			//"my-plugin": Data{},
		},
	}

	res = config.ToKV()
	return
}

func (r *TraefikEtcd) PutByTxt(mapValue map[string]string) (err error) {
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("registry server message to etcd")
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
	)
	// 申请一个5秒的租约
	if leaseGrantResp, err = r.Lease.Grant(context.TODO(), 30); err != nil {
		r.syslog.SetInfoType(base.LogLevelError).SystemOutPrintf(err.Error())
		return
	}

	ctx, cancelFunc := context.WithTimeout(r.ctx, 3*time.Second)
	defer cancelFunc()
	// 创建事务
	txn := clientv3.NewKV(r.Client).Txn(ctx)
	var listOptions = make([]clientv3.Op, 0, len(mapValue))
	var elseOptions = make([]clientv3.Op, 0, len(mapValue))
	for k, v := range mapValue {
		r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("%s = %s \n", k, v)
		listOptions = append(listOptions, clientv3.OpPut(k, v, clientv3.WithLease(leaseGrantResp.ID)))
		elseOptions = append(elseOptions, clientv3.OpGet(k))
	}

	_, err = txn.
		//If(cmpOptions...).
		Then(listOptions...).
		Else(elseOptions...). // 否则抢锁失败
		Commit()
	return
}
func (r *TraefikEtcd) Put(Key, val string) (res bool, err error) {
	var resp *clientv3.PutResponse

	if resp, r.Err = r.Client.Put(r.ctx, Key, val); r.Err != nil {
		r.cancel()
		return
	}
	r.cancel()

	_ = resp
	//if resp.OpResponse().Get().Count > 0 {
	//	res = true
	//}
	return
}

func (r *TraefikEtcd) Close() {
	r.Client.Close()
}
