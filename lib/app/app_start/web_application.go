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
	"time"
)

// HandleRouter 路由注册函数
type HandleRouter func(c *gin.Engine, urlPrefix string)
type RouterPath func(httpMethod, absolutePath, handlerName string, nuHandlers int)

var HandleFunc = make([]HandleRouter, 0)         // 路由函数切片
var HandleFuncIntranet = make([]HandleRouter, 0) // 内网路由函数切片
var HandleFuncOuterNet = make([]HandleRouter, 0) // 外网路由函数切片
var HandleFuncAdminNet = make([]HandleRouter, 0) // 外网路由函数切片
var HandleFuncPage = make([]HandleRouter, 0)     // 外网路由函数切片
var RoutePathInitCallBack RouterPath             //注册路由时调用
type WebApplication struct {
	GinEngine *gin.Engine
	syslog    *base.SystemOut
	//FlagMicro bool // 如果是支持微服务
	MicroOperate MicroOperateInterface
}

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

func (r *WebApplication) SetFlagMicro(micro MicroOperateInterface) (res *WebApplication) {
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
			SystemOutPrintln("注册外网访问接口路由....")
		for _, router := range HandleFuncAdminNet {
			r.GinEngine.Use(middlewares.AdminMiddlewares())
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

// Run 开始加载Gin 服务
func (r *WebApplication) Run(cTxs ...context.Context) (err error) {
	var ctx context.Context
	if len(cTxs) > 0 {
		ctx = cTxs[0]
	}
	if ctx == nil {
		ctx = context.Background()
	}
	appConfig := common.GetAppConfig()

	// // 如果支持优雅重启
	if appConfig.AppGraceReload > 0 {
		r.start(ctx)
		return
	}
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintln("General start ")
	if r.MicroOperate != nil { //如果实现了微服务注册与发现
		r.MicroOperate.RegisterMicro(r.GinEngine, ctx)
	}

	return
}

//将服务从注册中心拿掉
func (r *WebApplication) UnRegisterMicro() {
	if r.MicroOperate != nil {
		r.MicroOperate.UnRegisterMicro()
	}
	return
}

func (r *WebApplication) start(ctx context.Context) {

	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("Support grace reload")

	httpServer := &http.Server{
		Addr:    r.GetListenPortString(),
		Handler: r.GinEngine,
	}
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("Listen Addr  %s", httpServer.Addr)

	go func() { // 启动GIN服务动作

		if r.MicroOperate != nil {
			var err error
			_, err = r.MicroOperate.RegisterMicro(r.GinEngine, ctx)
			if err != nil {
				r.syslog.SetInfoType(base.LogLevelError).SystemOutFatalf("listen: %s\n", err.Error())
				return
			}
			return
		}

		// service connections
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.syslog.SetInfoType(base.LogLevelInfo).SystemOutFatalf("listen: %s\n", err)
		}

	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	defer func() {
		close(quit)
	}()
	signal.Notify(quit, os.Interrupt)
	<-quit
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
		r.syslog.SetInfoType(base.LogLevelError).SystemOutPrintln("Server exiting")
	}()

	if r.MicroOperate != nil { //如果开启了微服务,
		r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintln("将服务信息从注册中心移除")
		//则先将服务从注册中心拿掉
		r.UnRegisterMicro()
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		r.syslog.SetInfoType(base.LogLevelError).SystemOutFatalf("Server Shutdown:", err)
	}

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

func (r *WebApplication) getHealthPath(suffix string) (pathString string) {
	if app_obj.RouteTypeDefaultIntranet == "" {
		pathString = fmt.Sprintf("/%v/%v", app_obj.App.AppName, suffix)
		return
	}
	pathString = fmt.Sprintf("/%v/%v/%v", app_obj.App.AppName, app_obj.RouteTypeDefaultIntranet, suffix)
	return
}

func (r *WebApplication) registerDefaultRoute(UrlPrefix string) {
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("#注册健康检查路由...")
	// 注册健康检查请求地址
	r.GinEngine.GET(r.getHealthPath("health"), func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	r.GinEngine.HEAD(r.getHealthPath("health"), func(c *gin.Context) {
		c.String(http.StatusOK, "success")
		return
	})

	// 注册默认路径
	r.GinEngine.GET(r.getHealthPath("index"), func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, fmt.Sprintf("Welcome \"%s\" Server", UrlPrefix))
	})
	r.GinEngine.HEAD(r.getHealthPath("index"), func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, fmt.Sprintf("Welcome \"%s\" Server", UrlPrefix))
	})
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
