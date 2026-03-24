package micro_service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/library/common/tools"
	"net/http"
)

var io = base.NewSystemOut().
	SetInfoType(base.LogLevelInfo)

// 全局配置（可根据实际环境调整）

type (
	//注册服务参数到consul
	ConsulRegisterAndUnRegister struct {
		client       *api.Client    `json:"-"`
		ConsulConfig *ServiceConfig `json:"consul_config"`
	}

	// 服务配置
	ServiceConfig struct {
		ServiceName  string // 服务名（Traefik 会通过这个名字发现服务）
		Host         string // 服务监听地址
		Port         int    // 服务端口
		ConsulAddr   string // Consul 地址
		ServiceRoute string
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

func NewConsulRegisterAndUnRegister() (r app_start.MicroOperateInterface) {
	var (
		err   error
		ip, _ = tools.GetLocalIP() // 注意：如果 Traefik 在另一台机器，需改为本机内网 IP
		res   = &ConsulRegisterAndUnRegister{
			ConsulConfig: &ServiceConfig{
				Host:        ip,
				ServiceName: app_obj.App.AppName, // Traefik 会通过这个名称匹配服务
				Port:        app_obj.App.AppPort,
				ConsulAddr:  "127.0.0.1:8500", // Consul 默认地址
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

func (r *ConsulRegisterAndUnRegister) RegisterMicro(c *gin.Engine, cTxs ...context.Context) (ok bool, err error) {
	// 1. 初始化 Consul 客户端配置

	var traefikRouter = fmt.Sprintf("%v_router", r.ConsulConfig.ServiceName)
	var traefikEntry = "web"
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
		Tags: []string{
			"traefik.enable=true", // 告诉 Traefik 启用该服务
			fmt.Sprintf("traefik.http.routers.%v.rule=PathPrefix(`/%v`)", r.ConsulConfig.ServiceName, traefikRouter),             // Traefik 路由规则
			fmt.Sprintf("traefik.http.routers.%s.entrypoints=%s", traefikRouter, traefikEntry),                                   // 绑定入口点
			fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", r.ConsulConfig.ServiceName, r.ConsulConfig.Port), // 服务端口
		},
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
