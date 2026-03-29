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
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"time"
)

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
		ServiceNameNotPrefix string //无前缀的服务名
		ServiceName          string // 服务名（Traefik 会通过这个名字发现服务）
		Host                 string // 服务监听地址
		Port                 int    // 服务端口
		ConsulAddr           string // Consul 地址
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
				Host:                 ip,
				ServiceNameNotPrefix: app_obj.App.AppName,
				ServiceName:          fmt.Sprintf("%v-%v", app_obj.App.AppEnv, app_obj.App.AppName), // Traefik 会通过这个名称匹配服务
				Port:                 app_obj.App.AppPort,
				ConsulAddr:           strings.Join(app_obj.RegistryServiceConfig.Consul.Endpoints, ","), // Consul 默认地址
			},
		}
	)

	// 创建Consul客户端配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = res.ConsulConfig.ConsulAddr
	if res.client, err = api.NewClient(consulConfig); err != nil {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Failed to create Consul client: %v", err.Error())
		return res
	}
	return res
}

func (r *ConsulRegisterAndUnRegister) getServiceId() (serviceId string) {
	serviceId = fmt.Sprintf("%s_%s_%d", r.ConsulConfig.ServiceName, r.ConsulConfig.Host, r.ConsulConfig.Port)
	return
}

func (r *ConsulRegisterAndUnRegister) getMapAppConfig() (rules []*app_obj.ConsulAppRuleInfo) {
	return app_obj.RegistryServiceConfig.Consul.MapApp
}

func (r *ConsulRegisterAndUnRegister) parseMapAppRule() (mapAppRule map[string]*app_obj.ConsulAppRuleInfo) {
	var (
		rules = app_obj.RegistryServiceConfig.Consul.MapApp
		item  *app_obj.ConsulAppRuleInfo
	)
	mapAppRule = make(map[string]*app_obj.ConsulAppRuleInfo, len(rules))
	for _, item = range rules {
		mapAppRule[item.MicroAppName] = item
	}
	return
}

func (r *ConsulRegisterAndUnRegister) orgPathHost(appRule *app_obj.ConsulAppRuleInfo, ruleString *strings.Builder) {
	hostLength := len(appRule.Host)

	if hostLength <= 0 {
		return
	}

	if hostLength == 1 {
		ruleString.WriteString(fmt.Sprintf("Host(`%v`)", appRule.Host[0]))
		return
	}

	ruleString.WriteString("(")
	for k, item := range appRule.Host {
		if item == "" {
			continue
		}
		if k != 0 {
			ruleString.WriteString("||")
		}
		ruleString.WriteString(fmt.Sprintf("Host(`%v`)", item))
	}
	ruleString.WriteString(")")

	return
}

func (r *ConsulRegisterAndUnRegister) orgPathPrefix(appRule *app_obj.ConsulAppRuleInfo, ruleString *strings.Builder) {
	if appRule.PathPrefix == "" {
		return
	}
	if ruleString.Len() != 0 {
		ruleString.WriteString("&&")
	}
	ruleString.WriteString(fmt.Sprintf("PathPrefix(`%v`)", appRule.PathPrefix))
	return
}

func (r *ConsulRegisterAndUnRegister) orgRoutesRule(serviceRoute string) (res string) {

	var (
		appRule    *app_obj.ConsulAppRuleInfo
		ok         bool
		mapAppRule = r.parseMapAppRule()
		ruleString strings.Builder
	)

	if appRule, ok = mapAppRule[app_obj.App.AppName]; !ok {
		res = fmt.Sprintf("traefik.http.routers.%v.rule=PathPrefix(`/%v`)", serviceRoute, r.ConsulConfig.ServiceNameNotPrefix)
		return
	}

	r.orgPathHost(appRule, &ruleString)
	r.orgPathPrefix(appRule, &ruleString)
	res = fmt.Sprintf("traefik.http.routers.%v.rule=%v", serviceRoute, ruleString.String())
	return
}

func (r *ConsulRegisterAndUnRegister) GetMicroServiceTags() (tagList []string) {
	tagList = make([]string, 0, 30)
	serviceRouteWs := fmt.Sprintf("ws-%v", app_obj.App.AppName)
	middleWareCors := "cors-allow"
	tagList = make([]string, 0, 30)
	tagList = append(tagList, []string{
		"traefik.enable=true",           // 告诉 Traefik 启用该服务
		r.orgRoutesRule(serviceRouteWs), // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%v.service=%v", serviceRouteWs, r.ConsulConfig.ServiceName), // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=%s", serviceRouteWs, "web"),
		fmt.Sprintf("traefik.http.routers.%s.middlewares=%s", serviceRouteWs, middleWareCors),                                //设置中间件
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", r.ConsulConfig.ServiceName, r.ConsulConfig.Port), // 服务端口

		//配置跨域信息
		fmt.Sprintf("traefik.http.middlewares.%v.headers.accesscontrolalloworiginlist=*", middleWareCors),
		fmt.Sprintf("traefik.http.middlewares.%v.headers.accesscontrolallowmethods=*", middleWareCors),
		fmt.Sprintf("traefik.http.middlewares.%v.headers.accesscontrolallowheaders=*", middleWareCors),
	}...)

	//如果支持websocket
	if app_obj.RegistryServiceConfig.SupportWebsocket {
		tagList = append(tagList, r.orgWss(middleWareCors)...)
	}
	return
}

