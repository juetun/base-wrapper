// Package app_obj
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
	DebugFlag      = "req_debug"  // 当前请求是否为调试
	TraceId        = "trace_id"   // 请求上下文传递时 的唯一ID （日志链路追踪使用）
	AppLogKey      = "app"        // 当前操作所属的代码库微服务应用名（日志链路追踪使用）
	AppFieldKey    = "type"       // 日志类型字段KEY的值
	AppLogLoc      = "src"        // 代码所在位置
	HttpTraceId    = "X-Trace-Id" // 页面请求时的 传参或者nginx生成的trace_id的key
	HttpUserHid    = "X-User-Hid" // 页面请求时的 用户ID
	HttpTimestamp  = "X-Timestamp"
	HttpSign       = "X-Sign"
	HttpShopId     = "X-Shop-Id" // 店铺ID
	HttpHeaderInfo = "X-Info"    // 用户其他信息预埋参数

	HttpUserToken = "X-Auth-Token" // 页面请求时用户token

	//HttpHeaderTerminal        = "X-Terminal"    // 终端类型 android ,ios ,web weixin
	//HttpHeaderChannel         = "X-Channel"     // 渠道
	HttpHeaderApp             = "X-App"         // 接口请求头信息
	HttpHeaderVersion         = "X-App-version" // 接口请求版本信息
	HttpHeaderAdminToken      = "X-Console"     // 客服后台接口多的key值
	HttpResponseAdministrator = "X-Admin"       // 程序负责人、管理员

	WebSocketKey      = "Sec-Websocket-Key"
	WebSocketHeaderIp = "X-Forwarded-For" // 取Ip地址方法
	WebSocketUid      = "uid"

	DbNameKey         = "dbName"
	DbContextValueKey = "DbContextValue" // 数据库操作上下文传参保存的KEY
)
const (
	EnvDev  = "dev"  // 开发环境
	EnvTest = "test" // 测试环境
	EnvPre  = "pre"  // 预发布环境
	EnvProd = "prod" // 线上环境
)

var (
	EnvList = []string{
		EnvProd, EnvDev, EnvTest, EnvPre,
	}

	// BaseDirect 当前配置文件所在目录
	BaseDirect string

	App *Application

	RouteTypeDefaultIntranet = "in"
	RouteTypeDefaultOutranet = "out"
	RouteTypeDefaultAdminNet = "admin"
	RouteTypeDefaultPage     = ""
)

// Application 应用基本的配置结构体
type (
	Application struct {
		CronPort             int             `json:"cron_port" yaml:"cron_port"`                       // 定时任务客户端端口
		AppAlias             string          `json:"app_alias" yaml:"alias"`                           // 服务器别名
		AppSystemName        string          `json:"app_system_name" yaml:"system_name"`               // 系统名称
		AppEnv               string          `json:"app_env" yaml:"env"`                               // 当前运行环境
		AppName              string          `json:"app_name" yaml:"name"`                             // 应用名称
		AppVersion           string          `json:"app_version" yaml:"version"`                       // 应用版本以前缀v 开头
		AppApiVersion        string          `json:"app_api_version" yaml:"app_api_version"`           // 应用的API的版本号，用于api接口路由参数拼接
		AppPort              int             `json:"app_port" yaml:"port"`                             // 应用监听的端口
		AppGraceReload       int             `json:"grace_reload" yaml:"grace_reload"`                 // 应用是否支持优雅重启
		AppNeedPProf         bool            `json:"app_need_p_prof" yaml:"app_need_p_prof"`           // 是否需要内存分析
		AppSignDebug         bool            `json:"app_sign_debug" yaml:"app_sign_debug"`             // 是否支持签名调试
		AppTemplateDirectory string          `json:"app_template_directory" yaml:"template_directory"` // temp模板默认目录
		AppRouterPrefix      AppRouterPrefix `json:"app_router_prefix" yaml:"app_router_prefix"`       // 路由前缀
		NotValidateSign      []UrlFormat     `json:"not_validate_sign"  yaml:"not_validate_sign"`      //不需要签名验证的Uri
		NotSendHeader        []UrlFormat     `json:"not_send_header"  yaml:"not_send_header"`          //不需要设置headerInfo的Uri
		AppAdminToken        string          `json:"app_admin_token" yaml:"app_admin_token"`           // 客服后台接口多的token值
		Administrator        string          `json:"administrator" yaml:"administrator"`               //代码管理员信息
		AppRunTimerTask      bool            `json:"app_run_timer_task" yaml:"app_run_timer_task"`     //是否支持定时任务
		UseDefaultShopId     bool            `json:"use_default_shop_id" yaml:"use_default_shop_id"`   //测试环境调试数据使用的默认店铺ID
	}
	UrlFormat struct {
		IsPrefix bool            `json:"is_prefix" yaml:"is_prefix"`
		Uri      string          `json:"uri" yaml:"uri"`
		Method   map[string]bool `json:"method" yaml:"method"`
	}
	AppRouterPrefix struct {
		Intranet string `json:"intranet"` // 内网地址
		Outranet string `json:"outernet"` // 外网地址
		AdminNet string `json:"adminnet"` // 运营后台地址
		Page     string `json:"page"`     // 网页地址
	}
)

func (r *Application) ToString() string {
	v, _ := json.Marshal(r)
	return string(v)
}

func (r *Application) ToMap() (res map[string]interface{}) {
	s := r.ToString()
	_ = json.Unmarshal([]byte(s), &res)
	return
}
func (r *Application) Default() {
	if r.AppPort == 0 { // 默认80端口
		r.AppPort = 80
	}
	if r.AppEnv == "" {
		r.AppEnv = App.AppEnv
	}
	if r.CronPort == 0 {
		r.CronPort = 5921
	}
	r.AppVersion = "v" + r.AppVersion
	if r.AppRouterPrefix.Intranet == "" {
		r.AppRouterPrefix.Intranet = RouteTypeDefaultIntranet
	}
	if r.AppRouterPrefix.Outranet == "" {
		r.AppRouterPrefix.Outranet = RouteTypeDefaultOutranet
	}
	if r.AppRouterPrefix.AdminNet == "" {
		r.AppRouterPrefix.AdminNet = RouteTypeDefaultAdminNet
	}
	if r.AppRouterPrefix.Page == "" {
		r.AppRouterPrefix.Page = RouteTypeDefaultPage
	}
}
