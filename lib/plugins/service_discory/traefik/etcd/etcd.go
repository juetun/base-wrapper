// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/etcd-io/etcd/clientv3"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/micro_service"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory/traefik/discovery"
	"github.com/juetun/base-wrapper/lib/utils"
)

// 安装、学习文档:http://blueskykong.com/2020/06/06/etcd-3/
// https://pkg.go.dev/go.etcd.io/etcd/clientv3#pkg-overview
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
	var endpoints = make([]string, 0, len(serverRegistry.EtcdEndPoints))
	for _, endpoint := range serverRegistry.Endpoints {
		endpoints = append(endpoints, endpoint)
	}
	res.Client, res.Err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	res.ctx = context.Background()
	res.Lease = clientv3.NewLease(res.Client)
	res.Dir = serverRegistry.Dir
	res.syslog = syslog
	return
}

type KeyStruct struct {
	Key   []string `json:"key"`
	Value string   `json:"value"`
}

func (r *TraefikEtcd) getChildString(slice []KeyStruct) (res string) {
	var dt = make(map[string][]KeyStruct, len(slice))
	for _, keyStruct := range slice {

		if len(keyStruct.Key) <= 0 {
			switch keyStruct.Value {
			case "true", "false":
				res = fmt.Sprintf(`%s`, keyStruct.Value)
			default:
				res = fmt.Sprintf(`"%s"`, keyStruct.Value)
			}
			return
		}
		if _, ok := dt[keyStruct.Key[0]]; !ok {
			dt[keyStruct.Key[0]] = make([]KeyStruct, 0, 20)
		}
		v1 := KeyStruct{
			Key:   keyStruct.Key[1:],
			Value: keyStruct.Value,
		}
		dt[keyStruct.Key[0]] = append(dt[keyStruct.Key[0]], v1)

	}
	sSlice := make([]string, 0, len(slice))
	for k, it := range dt {
		sSlice = append(sSlice, fmt.Sprintf(`"%s":%s`, k, r.getChildString(it)))
	}

	res = fmt.Sprintf(`{%s}`, strings.Join(sSlice, ","))
	return
}
func (r *TraefikEtcd) parseMapToJsonByte(prefix string, mapv map[string]string) (res discovery.HttpTraefik, err error) {
	res = discovery.HttpTraefik{}
	var slice = make([]KeyStruct, 0, len(mapv))

	for key, v := range mapv {
		curString := strings.TrimPrefix(key, prefix)
		keySlice := strings.Split(curString, "/")
		slice = append(slice, KeyStruct{
			Key:   keySlice,
			Value: v,
		})
	}
	stringJson := r.getChildString(slice)
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("从etcd中获取json为: %s", stringJson)
	if err = json.Unmarshal([]byte(stringJson), &res); err != nil {
		return
	}
	return
}

type MicroServerSingleton struct {
	CurrentServer *discovery.HttpTraefik `json:"current_server"`
	ServiceName   string                 `json:"service_name"`
	KeyPrefixs    []string               `json:"key_prefixs"`
}

var serverConfig *MicroServerSingleton

func (r *TraefikEtcd) Action() (err error) {

	if serverConfig == nil { // 只执行一次动作
		serverConfig = &MicroServerSingleton{}
		serverConfig.CurrentServer, serverConfig.ServiceName, serverConfig.KeyPrefixs = r.readyServerData()
	}

	// 实现悲观锁锁定数据
	lid := r.lockService(serverConfig.ServiceName)
	defer r.unLockService(serverConfig.ServiceName, lid)

	var etcdMapValue map[string]string
	var etcdObject *discovery.TraefikConfig

	if etcdMapValue, err = r.getAllKey(serverConfig.KeyPrefixs); err != nil {
		return
	}
	etcdObject = &discovery.TraefikConfig{}
	etcdObject.Http, err = r.parseMapToJsonByte(discovery.HttpPrefix, etcdMapValue)
	if err != nil {
		return
	}
	mapValue := r.getTraefikConfigToKeyValue(etcdObject, serverConfig.CurrentServer)
	err = r.PutByTxt(mapValue)
	return
}

// 分布式锁
func (r *TraefikEtcd) lockService(serviceName string) (res clientv3.LeaseID) {
	return
}

// 分布式锁解锁
func (r *TraefikEtcd) unLockService(serviceName string, leaseID clientv3.LeaseID) {

	return
}
func (r *TraefikEtcd) mergeData(mapValue, nowData map[string]string) (res map[string]string) {

	return
}

