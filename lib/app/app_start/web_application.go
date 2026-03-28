// Package app_start
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package app_start

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	microService2 "github.com/juetun/base-wrapper/lib/app/app_start/micro_register"
	"github.com/juetun/base-wrapper/lib/app/middlewares"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	HandleFunc            = make([]HandleRouter, 0)        // 路由函数切片
	HandleFuncIntranet    = make([]HandleRouter, 0)        // 内网路由函数切片
	HandleFuncOuterNet    = make([]HandleRouter, 0)        // 外网路由函数切片
	HandleFuncAdminNet    = make([]HandleRouter, 0)        // 外网路由函数切片
	HandleFuncPage        = make([]HandleRouter, 0)        // 外网路由函数切片
	RoutePathInitCallBack RouterPath                       //注册路由时调用
	PermitAdminUrlPath    = make([]*PermitUrlPath, 0, 200) //收集应用内路由信息
	AdminNetHandlerFunc   = make([]gin.HandlerFunc, 0, 5)  //管理后台操作中间件
)

type (
	WebApplication struct {
		GinEngine *gin.Engine
		syslog    *base.SystemOut
		//FlagMicro bool // 如果是支持微服务
		MicroOperate microService2.MicroOperateInterface
	}
	PermitUrlPath struct {
		Method string `json:"method"`
		Uri    string `json:"uri"`
	}
	// HandleRouter 路由注册函数
	HandleRouter func(c *gin.Engine, urlPrefix string)
	RouterPath   func(httpMethod, absolutePath, handlerName string, nuHandlers int)
)

// NewWebApplication privateMiddleWares 项目自定义的GIN中间件
func NewWebApplication(privateMiddleWares ...gin.HandlerFunc) *WebApplication {
	switch strings.ToLower(common.GetAppConfig().AppEnv) {
	case app_obj.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	case app_obj.EnvTest:
		gin.SetMode(gin.TestMode)
	case app_obj.EnvDev:
		gin.SetMode(gin.DebugMode)
	case app_obj.EnvPre:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	webApp := &WebApplication{
		GinEngine: gin.New(),
		syslog:    base.NewSystemOut(),
	}
	if RoutePathInitCallBack != nil {
		gin.DebugPrintRouteFunc = RoutePathInitCallBack
	} else {
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			PermitAdminUrlPath = append(PermitAdminUrlPath, &PermitUrlPath{Method: httpMethod, Uri: absolutePath})
			//路由结构修改
			webApp.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("ROUTE_PATH: %v %v %v (%v Handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
		}
	}

	// 加载GIN框架 中间件
	middlewares.LoadMiddleWare(privateMiddleWares...)

	// gin加载中间件
	webApp.GinEngine.Use(middlewares.MiddleWareComponent..., )

	return webApp
}

type RouterHandler func(r *gin.Engine) (err error)

func (r *WebApplication) SetFlagMicro(micro microService2.MicroOperateInterface) (res *WebApplication) {
	res = r
	r.MicroOperate = micro
	return
}

// LoadRouter 加载API路由
func (r *WebApplication) LoadRouter(routerHandler ...RouterHandler) (res *WebApplication) {
	res = r
	var err error
	defer func() {
		if err != nil {
			r.syslog.SetInfoType(base.LogLevelError).
				SystemOutPrintf("Load router err  %s", err.Error())
		}
	}()
	appConfig := common.GetAppConfig()
	// var UrlPrefix = fmt.Sprintf("%s/%s", appConfig.AppName, appConfig.AppApiVersion)
	var UrlPrefix = fmt.Sprintf("%s", appConfig.AppName)

	fmt.Printf("\n\n")
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("***********************开始注册路由(app_name:%s)*****************", UrlPrefix)
	defer func() {
		fmt.Printf("\n\n")
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("***********************注册路由结束(app_name:%s)***********************", UrlPrefix)
	}()

	// 工具路由注册（心跳检测、性能分析等）
	r.toolRouteRegister(appConfig, UrlPrefix)

	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("注册业务路由....\n\n")

	// 操作路由相关动作
	for _, handler := range routerHandler {
		if err = handler(r.GinEngine); err != nil {
			return
		}
	}
	// URL路由注册操作
	for _, router := range HandleFunc {
		router(r.GinEngine, UrlPrefix)
	}

	if len(HandleFuncIntranet) > 0 {
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintln("注册内网访问接口路由....")
		for _, router := range HandleFuncIntranet {
			router(r.GinEngine, fmt.Sprintf("%s/%s", UrlPrefix, app_obj.App.AppRouterPrefix.Intranet))
		}
	}
	if len(HandleFuncOuterNet) > 0 {
		fmt.Printf("\n")
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintln("注册外网访问接口路由....")
		for _, router := range HandleFuncOuterNet {
			router(r.GinEngine, fmt.Sprintf("%s/%s", UrlPrefix, app_obj.App.AppRouterPrefix.Outranet))
		}
	}
	if len(HandleFuncAdminNet) > 0 {
		fmt.Printf("\n")
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintln("注册客服后台访问接口路由....")
		AdminNetHandlerFunc = append([]gin.HandlerFunc{middlewares.AdminMiddlewares()}, AdminNetHandlerFunc...)
		for _, router := range HandleFuncAdminNet {
			r.GinEngine.Use(AdminNetHandlerFunc...)
			router(r.GinEngine, fmt.Sprintf("%s/%s", UrlPrefix, app_obj.App.AppRouterPrefix.AdminNet))
		}
	}
	if len(HandleFuncPage) > 0 {
		fmt.Printf("\n")
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("注册网页界面访问路由(%s).... \n", UrlPrefix)
		pr := app_obj.App.AppRouterPrefix.Page
		if pr != "" {
			pr = "/" + pr
		}
		for _, router := range HandleFuncPage {
			router(r.GinEngine, fmt.Sprintf("%s%s", UrlPrefix, pr))
		}
	}
	return
}

func (r *WebApplication) getCtx(cTxs ...context.Context) (ctx context.Context) {
	if len(cTxs) > 0 {
		ctx = cTxs[0]
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return
}

//微服务注册信息
func (r *WebApplication) loadRegistryMicroSrv() {
	if !app_obj.RegistryServiceConfig.OpenMicService {
		return
	}

	if r.MicroOperate == nil {
		//如果开启了微服务注册
		switch app_obj.RegistryServiceConfig.MicServiceType {
		case app_obj.RegisterCenterConsul: //注册中心为Consul
			r.SetFlagMicro(microService2.NewConsulRegisterAndUnRegister())
		case app_obj.RegisterCenterETCD: //注册中心为 ETCD
			r.SetFlagMicro(microService2.NewETCDRegisterAndUnRegister())
		}
	}

	return
}

// Run 开始加载Gin 服务
func (r *WebApplication) Run(cTxs ...context.Context) (err error) {

	ctx := r.getCtx(cTxs...)
	// 5. 启动 Gin 服务
	var server = &http.Server{
		Addr:    fmt.Sprintf(":%d", app_obj.App.AppPort),
		Handler: r.GinEngine,
	}
	// // 如果支持优雅重启（微服务启动）
	if app_obj.RegistryServiceConfig.OpenMicService {
		r.loadRegistryMicroSrv()
		r.startWithMicro(ctx, server)
		return
	}

	//普通启动
	r.startGeneral(server)

	return
}

func (r *WebApplication) startGeneral(server *http.Server) {

	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("开始普通启动信息%v \n", server.Addr)
	var err error
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("服务启动失败: %v \n", err)
	}
	return
}

func (r *WebApplication) startWithMicro(ctx context.Context, server *http.Server) {
	if r.MicroOperate == nil {
		return
	}
	if _, err := r.MicroOperate.RegisterMicro(r.GinEngine, ctx); err != nil {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("注册服务失败: %v", err)
		return
	}

	// 确保程序退出时注销服务
	defer r.MicroOperate.UnRegisterMicro()

	// 6. 优雅退出（监听系统信号）
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("开始优雅关闭服务...")

		// 关闭 HTTP 服务
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("服务关闭失败: %v", err)
			return
		}
		base.Io.SetInfoType(base.LogLevelInfo).SystemOutPrintln("服务已优雅关闭")
	}()

	// 7. 启动 HTTP 服务
	base.Io.
		SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Gin 服务启动成功，监听地址: :%d \n", app_obj.App.AppPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		base.Io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("服务启动失败: %v", err)

	}
	return
}

