package app_start

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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

func NewWebApplication(privateMiddleWares ...gin.HandlerFunc) *WebApplication {
	if false && common.GetAppConfig().AppEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	webApp := &WebApplication{
		GinEngine: gin.Default(),
		syslog:    common.NewSystemOut(),
	}
	// 加载GIN框架 中间件
	middlewares.LoadMiddleWare(privateMiddleWares...)

	webApp.GinEngine.Use(middlewares.MiddleWareComponent...)
	return webApp
}
func RunLoadRouter(c *gin.Engine) (err error) {
	appConfig := common.GetAppConfig()
	var UrlPrefix = appConfig.AppName + "/" + appConfig.AppApiVersion
	io := common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	io.SystemOutPrintf("Start load route config.... %s", UrlPrefix)
	defer func() {
		io.SystemOutPrintln("Load route finished")
	}()
	for _, router := range HandleFunc {
		router(c, UrlPrefix)
	}

	return
}
func (r *WebApplication) LoadRouter() *WebApplication {

	err := RunLoadRouter(r.GinEngine) // 注册Gin路由组件
	if err != nil {
		r.syslog.SetInfoType(common.LogLevelError).
			SystemOutPrintf("Load router err  %s", err.Error())
	}
	return r
}

// 开始加载Gin 服务
func (r *WebApplication) Run() (err error) {
	appConfig := common.GetAppConfig()
	if appConfig.AppNeedPProf { // pprof开启后，每隔一段时间(10ms)就会收集当前的堆栈信息，获取各个函数占用的CPU以及内存资源，然后通过对这些采样数据进行分析，形成一个性能分析报告
		r.syslog.SetInfoType(common.LogLevelInfo).
			SystemOutPrintln("开启性能分析 ")
		pprof.Register(r.GinEngine) // 性能分析用代码
	}
	defaultEngine(r.GinEngine)

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

func defaultEngine(r *gin.Engine) {
	r.GET("/index", func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
}
