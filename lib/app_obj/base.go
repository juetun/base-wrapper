/**
* @Author:changjiang
* @Description:
* @File:base
* @Version: 1.0.0
* @Date 2020/10/15 12:12 上午
 */
package app_obj

const (
	TRACE_ID        = "trace_id"     // 请求上下文传递时 的唯一ID （日志链路追踪使用）
	APP_LOG_KEY     = "app"          // 当前操作所属的代码库微服务应用名（日志链路追踪使用）
	HTTP_TRACE_ID   = "X-Trace-Id"   // 页面请求时的 传参或者nginx生成的trace_id的key
	HTTP_USER_TOKEN = "X-Auth-Token" // 页面请求时用户token
)

// 当前配置文件所在目录
var BaseDirect string
