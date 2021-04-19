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
	"github.com/juetun/base-wrapper/lib/plugins/service_discory"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory/traefik/discovery"
	"github.com/juetun/base-wrapper/lib/utils"
	"log"
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

	err   error
	Lease clientv3.Lease

	Dir string
}

func NewTraefikEtcd(serverRegistry *service_discory.ServerRegistry) (res *TraefikEtcd) {
	res = &TraefikEtcd{}
	res.Client, res.Err = clientv3.New(clientv3.Config{
		Endpoints:   serverRegistry.Endpoints,
		DialTimeout: 2 * time.Second,
	})
	res.ctx = context.Background()
	res.Lease = clientv3.NewLease(res.Client)
	res.Dir = serverRegistry.Dir
	log.Printf("etcd dir is:'%s'\n", res.Dir)
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
				Rule:        fmt.Sprintf("Host(`api.test.com`) && PathPrefix(`/%s`)", app_obj.App.AppName),
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
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
	)

	// 申请一个5秒的租约
	if leaseGrantResp, err = r.Lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// 准备一个用于取消自动续租的context
	ctx, cancelFunc := context.WithCancel(r.ctx)
	_ = cancelFunc
	// 确保函数退出后, 自动续租会停止
	defer func() {
		log.Printf("结束任务 \n")
		//cancelFunc()
	}()
	//defer r.Lease.Revoke(ctx, leaseId)

	// 5秒后会取消自动续租
	if keepRespChan, err = r.Lease.KeepAlive(ctx, leaseGrantResp.ID); err != nil {
		fmt.Println(err)
		return
	}

	// 处理续约应答的协程
	go func() {
		var keepResp *clientv3.LeaseKeepAliveResponse
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else { // 每秒会续租一次, 所以就会收到一次应答
					fmt.Println("收到自动续租应答:", keepResp.ID)
					break
				}
			}
		}
	END:
	}()

	//  if 不存在key， then 设置它, else 抢锁失败

	// 创建事务
	txn := clientv3.NewKV(r.Client).Txn(context.TODO())

	var listOptions = make([]clientv3.Op, 0, len(mapValue))
	var elseOptions = make([]clientv3.Op, 0, len(mapValue))
	var cmpOptions = make([]clientv3.Cmp, 0, len(mapValue))
	for k, v := range mapValue {
		log.Printf("%s = %s \n", k, v)
		listOptions = append(listOptions, clientv3.OpPut(k, v, clientv3.WithLease(leaseGrantResp.ID)))
		cmpOptions = append(cmpOptions, clientv3.Compare(clientv3.CreateRevision(k), "=", 0))
		elseOptions = append(elseOptions, clientv3.OpGet(k))
	}

	// 如果key不存在
	txn.If(cmpOptions...).
		Then(listOptions...).
		Else(elseOptions...) // 否则抢锁失败

	var txnResp *clientv3.TxnResponse
	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println("提交结果",err)
		return // 没有问题
	}
	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Printf("锁被占用:%#v  释放锁\n", txnResp.Responses[0].GetResponseRange().Kvs)
		return
	}
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