func (r *WebApplication) GetListenPortString() string {
	return ":" + strconv.Itoa(common.GetAppConfig().AppPort)
}

// 工具路由注册（心跳检测、性能分析等）
// 每个系统自动支持 /health 和 /index 访问
func (r *WebApplication) toolRouteRegister(appConfig *app_obj.Application, UrlPrefix string) {
	// 注册默认的公共路由，如健康检查
	r.registerDefaultRoute(UrlPrefix)

	// 是否开启性能分析工具
	r.pProf(appConfig)

	// 注册swagger路由
	r.registerSwagger(appConfig)

}

func (r *WebApplication) registerDefaultRoute(UrlPrefix string) {

	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("#注册健康检查路由...")

	// 注册健康检查请求地址
	r.GinEngine.GET(app_obj.GetHealthPath("health"), func(c *gin.Context) {
		app_obj.SetLastRegisterTime() //设置心跳检测的时间
		c.String(http.StatusOK, "success")
	})

	r.GinEngine.HEAD(app_obj.GetHealthPath("health"), func(c *gin.Context) {
		app_obj.SetLastRegisterTime() //设置心跳检测的时间
		c.String(http.StatusOK, "success")
		return
	})

	// 注册默认路径
	r.GinEngine.GET(app_obj.GetHealthPath("index"), func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, fmt.Sprintf("Welcome \"%s\" Server", UrlPrefix))
	})

	r.GinEngine.HEAD(app_obj.GetHealthPath("index"), func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, fmt.Sprintf("Welcome \"%s\" Server", UrlPrefix))
	})
	return
}

func (r *WebApplication) registerSwagger(appConfig *app_obj.Application) {
	// 如果非线上(release)环境，则可以直接使用
	if app_obj.App.AppEnv == app_obj.EnvProd {
		return
	}

	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("#集成swagger路由...")

	// 文档界面访问URL
	// http://127.0.0.1:8080/swagger/index.html
	r.GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

// 是否开启性能分析工具
func (r *WebApplication) pProf(appConfig *app_obj.Application) {
	if !appConfig.AppNeedPProf {
		return
	}
	// pprof开启后，每隔一段时间(10ms)就会收集当前的堆栈信息，获取各个函数占用的CPU以及内存资源，然后通过对这些采样数据进行分析，形成一个性能分析报告
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("-注册性能分析路由开始-")
	defer func() {
		r.syslog.SetInfoType(base.LogLevelInfo).
			SystemOutPrintln("-注册性能分析路由结束-")
	}()
	pprof.Register(r.GinEngine) // 性能分析用代码

}
