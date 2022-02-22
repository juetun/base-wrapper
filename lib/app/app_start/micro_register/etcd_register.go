package micro_register

import (
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juetun/base-wrapper/lib/app/micro_service"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"net/http"
	"strconv"
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

type ETCDRegister struct {
	syslog *base.SystemOut
}

func NewETCDRegister() (res *ETCDRegister) {
	res = &ETCDRegister{syslog:base.NewSystemOut()}
	return
}

func (r *ETCDRegister) RegisterMicro(c *gin.Engine) (ok bool, err error) {
	if r.microRun(c) {
		return
	}
	// listen and serve on 0.0.0.0:8080
	if err = c.Run(r.GetListenPortString()); err != nil {

		r.syslog.SetInfoType(base.LogLevelError).SystemOutPrintf("start err :%s", err.Error())
	}

	return
}

func (r *ETCDRegister) GetListenPortString() string {
	return ":" + strconv.Itoa(common.GetAppConfig().AppPort)
}

func (r *ETCDRegister) microRun(engine *gin.Engine) (res bool) {
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("Run as micro!")
	r.runAsMicro(engine)
	return
}

func (r *ETCDRegister) UnRegisterMicro() {

}

// RunAsMicro 使用go-micro实现服务注册与发现
func (r *ETCDRegister) runAsMicro(gin *gin.Engine) {
	var err error
	address := r.GetListenPortString()

	microServer := httpServer.NewServer(
		server.Name(common.GetAppConfig().AppName),
		server.Address(address),
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
		micro.Registry(
			newETCDRegistry(
				r.syslog,
				registry.Addrs(micro_service.ServiceConfig.Endpoints...),
				registry.Timeout(20*time.Second),
				registry.Secure(true),
			),
		))
	service.Init()
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Server(address:`%s`) init finished", address)
	if err = service.Run(); err != nil {
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("Server(address:`%s`) init Error:", err.Error())
	}
	return
}
