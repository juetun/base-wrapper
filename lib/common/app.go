package common

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
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

func GetEnv() string {
	return os.Getenv("GO_ENV")
}

// 获取当前应用的基本配置
func GetAppConfig() *app_obj.Application {
	return app_obj.App
}

func PluginsApp() (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	app_obj.App = NewApplication()
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintln("Load app config start")
	defer func() {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("app config is: '%v' \n", app_obj.App.ToString())
		io.SystemOutPrintf("load app config finished \n")
	}()
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	dir := GetConfigFileDirectory()
	io.SystemOutPrintf("config directory is : '%s' ", dir)
	viper.AddConfigPath(dir + "/../") // path to look for the config file in
	err = viper.ReadInConfig()        // Find and read the config file
	if err != nil { // Handle errors reading the config file
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) { // 热加载
		io.SystemOutPrintf("Config file changed:", e.Name)
	})

	version := viper.GetString("app.version")
	app_obj.App.AppVersion = "v" + version
	app_obj.App.AppPort = viper.GetInt("app.port")
	if app_obj.App.AppPort == 0 { // 默认80端口
		app_obj.App.AppPort = 80
	}
	app_obj.App.AppName = viper.GetString("app.name")
	app_obj.App.AppGraceReload = viper.GetBool("app.grace_reload")
	app_obj.App.AppSystemName = viper.GetString("app.system_name")
	app_obj.App.AppNeedPProf = viper.GetBool("app.app_need_p_prof")
	app_obj.App.AppTemplateDirectory = defaultAppTemplateDirectory(io, viper.GetString("app.app_template_directory"))
	return
}
func defaultAppTemplateDirectory(io *base.SystemOut, dir string) (res string) {
	if dir != "" {
		res = dir
	}
	var err error
	if dir, err = os.Getwd(); err != nil {
		return
	}
	res = fmt.Sprintf("%s/web/views/", dir)
	io.SystemOutPrintf("Template default directory is :'%s'", res)
	return
}

// 获得配置文件所在目录
func GetConfigFileDirectory() string {
	var env = ""
	if app_obj.App.AppEnv != "" {
		env = app_obj.App.AppEnv + "/"
	}
	if app_obj.BaseDirect == "" {
		return "./config/" + env
	}
	return fmt.Sprintf("%s/config/%s", app_obj.BaseDirect, env)
}

// 初始化应用信息
func NewApplication() *app_obj.Application {
	var env = GetEnv()
	if env == "" { // 默认环境为线上环境
		env = ENV_RELEASE
	}
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintf("Env:'%s'(You can set environment variable with 'export \"GO_ENV=%s\")", env, strings.Join([]string{ENV_DEVELOP, ENV_TEST, ENV_TEST, ENV_DEMO, ENV_RELEASE,}, "|"), )
	return &app_obj.Application{
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
