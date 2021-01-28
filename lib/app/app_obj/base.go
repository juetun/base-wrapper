/**
* @Author:changjiang
* @Description:
* @File:base
* @Version: 1.0.0
* @Date 2020/10/15 12:12 上午
 */
package app_obj

import (
	"encoding/json"
)

const (
	TRACE_ID        = "trace_id"     // 请求上下文传递时 的唯一ID （日志链路追踪使用）
	APP_LOG_KEY     = "app"          // 当前操作所属的代码库微服务应用名（日志链路追踪使用）
	APP_FIELD_KEY   = "type"         // 日志类型字段KEY的值
	APP_LOG_LOC     = "src"          // 代码所在位置
	HTTP_TRACE_ID   = "X-Trace-Id"   // 页面请求时的 传参或者nginx生成的trace_id的key
	HTTP_USER_TOKEN = "X-Auth-Token" // 页面请求时用户token

)

// 当前配置文件所在目录
var BaseDirect string

var App *Application

// 应用基本的配置结构体
type Application struct {
	CronPort             int    `json:"cron_port" yaml:"cron_port"`                       // 定时任务客户端端口
	AppAlias             string `json:"app_alias" yaml:"alias"`                           // 服务器别名
	AppSystemName        string `json:"app_system_name" yaml:"system_name"`               // 系统名称
	AppEnv               string `json:"app_env" yaml:"env"`                               // 当前运行环境
	AppName              string `json:"app_name" yaml:"name"`                             // 应用名称
	AppVersion           string `json:"app_version" yaml:"version"`                       // 应用版本以前缀v 开头
	AppApiVersion        string `json:"app_api_version" yaml:"app_api_version"`           // 应用的API的版本号，用于api接口路由参数拼接
	AppPort              int    `json:"app_port" yaml:"port"`                             // 应用监听的端口
	AppGraceReload       int    `json:"grace_reload" yaml:"grace_reload"`                 // 应用是否支持优雅重启
	AppNeedPProf         bool   `json:"app_need_p_prof" yaml:"app_need_p_prof"`           // 是否需要内存分析
	AppTemplateDirectory string `json:"app_template_directory" yaml:"template_directory"` // temp模板默认目录
}

func (r *Application) ToString() string {
	v, _ := json.Marshal(r)
	return string(v)
}

func (r *Application) Default() {
	if r.AppPort == 0 { // 默认80端口
		r.AppPort = 80
	}
	if r.AppEnv != "" && r.AppEnv == "" {
		r.AppEnv = App.AppEnv
	}
	if r.CronPort == 0 {
		r.CronPort = 5921
	}
	r.AppVersion = "v" + r.AppVersion

}
