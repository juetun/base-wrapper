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
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/middlewares"
)

// 路由注册函数
type HandleRouter func(c *gin.Engine, urlPrefix string)

// 路由函数数组
var HandleFunc = make([]HandleRouter, 0)

type WebApplication struct {
	GinEngine *gin.Engine
	syslog    *common.SystemOut
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
		syslog:    common.NewSystemOut(),
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
			r.syslog.SetInfoType(common.LogLevelError).
				SystemOutPrintf("Load router err  %s", err.Error())
		}
	}()
	appConfig := common.GetAppConfig()
	var UrlPrefix = appConfig.AppName + "/" + appConfig.AppApiVersion
	io := common.
		NewSystemOut().
		SetInfoType(common.LogLevelInfo)

	io.SystemOutPrintf("Start route(app_name:%s) register url config.... ", UrlPrefix)
	defer func() {
		io.SystemOutPrintln("Load route register url finished")
	}()
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("开始注册路由...")

	// 工具路由注册（心跳检测、性能分析等）
	r.toolRouteRegister(appConfig, UrlPrefix)

	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("***********************注册业务路由***********************\n\n")
	// URL路由注册操作
	for _, router := range HandleFunc {
		router(r.GinEngine, UrlPrefix)
	}
	fmt.Printf("\n\n")
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("***********************注册路由结束***********************")

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
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("General start ")

	r.GinEngine.Run(r.getListenPortString()) // listen and serve on 0.0.0.0:8080
	return
}
func (r *WebApplication) start() {
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("Support grace reload")
	srv := &http.Server{
		Addr:    r.getListenPortString(),
		Handler: r.GinEngine,
	}
	r.syslog.SetInfoType(common.LogLevelInfo).SystemOutPrintf("Listen Addr  %s", srv.Addr)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.syslog.SetInfoType(common.LogLevelInfo).SystemOutFatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	r.syslog.SetInfoType(common.LogLevelInfo).SystemOutPrintln("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		r.syslog.SetInfoType(common.LogLevelError).SystemOutFatalf("Server Shutdown:", err)
	}
	r.syslog.SetInfoType(common.LogLevelError).SystemOutPrintln("Server exiting")
	close(quit)
}
func (r *WebApplication) getListenPortString() string {
	return ":" + strconv.Itoa(common.GetAppConfig().AppPort)
}

// 工具路由注册（心跳检测、性能分析等）
func (r *WebApplication) toolRouteRegister(appConfig *common.Application, UrlPrefix string) {
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("-注册健康检查路由-")
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
func (r *WebApplication) pProf(appConfig *common.Application) {
	if !appConfig.AppNeedPProf {
		return
	}
	// pprof开启后，每隔一段时间(10ms)就会收集当前的堆栈信息，获取各个函数占用的CPU以及内存资源，然后通过对这些采样数据进行分析，形成一个性能分析报告
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("-注册性能分析路由-")
	pprof.Register(r.GinEngine) // 性能分析用代码
	r.syslog.SetInfoType(common.LogLevelInfo).
		SystemOutPrintln("-注册性能分析路由-")
}