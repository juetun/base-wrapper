package micro_register

import (
	"context"
	"crypto/tls"
	"fmt"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/util/addr"
	"github.com/gin-gonic/gin"
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
)

var (
	DefaultId = uuid.New().String()
)

func newETCDRegistry(syslog *base.SystemOut, opts ...registry.Option) registry.Registry {
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

type (
	ETCDRegister struct {
		syslog      *base.SystemOut
		microServer server.Server
	}
	EtcdRegistry struct {
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

		etcdTraefik *etcd.TraefikEtcd
	}
)

func NewETCDRegister() (res *ETCDRegister) {
	res = &ETCDRegister{
		syslog: base.NewSystemOut(),
	}
	return
}

func (r *ETCDRegister) RegisterMicro(c *gin.Engine, cTxs ...context.Context) (ok bool, err error) {
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("Run as micro!")
	r.microServer = r.runAsMicro(c, cTxs...)
	return
}

func (r *ETCDRegister) UnRegisterMicro() {
	//停止微服务注册
	if r.microServer != nil {
		r.microServer.Stop()
	}
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("停止微服务注册!")
}

func (r *ETCDRegister) GetListenPortString() string {
	return ":" + strconv.Itoa(common.GetAppConfig().AppPort)
}

// RunAsMicro 使用go-micro实现服务注册与发现
func (r *ETCDRegister) runAsMicro(gin *gin.Engine, cTxs ...context.Context) (microServer server.Server) {
	var err error
	address := r.GetListenPortString()
	etcdRegistry := newETCDRegistry(
		r.syslog,
		registry.Addrs(micro_service.ServiceConfig.Endpoints...),
		registry.Timeout(20*time.Second),
		registry.Secure(true),
	)
	var ctx context.Context
	if len(cTxs) > 0 {
		ctx = cTxs[0]
	}
	microServer = httpServer.NewServer(
		server.Name(common.GetAppConfig().AppName),
		server.Address(address),
		server.Context(ctx),
		server.RegisterTTL(time.Second*10),
		server.RegisterInterval(time.Second*5),
	)

	if err = microServer.Handle(microServer.NewHandler(gin)); err != nil {
		r.syslog.SetInfoType(base.LogLevelFatal).
			SystemOutFatalf("Register micro router failure!")
		return
	}

	service := micro.NewService(
		micro.Server(microServer),
		micro.Registry(etcdRegistry))
	service.Init()
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Server(address:`%s`) init finished", address)
	if err = service.Run(); err != nil {
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("Server(address:`%s`) init Error:", err.Error())
	}
	return
}

func NewOptions(opts ...registry.Option) registry.Options {
	opt := registry.Options{
		Addrs:     []string{""},
		Timeout:   22 * time.Second,
		Secure:    false,
		TLSConfig: &tls.Config{},
		Context:   context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func (r *EtcdRegistry) GenSrv() {

	// default host:port
	parts := strings.Split(r.Address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	// check the advertisement address first
	// if it exists then use it, otherwise
	// use the address
	if len(r.Advertise) > 0 {
		parts = strings.Split(r.Advertise, ":")

		// we have Host:port
		if len(parts) > 1 {
			// set the host
			host = strings.Join(parts[:len(parts)-1], ":")

			// get the port
			if addressPort, _ := strconv.Atoi(parts[len(parts)-1]); addressPort > 0 {
				port = addressPort
			}
		} else {
			host = parts[0]
		}
	}

	var addressIp string
	var err error

	if addressIp, err = addr.Extract(host); err != nil {
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

func (r *EtcdRegistry) Init(opts ...registry.Option) error {
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintln("EtcdRegistry implement me Init")
	for _, o := range opts {
		o(&r.opts)
	}

	return nil
}

func (r *EtcdRegistry) Options() registry.Options {
	return r.opts
}

func (r *EtcdRegistry) Register(service *registry.Service, option ...registry.RegisterOption) (err error) {
	_ = service
	_ = option
	if r.etcdTraefik == nil {
		r.etcdTraefik = etcd.NewTraefikEtcd(&micro_service.ServiceConfig, r.syslog)
	}
	if err = r.etcdTraefik.Action(); err != nil {
		r.syslog.SetInfoType(base.LogLevelFatal).SystemOutFatalf("registry server err(%#v) \n", err)
	}
	return
}

func (r *EtcdRegistry) Deregister(service *registry.Service, option ...registry.DeregisterOption) (err error) {
	_ = option
	err = registry.Deregister(service)
	return
}

func (r *EtcdRegistry) GetService(s string, option ...registry.GetOption) (res []*registry.Service, err error) {
	_ = s
	_ = option
	log.Println("implement me GetService")
	return
}

func (r *EtcdRegistry) ListServices(option ...registry.ListOption) (res []*registry.Service, err error) {
	_ = option
	log.Println("implement me ListServices")
	return
}

func (r *EtcdRegistry) Watch(option ...registry.WatchOption) (res registry.Watcher, err error) {
	_ = option
	log.Println("implement me Watch")
	return
}

func (r *EtcdRegistry) String() string {
	return "EtcdRegistry"
}
