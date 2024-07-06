package common

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
)

var ExecPath = ""

func GetEnv() string {
	return os.Getenv("GO_ENV")
}

// GetAppConfig 获取当前应用的基本配置
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
func getEnvPath() (env string) {
	if app_obj.App != nil && app_obj.App.AppEnv != "" {
		env = app_obj.App.AppEnv + "/"
	}
	return
}

// GetConfigFileDirectory 获得配置文件所在目录
func GetConfigFileDirectory(notEnv ...bool) (res string) {
	var (
		dir = ExecPath
		err error
		io  = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	)
	env := getEnvPath()

	if app_obj.BaseDirect == "" {

		if ExecPath == "" {
			if dir, err = os.Getwd(); err != nil {
				io.SystemOutPrintf("Template GetConfigFileDirectory is :'%s'", res)
			}
		}

		if len(notEnv) > 0 && notEnv[0] {
			res = fmt.Sprintf("%s/config/", dir)
			return
		}
		res = fmt.Sprintf("%s/config/apps/%s/%s", dir, app_obj.App.AppName, env)
		return

	}

	if len(notEnv) > 0 && notEnv[0] {
		res = fmt.Sprintf("%s/config/", app_obj.BaseDirect)
		return
	}

	res = fmt.Sprintf("%s/config/%s", app_obj.BaseDirect, env)
	return
}

//获取common_config/{env}目录位置
func GetCommonConfigDirectory(notEnv ...bool) (res string) {
	env := getEnvPath()
	res = fmt.Sprintf("%v%v", GetCommonConfigDir(notEnv...), env)
	return
}

//获取common_config目录位置
func GetCommonConfigDir(notEnv ...bool) (res string) {
	res = fmt.Sprintf("%vapps/common_config/", GetConfigFileDirectory(notEnv...))
	return
}

// GetConfigFilePath 获取配置文件的路径
func GetCommonConfigFilePath(fileName string, notEnv ...bool) (res string) {
	dir := GetCommonConfigDirectory(notEnv...)
	res = getConfigFilePathContent(dir, fileName, notEnv...)
	return
}

// GetConfigFilePath 获取配置文件的路径
func GetConfigFilePath(fileName string, notEnv ...bool) (res string) {
	dir := GetConfigFileDirectory(notEnv...)
	res = getConfigFilePathContent(dir, fileName, notEnv...)
	return
}

func getConfigFilePathContent(dir, fileName string, notEnv ...bool) (res string) {
	res = fmt.Sprintf("%s%s", dir, fileName)
	extString := path.Ext(fileName)
	var ext string
	if extString != "" {
		ext = strings.TrimLeft(extString, ".")
	}
	io := base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintf("dir:%v,fileName:%v,ext:%v,config_file_path:%v", dir, fileName, ext, res)
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

// NewApplication 初始化应用信息
func NewApplication() *app_obj.Application {
	var env = GetEnv()
	if env == "" { // 默认环境为线上环境
		env = app_obj.EnvProd
	}
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintf("Env:'%s'(You can set environment variable with 'export \"GO_ENV=%s\")", env, strings.Join(app_obj.EnvList, "|"))
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
