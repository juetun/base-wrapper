package app_start

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/middlewares"
)

// 路由注册函数
type HandleRouter func(c *gin.Engine, urlPrefix string)

// 路由函数数组
var HandleFunc = make([]HandleRouter, 0)

type WebApplication struct {
	GinEngine *gin.Engine
	syslog    *base.SystemOut
}

// privateMiddleWares 项目自定义的GIN中间件
func NewWebApplication(privateMiddleWares ...gin.HandlerFunc) *WebApplication {
	switch strings.ToLower(common.GetAppConfig().AppEnv) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	webApp := &WebApplication{
		GinEngine: gin.New(),
		syslog:    base.NewSystemOut(),
	}

	// 加载GIN框架 中间件
	middlewares.LoadMiddleWare(privateMiddleWares...)

	// gin加载中间件
	webApp.GinEngine.Use(middlewares.MiddleWareComponent...)

	logger := app_log.GetLog()

	// 日志对象获取
	webApp.GinEngine.Use(middlewares.GinLogCollect(logger))

	return webApp
}

// 加载API路由
func (r *WebApplication) LoadRouter() *WebApplication {
	var err error
	defer func() {
		if err != nil {
			r.syslog.SetInfoType(base.LogLevelError).
				SystemOutPrintf("Load router err  %s", err.Error())
		}
	}()
	appConfig := common.GetAppConfig()
	var UrlPrefix = appConfig.AppName + "/" + appConfig.AppApiVersion

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

	// URL路由注册操作
	for _, router := range HandleFunc {
		router(r.GinEngine, UrlPrefix)
	}

	return r
}

// 开始加载Gin 服务
func (r *WebApplication) Run() (err error) {
	appConfig := common.GetAppConfig()

	// // 如果支持优雅重启
	if appConfig.AppGraceReload {
		r.start()
		return
	}
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("General start ")

	r.GinEngine.Run(r.getListenPortString()) // listen and serve on 0.0.0.0:8080
	return
}
func (r *WebApplication) start() {
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("Support grace reload")
	srv := &http.Server{
		Addr:    r.getListenPortString(),
		Handler: r.GinEngine,
	}
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintf("Listen Addr  %s", srv.Addr)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.syslog.SetInfoType(base.LogLevelInfo).SystemOutFatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	r.syslog.SetInfoType(base.LogLevelInfo).SystemOutPrintln("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		r.syslog.SetInfoType(base.LogLevelError).SystemOutFatalf("Server Shutdown:", err)
	}
	r.syslog.SetInfoType(base.LogLevelError).SystemOutPrintln("Server exiting")
	close(quit)
}
func (r *WebApplication) getListenPortString() string {
	return ":" + strconv.Itoa(common.GetAppConfig().AppPort)
}

// 工具路由注册（心跳检测、性能分析等）
func (r *WebApplication) toolRouteRegister(appConfig *app_obj.Application, UrlPrefix string) {
	r.syslog.SetInfoType(base.LogLevelInfo).
		SystemOutPrintln("1、注册健康检查路由...")
	// 注册健康检查请求地址
	r.GinEngine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})
	// 注册默认路径
	r.GinEngine.GET("/index", func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, fmt.Sprintf("Welcome \"%s\" Server", UrlPrefix))
	})
	// 是否开启性能分析工具
	r.pProf(appConfig)
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
