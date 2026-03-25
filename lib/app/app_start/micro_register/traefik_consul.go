package micro_service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
	"net/http"
	"strings"
)

var io = base.NewSystemOut().
	SetInfoType(base.LogLevelInfo)

// 全局配置（可根据实际环境调整）

type (
	//微服务注册逻辑
	MicroOperateInterface interface {
		//将服务信息注册入注册中心
		RegisterMicro(c *gin.Engine, cTxs ...context.Context) (ok bool, err error)

		GetMicroServiceTags() (tagList []string)

		//将服务信息从注册中心拿掉
		UnRegisterMicro()
	}

	//注册服务参数到consul
	ConsulRegisterAndUnRegister struct {
		client       *api.Client         `json:"-"`
		ConsulConfig *MicroServiceConfig `json:"consul_config"`
	}

	// 服务配置
	MicroServiceConfig struct {
		ServiceName string // 服务名（Traefik 会通过这个名字发现服务）
		Host        string // 服务监听地址
		Port        int    // 服务端口
		ConsulAddr  string // Consul 地址
	}
)

func (r *ConsulRegisterAndUnRegister) UnRegisterMicro() {
	if err := r.client.Agent().ServiceDeregister(r.getServiceId()); err != nil {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("注销服务失败: %v", err)
		return
	}
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("服务 %s 已从 Consul 注销\n", r.ConsulConfig.ServiceName)
	return
}

func NewConsulRegisterAndUnRegister() (r MicroOperateInterface) {
	var (
		err   error
		ip, _ = utils.GetLocalIP() // 注意：如果 Traefik 在另一台机器，需改为本机内网 IP
		res   = &ConsulRegisterAndUnRegister{
			ConsulConfig: &MicroServiceConfig{
				Host:        ip,
				ServiceName: app_obj.App.AppName, // Traefik 会通过这个名称匹配服务
				Port:        app_obj.App.AppPort,
				ConsulAddr:  strings.Join(app_obj.RegistryServiceConfig.Consul.Endpoints, ","), // Consul 默认地址
			},
		}
	)

	// 创建Consul客户端配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = res.ConsulConfig.ConsulAddr
	if res.client, err = api.NewClient(consulConfig); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Failed to create Consul client: %v", err.Error())
		return res
	}
	return res
}

func (r *ConsulRegisterAndUnRegister) getServiceId() (serviceId string) {
	serviceId = fmt.Sprintf("%s_%s_%d", r.ConsulConfig.ServiceName, r.ConsulConfig.Host, r.ConsulConfig.Port)
	return
}

func (r *ConsulRegisterAndUnRegister) GetMicroServiceTags() (tagList []string) {
	tagList = make([]string, 0, 30)
	serviceRouteWs := fmt.Sprintf("ws-%v", app_obj.App.AppName)
	websocketMiddleWareName := "websocket-long"
	middleWareCors := "cors-allow"

	tagList = append(tagList, []string{
		"traefik.enable=true", // 告诉 Traefik 启用该服务
		fmt.Sprintf("traefik.http.middlewares.%v.websocket.timeout=%v", websocketMiddleWareName, app_obj.RegistryServiceConfig.Consul.WebsocketConnectTimeOut), // 24小时不断开（单位：秒）
		fmt.Sprintf("traefik.http.routers.%v.rule=PathPrefix(`/%v`)", serviceRouteWs, r.ConsulConfig.ServiceName),                                              // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%v.service=%v", serviceRouteWs, r.ConsulConfig.ServiceName),                                                          // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=%s", serviceRouteWs, "web,websecure"),                                                                 //支持80和443端口
		fmt.Sprintf("traefik.http.routers.%s.middlewares=%v", serviceRouteWs, fmt.Sprintf("%v,%v", middleWareCors, websocketMiddleWareName)),                   //配置中间件信息
		fmt.Sprintf("traefik.http.routers.%s.tls=true", serviceRouteWs),                                                                                        // 绑定入口点 // 绑定入口点
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", r.ConsulConfig.ServiceName, r.ConsulConfig.Port),                                   // 服务端口

		//配置跨域信息
		fmt.Sprintf("traefik.http.middlewares.%v.headers.accesscontrolalloworiginlist=*", middleWareCors),
		fmt.Sprintf("traefik.http.middlewares.%v.headers.accesscontrolallowmethods=*", middleWareCors),
		fmt.Sprintf("traefik.http.middlewares.%v.headers.accesscontrolallowheaders=*", middleWareCors),
	}...)

	//如果支持websocket
	if app_obj.RegistryServiceConfig.SupportWebsocket {
		tagList = append(tagList, r.orgWss(websocketMiddleWareName, middleWareCors)...)
	}
	return
}

func (r *ConsulRegisterAndUnRegister) orgWss(websocketMiddleWareName, middleWareCors string) (tagList []string) {
	serviceRouteWss := fmt.Sprintf("wss-%v", app_obj.App.AppName)
	tagList = []string{
		fmt.Sprintf("traefik.http.routers.%v.rule=PathPrefix(`/%v`)", serviceRouteWss, r.ConsulConfig.ServiceName),                            // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%v.service=%v", serviceRouteWss, r.ConsulConfig.ServiceName),                                        // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=%s", serviceRouteWss, "web,websecure"),                                               //支持80和443端口
		fmt.Sprintf("traefik.http.routers.%s.middlewares=%v", serviceRouteWss, fmt.Sprintf("%v,%v", middleWareCors, websocketMiddleWareName)), //配置中间件信息
		fmt.Sprintf("traefik.http.routers.%s.tls=true", serviceRouteWss),                                                                      // 绑定入口点
	}
	return
}

func (r *ConsulRegisterAndUnRegister) RegisterMicro(c *gin.Engine, cTxs ...context.Context) (ok bool, err error) {
	// 1. 初始化 Consul 客户端配置

	var checkUri = fmt.Sprintf("http://%v:%v/%v/%v/heart_beat", r.ConsulConfig.Host, r.ConsulConfig.Port, r.ConsulConfig.ServiceName, app_obj.RouteTypeDefaultIntranet)

	// 3. 构建服务注册信息
	registration := &consulapi.AgentServiceRegistration{
		ID:      r.getServiceId(),           // 服务唯一 ID
		Name:    r.ConsulConfig.ServiceName, // 服务名（核心，Traefik 依赖这个）
		Address: r.ConsulConfig.Host,        // 服务地址
		Port:    r.ConsulConfig.Port,        // 服务端口
		Check: &consulapi.AgentServiceCheck{ // 健康检查配置（Consul 定期调用，确保服务存活）
			HTTP:     checkUri, // 健康检查接口
			Method:   http.MethodHead,
			Interval: "5s", // 检查间隔
			Timeout:  "3s", // 超时时间
			//TTL:                            "10s", // 存活时间
			DeregisterCriticalServiceAfter: "30s", // 不健康后 30s 注销服务
		},

		// 自定义标签（可选，Traefik 可通过标签过滤服务）
		Tags: r.GetMicroServiceTags(),
	}
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("健康检查地址:%v(%v) \n", checkUri, registration.Check.Method)
	// 4. 注册服务到 Consul
	if err = r.client.Agent().ServiceRegister(registration); err != nil {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("注册服务到 Consul 失败: %v", err)
		return
	}
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("服务 %s 已成功注册到 Consul (地址: %s:%d)\n", r.ConsulConfig.ServiceName, r.ConsulConfig.Host, r.ConsulConfig.Port)
	ok = true
	return
}
