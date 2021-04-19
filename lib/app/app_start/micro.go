// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package app_start

import (
	"context"
	"crypto/tls"
	"fmt"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/util/addr"
	"github.com/google/uuid"
	"github.com/juetun/base-wrapper/lib/app/micro_service"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins/service_discory/traefik/etcd"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"
)

var (
	DefaultId = uuid.New().String()
)

//使用go-micro实现服务注册与发现
func (r *WebApplication) RunAsMicro(gin *gin.Engine) {
	var err error
	address := r.GetListenPortString()
	srv := httpServer.NewServer(

		server.Name(common.GetAppConfig().AppName),
		server.Address(address),
		server.RegisterTTL(time.Second*20),
		server.RegisterInterval(time.Second*15),
	)
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Server address:%s", address)
	hd := srv.NewHandler(gin)
	if err = srv.Handle(hd); err != nil {
		r.syslog.SetInfoType(base.LogLevelFatal).
			SystemOutFatalf("Register micro router failure!")
		return
	}

	service := micro.NewService(
		micro.Server(srv),
		micro.Registry(
			newEtcdRegistry(
				r.syslog,
				registry.Addrs(micro_service.ServiceConfig.Endpoints...),
				registry.Timeout(20*time.Second),
				registry.Secure(true),
			),
		),
	)
	//if err = etcd.NewTraefikEtcd(&micro_service.ServiceConfig).Action(); err != nil {
	//	r.syslog.SetInfoType(base.LogLevelFatal).SystemOutFatalf("registry server err(%#v) \n", err)
	//}
	service.Init()
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Server init finished")
	fmt.Println("")
	fmt.Println("")
	service.Run()

}

type EtcdRegistry struct {
	AppName    string `json:"app_name"`
	Address    string `json:"address"`
	Advertise  string `json:"advertise"`
	AppVersion string `json:"app_version"`
	ServerId   string `json:"server_id"`
	Metadata   map[string]string

	opts registry.Options
	mux  *http.ServeMux
	srv  *registry.Service

	sync.Mutex
	running bool
	static  bool
	exit    chan chan error

	syslog *base.SystemOut
}

func newEtcdRegistry(syslog *base.SystemOut, opts ...registry.Option) registry.Registry {
	options := NewOptions(opts...)
	res := &EtcdRegistry{
		ServerId: DefaultId,
		opts:     options,
		mux:      http.NewServeMux(),
		static:   true,
		syslog:   syslog,
	}
	res.GenSrv()
	return res
}

func NewOptions(opts ...registry.Option) registry.Options {
	opt := registry.Options{
		Addrs:   []string{""},
		Timeout: 30 * time.Second,
		Secure:  false,
		TLSConfig: &tls.Config{

		},
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func (r EtcdRegistry) GenSrv() {

	// default host:port
	parts := strings.Split(r.Address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(r.Advertise) > 0 {
		parts = strings.Split(r.Advertise, ":")

		// we have host:port
		if len(parts) > 1 {
			// set the host
			host = strings.Join(parts[:len(parts)-1], ":")

			// get the port
			if aport, _ := strconv.Atoi(parts[len(parts)-1]); aport > 0 {
				port = aport
			}
		} else {
			host = parts[0]
		}
	}

	addressIp, err := addr.Extract(host)
	if err != nil {
		addressIp = "127.0.0.1"
	}

	r.srv = &registry.Service{
		Name:    r.AppName,
		Version: r.AppVersion,
		Nodes: []*registry.Node{{
			Id:       r.ServerId,
			Address:  fmt.Sprintf("%s:%d", addressIp, port),
			Metadata: r.Metadata,
		}},
	}
	return
}

func (r EtcdRegistry) Init(opts ...registry.Option) error {

	//for _, o := range opts {
	//	o(&e.opts)
	//}
	//
	//serviceOpts := []micro.Option{}
	//
	//if len(e.opts.Flags) > 0 {
	//	serviceOpts = append(serviceOpts, micro.Flags(s.opts.Flags...))
	//}
	//
	//if s.opts.Registry != nil {
	//	serviceOpts = append(serviceOpts, micro.Registry(s.opts.Registry))
	//}
	//
	//serviceOpts = append(serviceOpts, micro.Action(func(ctx *cli.Context) {
	//	if ttl := ctx.Int("register_ttl"); ttl > 0 {
	//		e.opts.RegisterTTL = time.Duration(ttl) * time.Second
	//	}
	//
	//	if interval := ctx.Int("register_interval"); interval > 0 {
	//		e.opts.RegisterInterval = time.Duration(interval) * time.Second
	//	}
	//
	//	if name := ctx.String("server_name"); len(name) > 0 {
	//		e.opts.Name = name
	//	}
	//
	//	if ver := ctx.String("server_version"); len(ver) > 0 {
	//		e.opts.Version = ver
	//	}
	//
	//	if id := ctx.String("server_id"); len(id) > 0 {
	//		e.opts.Id = id
	//	}
	//
	//	if addr := ctx.String("server_address"); len(addr) > 0 {
	//		e.opts.Address = addr
	//	}
	//
	//	if adv := ctx.String("server_advertise"); len(adv) > 0 {
	//		e.opts.Advertise = adv
	//	}
	//
	//	if e.opts.Action != nil {
	//		e.opts.Action(ctx)
	//	}
	//}))
	//
	//e.opts.Service.Init(serviceOpts...)
	//srv := e.genSrv()
	//srv.Endpoints = e.srv.Endpoints
	//e.srv = srv

	return nil
}

func (r EtcdRegistry) Options() registry.Options {
	return r.opts
}

func (r EtcdRegistry) Register(service *registry.Service, option ...registry.RegisterOption) (err error) {
	log.Println("implement me Register")
	if err = etcd.NewTraefikEtcd(&micro_service.ServiceConfig).Action(); err != nil {
		r.syslog.SetInfoType(base.LogLevelFatal).SystemOutFatalf("registry server err(%#v) \n", err)
	}
	return
}

func (r EtcdRegistry) Deregister(service *registry.Service, option ...registry.DeregisterOption) (err error) {
	registry.Deregister(service)
	return
}

func (r EtcdRegistry) GetService(s string, option ...registry.GetOption) (res []*registry.Service,err error) {
	log.Println("implement me GetService")
	return
}

func (r *EtcdRegistry) ListServices(option ...registry.ListOption) (res []*registry.Service, err error) {
	log.Println("implement me ListServices")
	return
}

func (r *EtcdRegistry) Watch(option ...registry.WatchOption) (res registry.Watcher,err error) {
	log.Println("implement me Watch")
	return
}

func (r *EtcdRegistry) String() string {
	return "EtcdRegistry"
}
