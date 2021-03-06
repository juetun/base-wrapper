package common

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
	"os"
	"path"
	"strings"
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

func DefaultAppTemplateDirectory(io *base.SystemOut) (res string) {
	var dir string
	var err error
	if dir, err = os.Getwd(); err != nil {
		return
	}
	res = fmt.Sprintf("%s/web/views/", dir)
	io.SystemOutPrintf("Template default directory is :'%s'", res)
	return
}

// 获得配置文件所在目录
func GetConfigFileDirectory(notEnv ...bool) (res string) {

	var env = ""
	if app_obj.App.AppEnv != "" {
		env = app_obj.App.AppEnv + "/"
	}
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

	if app_obj.BaseDirect == "" {
		var (
			dir string
			err error
		)
		if dir, err = os.Getwd(); err != nil {
			io.SystemOutPrintf("Template GetConfigFileDirectory is :'%s'", res)
		}
		if len(notEnv) > 0 && notEnv[0] {
			return fmt.Sprintf("%s/config/", dir)
		} else {
			return fmt.Sprintf("%s/config/%s", dir, env)
		}

	}
	if len(notEnv) > 0 && notEnv[0] {
		res = fmt.Sprintf("%s/config/", app_obj.BaseDirect)
		return
	} else {
		return fmt.Sprintf("%s/config/%s", app_obj.BaseDirect, env)
	}

}

// 获取配置文件的路径
func GetConfigFilePath(fileName string, notEnv ...bool) (res string) {
	dir := GetConfigFileDirectory(notEnv...)
	res = fmt.Sprintf("%s%s", dir, fileName)

	extString := path.Ext(fileName)
	var ext string
	if extString != "" {
		ext = strings.TrimLeft(extString, ".")
	}
	switch ext {
	case "yaml":
		if ok, _ := utils.PathExists(res); ok {
			return
		}
		res = fmt.Sprintf("%s%s.yml", dir, strings.TrimSuffix(path.Base(fileName), extString))
		return
	case "yml":
		if ok, _ := utils.PathExists(res); ok {
			return
		}
		res = fmt.Sprintf("%s%s.yaml", dir, strings.TrimSuffix(path.Base(fileName), extString))
		return
	}

	return
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
		AppGraceReload: 0,
		AppNeedPProf:   false,
	}
}