// serviceName分布式锁的作用域
func (r *TraefikEtcd) getAllKey(prefixs []string) (res map[string]string, err error) {
	res = make(map[string]string, 50)
	var dt *clientv3.GetResponse
	for _, prefix := range prefixs {
		if dt, err = r.Client.Get(r.ctx, prefix, clientv3.WithPrefix()); err != nil {
			return
		}
		for _, it := range dt.Kvs {
			k := string(it.Key)
			v := string(it.Value)
			r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("etcd value:`%s = %v` [prefix:%s]\n", prefix, k, v)
			res[k] = v
		}
	}
	return
}
func (r *TraefikEtcd) sliceToMap(data []string) (res map[string]string) {
	res = make(map[string]string, len(data))
	for i, name := range data {
		res[strconv.Itoa(i)] = name
	}
	return
}
func (r *TraefikEtcd) getRouter(serviceName string, middlewaresNames ...string) (res map[string]discovery.HttpTraefikRouters, routerName string) {
	routerName = fmt.Sprintf("go-%s", serviceName)
	router := discovery.HttpTraefikRouters{
		EntryPoints: r.sliceToMap(micro_service.ServiceConfig.EtcdEndPoints),
		Rule:        fmt.Sprintf("Host(`%s`) && PathPrefix(`/%s`)", micro_service.ServiceConfig.Host, app_obj.App.AppName),
		Service:     serviceName,
		Middlewares: r.sliceToMap(middlewaresNames),
	}
	res = map[string]discovery.HttpTraefikRouters{
		routerName: router,
	}
	return
}

func (r *TraefikEtcd) getServices() (res map[string]discovery.HttpTraefikServiceConfig, serviceName string) {
	serviceName = app_obj.App.AppName
	ip, _ := utils.GetLocalIP()
	service := discovery.HttpTraefikServiceConfig{
		LoadBalancer: &discovery.HttpLoadBalancer{
			Servers: map[string]discovery.HttpLoadBalancerServer{
				"0": {
					Url: fmt.Sprintf("http://%s:%d", ip, app_obj.App.AppPort),
				},
			},
			PassHostHeader: true,
			HealthCheck: &discovery.HttpHealthCheck{
				Scheme:          "http://",
				Path:            "/health",
				Port:            strconv.Itoa(app_obj.App.AppPort),
				Hostname:        ip,
				FollowRedirects: true,
				Headers:         nil,
				Interval:        3 * time.Second,
				Timeout:         100 * time.Millisecond,
			},
		},
	}
	res = map[string]discovery.HttpTraefikServiceConfig{serviceName: service}
	return
}
func (r *TraefikEtcd) getServersTransports() (res map[string]discovery.HttpTraefikServersTransports, middlewaresName []string) {
	res = map[string]discovery.HttpTraefikServersTransports{}
	return
}
func (r *TraefikEtcd) getMiddlewares() (res map[string]discovery.HttpTraefikMiddleware, middlewaresName []string) {
	res = map[string]discovery.HttpTraefikMiddleware{}
	return
}

func (r *TraefikEtcd) readyServerData() (res *discovery.HttpTraefik, serviceName string, keyPrefix []string) {
	res = &discovery.HttpTraefik{}
	var routerName string
	var middlewaresName, serversTransportsName []string

	res.Services, serviceName = r.getServices()
	res.Middlewares, middlewaresName = r.getMiddlewares()
	res.Routers, routerName = r.getRouter(serviceName, middlewaresName...)
	res.ServersTransports, serversTransportsName = r.getServersTransports()

	// 获取要更新的Key前缀
	keyPrefix = r.getPrefixKeys(serviceName, routerName, middlewaresName, serversTransportsName)
	return
}

// 获取需要设置的参数
func (r *TraefikEtcd) getTraefikConfigToKeyValue(etcdTraefikConfig *discovery.TraefikConfig, currentServer *discovery.HttpTraefik) (res map[string]string) {
	etcdTraefikConfig.Http.MergeRouters(currentServer.Routers)
	etcdTraefikConfig.Http.MergeServices(currentServer.Services)
	etcdTraefikConfig.Http.MergeMiddlewares(currentServer.Middlewares)
	etcdTraefikConfig.Http.MergeServersTransports(currentServer.ServersTransports)
	res = etcdTraefikConfig.ToKV()
	return
}

func (r *TraefikEtcd) getPrefixKeys(serviceName, routerName string, middlewaresName, serversTransportsName []string) (keyPrefix []string) {
	keyPrefix = make([]string, 0, 2+len(middlewaresName)+len(serversTransportsName))
	keyPrefix = append(keyPrefix, fmt.Sprintf("traefik/http/services/%s", serviceName))
	keyPrefix = append(keyPrefix, fmt.Sprintf("traefik/http/routers/%s", routerName))
	for _, it := range middlewaresName {
		keyPrefix = append(keyPrefix, fmt.Sprintf("traefik/http/middlewares/%s", it))
	}
	for _, it := range serversTransportsName {
		keyPrefix = append(keyPrefix, fmt.Sprintf("traefik/http/serversTransports/%s", it))
	}
	return
}

// 将数据通过事务的方式提交到ETCD
func (r *TraefikEtcd) PutByTxt(mapValue map[string]string) (err error) {
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("开始将参数注册到ETCD【START】")
	defer func() {
		r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("将参数注册到ETCD【OVER】")
	}()
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
	)
	// 申请一个5秒的租约
	if leaseGrantResp, err = r.Lease.Grant(context.TODO(), 60); err != nil {
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
		// If(cmpOptions...).
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
	// if resp.OpResponse().Get().Count > 0 {
	//	res = true
	// }
	return
}

func (r *TraefikEtcd) Close() {
	r.Client.Close()
}
