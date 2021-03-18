// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package app_start

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/util/addr"
	"github.com/google/uuid"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	DefaultId = uuid.New().String()
)

func (r *WebApplication) RunAsMicro() {
	service := micro.NewService(
		micro.Name(app_obj.App.AppName), //注册的服务名称
		micro.Registry(newEtcdRegistry()),
		micro.RegisterInterval(time.Second*15), //每隔15秒重新注册一次
		micro.RegisterTTL(time.Second*30),      //注册服务的过期时间
	)

	service.Init()

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
}

func newEtcdRegistry(opts ...registry.Option) (res EtcdRegistry) {
	options := NewOptions(opts...)
	res = EtcdRegistry{
		ServerId: DefaultId,
		opts:     options,
		mux:      http.NewServeMux(),
		static:   true,
	}
	res.srv = res.genSrv()
	return res
}
func (e *EtcdRegistry) genSrv() *registry.Service {
	// default host:port
	parts := strings.Split(e.Address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(e.Advertise) > 0 {
		parts = strings.Split(e.Advertise, ":")

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
		// best effort localhost
		addressIp = "127.0.0.1"
	}

	return &registry.Service{
		Name:    e.AppName,
		Version: e.AppVersion,
		Nodes: []*registry.Node{{
			Id:       e.ServerId,
			Address:  fmt.Sprintf("%s:%d", addressIp, port),
			Metadata: e.Metadata,
		}},
	}
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

func (e EtcdRegistry) Init(opts ...registry.Option) error {
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

func (e EtcdRegistry) Options() registry.Options {
	return e.opts
}

func (e EtcdRegistry) Register(service *registry.Service, option ...registry.RegisterOption) (err error) {
	panic("implement me")
}

func (e EtcdRegistry) Deregister(service *registry.Service, option ...registry.DeregisterOption) (err error) {
	registry.Deregister(service)
	return
}

func (e EtcdRegistry) GetService(s string, option ...registry.GetOption) ([]*registry.Service, error) {
	panic("implement me")
}

func (e EtcdRegistry) ListServices(option ...registry.ListOption) ([]*registry.Service, error) {
	panic("implement me")
}

func (e EtcdRegistry) Watch(option ...registry.WatchOption) (registry.Watcher, error) {
	panic("implement me")
}

func (e EtcdRegistry) String() string {
	panic("implement me")
}