func (r *ConsulRegisterAndUnRegister) orgWss(middleWareCors string) (tagList []string) {
	serviceRouteWss := fmt.Sprintf("wss-%v", app_obj.App.AppName)
	tagList = []string{
		r.orgRoutesRule(serviceRouteWss),
 		fmt.Sprintf("traefik.http.routers.%v.service=%v", serviceRouteWss, r.ConsulConfig.ServiceName),                      // Traefik 路由规则
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=%s", serviceRouteWss, "websecure"),                                 // 支持websocket
		fmt.Sprintf("traefik.http.routers.%s.middlewares=%s", serviceRouteWss, middleWareCors),                              // 设置中间件
		//fmt.Sprintf("traefik.http.routers.%s.middlewares=%v", serviceRouteWss, fmt.Sprintf("%v,%v", middleWareCors,websocketMiddleWareName)),
		fmt.Sprintf("traefik.http.routers.%s.tls=false", serviceRouteWss),           // 绑定入口点
		fmt.Sprintf("traefik.tcp.routers.%v.tls.passthrough=true", serviceRouteWss), // 开启TLS 透传 Traefik 不解密 HTTPS/WSS 流量，直接把加密数据【原封不动】转发给后端服务。
	}
	return
}
func (r *ConsulRegisterAndUnRegister) registerAction(cTxs ...context.Context) (err error) {

	// 3. 构建服务注册信息
	registration := &consulapi.AgentServiceRegistration{
		ID:      r.getServiceId(),           // 服务唯一 ID
		Name:    r.ConsulConfig.ServiceName, // 服务名（核心，Traefik 依赖这个）
		Address: r.ConsulConfig.Host,        // 服务地址
		Port:    r.ConsulConfig.Port,        // 服务端口
		Check: &consulapi.AgentServiceCheck{ // 健康检查配置（Consul 定期调用，确保服务存活）
			HTTP:     fmt.Sprintf("http://%v:%v%v", r.ConsulConfig.Host, r.ConsulConfig.Port, app_obj.GetHealthPath("health")), // 健康检查接口
			Method:   http.MethodHead,
			Interval: "5s", // 检查间隔
			Timeout:  "3s", // 超时时间
			//TTL:                            "10s", // 存活时间
			DeregisterCriticalServiceAfter: "30s", // 不健康后 30s 注销服务
		},

		// 自定义标签（可选，Traefik 可通过标签过滤服务）
		Tags: r.GetMicroServiceTags(),
	}
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("健康检查地址:%v(%v) Interval:%v  Timeout:%v\n",
		registration.Check.HTTP,
		registration.Check.Method,
		registration.Check.Interval,
		registration.Check.Timeout,
	)
	// 4. 注册服务到 Consul
	if err = r.client.Agent().ServiceRegister(registration); err != nil {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("注册服务到 Consul 失败: %v", err)
		return
	}
	base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("服务 %s 已成功注册到 Consul (地址: %s:%d)\n", r.ConsulConfig.ServiceName, r.ConsulConfig.Host, r.ConsulConfig.Port)

	return
}

func (r *ConsulRegisterAndUnRegister) RegisterMicro(c *gin.Engine, cTxs ...context.Context) (ok bool, err error) {

	if err = r.registerAction(cTxs...); err != nil {
		return
	}

	//设置注册微服务更新时间
	app_obj.SetLastRegisterTime()
	time.AfterFunc(3*time.Second, func() {
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("延迟4秒启动%v检测微服务心跳任务\n", app_obj.App.AppName)
		c := cron.New()
		// 每隔2秒执行一次任务，可以使用 Cron 表达式来定义更复杂的调度规则，例如 "*/4 * * * * *" 每4秒执行一次。
		c.AddFunc("*/4 * * * * *", func() {

			//如果consul不主动发健康检查请求,侧服务器尝试再次注册服务
			if time.Now().Unix()-app_obj.GetLastRegisterTime() < 10 { //如果心跳检测时间小于10秒,不主动发起请求
				return
			}

			//重新注册服务
			_ = r.registerAction(cTxs...)

		}) // 注意：这里的 Cron 表达式在不同的环境下可能有细微差别，例如在某些系统上可能需要使用 "0/2 * * * * *"。请根据实际环境调整。
		c.Start()
	})
	ok = true
	return
}
