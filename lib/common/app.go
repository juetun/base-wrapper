package common

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/spf13/viper"
)

const (
	// 开发环境
	ENV_DEVELOP = "dev"

	// 测试环境
	ENV_TEST = "test"

	// demo
	ENV_DEMO = "demo"

	// 线上环境
	ENV_RELEASE = "release"
)

var app *Application

// 应用基本的配置结构体
type Application struct {
	AppSystemName  string `json:"app_system_name"`
	AppEnv         string `json:"app_env"`         // 当前运行环境
	AppName        string `json:"app_name"`        // 应用名称
	AppVersion     string `json:"app_version"`     // 应用版本以前缀v 开头
	AppApiVersion  string `json:"app_api_version"` // 应用的API的版本号，用于api接口路由参数拼接
	AppPort        int    `json:"app_port"`        // 应用监听的端口
	AppGraceReload bool   `json:"grace_reload"`    // 应用是否支持优雅重启
	AppNeedPProf   bool   `json:"app_need_p_prof"` // 是否需要内存分析
}

func (r *Application) ToString() string {
	v, _ := json.Marshal(r)
	return string(v)
}

func GetEnv() string {
	return os.Getenv("GO_ENV")
}

// 获取当前应用的基本配置
func GetAppConfig() *Application {
	return app
}

func PluginsApp() (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	app = NewApplication()
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintln("Load app config start")
	defer func() {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("app config is: '%v' \n", app.ToString())
		io.SystemOutPrintf("load app config finished \n")
	}()
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	dir := GetConfigFileDirectory()
	io.SystemOutPrintf("config directory is : '%s' ", dir)
	viper.AddConfigPath(dir + "/../") // path to look for the config file in
	err = viper.ReadInConfig()        // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) { // 热加载
		fmt.Println("Config file changed:", e.Name)
	})

	version := viper.GetString("app.version")
	app.AppVersion = "v" + version
	app.AppPort = viper.GetInt("app.port")
	if app.AppPort == 0 { // 默认80端口
		app.AppPort = 80
	}
	app.AppName = viper.GetString("app.name")
	app.AppGraceReload = viper.GetBool("app.grace_reload")
	app.AppSystemName = viper.GetString("app.system_name")
	app.AppNeedPProf = viper.GetBool("app.app_need_p_prof")
	return
}

// 获得配置文件所在目录
func GetConfigFileDirectory() string {
	var env = ""
	if app.AppEnv != "" {
		env = app.AppEnv + "/"
	}
	return "./config/" + env
}

// 初始化应用信息
func NewApplication() *Application {
	var env = GetEnv()
	if env == "" { // 默认环境为线上环境
		env = ENV_RELEASE
	}
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintf("Env is: '%s'", env)
	return &Application{
		AppSystemName:  "",
		AppName:        "app",
		AppVersion:     "v1.0",
		AppApiVersion:  "v1",
		AppPort:        8080,
		AppEnv:         env,
		AppGraceReload: false,
		AppNeedPProf:   false,
	}
}
