/**
* @Author:changjiang
* @Description:
* @File:Context
* @Version: 1.0.0
* @Date 2020/3/28 10:18 上午
 */
package base

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/juetun/base-wrapper/lib/app_obj"
)

func GetControllerBaseContext(controller *ControllerBase, c *gin.Context) (res *Context) {
	return NewContext(Log(controller.Log), GinContext(c))
}
func NewContext(contextOption ...ContextOption) *Context {
	context := &Context{}

	// 初始化默认值
	context.Init()

	// 为数据指定初始化值
	for _, option := range contextOption {
		option(context)
	}
	return context
}

type Context struct {
	Log         *app_log.AppLog `json:"log"`
	Db          *gorm.DB        `json:"db"`
	CacheClient *redis.Client   `json:"cache_client"`
	GinContext  *gin.Context
}
type ContextOption func(context *Context)

func Log(opt *app_log.AppLog) ContextOption {
	return func(context *Context) {
		context.Log = opt
	}
}
func Db(opt *gorm.DB) ContextOption {
	return func(context *Context) {
		context.Db = opt
	}
}
func CacheClient(opt *redis.Client) ContextOption {
	return func(context *Context) {
		context.CacheClient = opt
	}
}
func GinContext(opt *gin.Context) ContextOption {
	return func(context *Context) {
		context.GinContext = opt
	}
}
func (r *Context) Init() (c *Context) {
	if r.Log == nil {
		r.Log = app_log.GetLog()
	}
	if r.Db == nil {
		r.Db = app_obj.GetDbClient()
		s := ""
		if nil == r.GinContext {
			var io = NewSystemOut().SetInfoType(LogLevelInfo)
			io.SetInfoType("WARN").SystemOutFatalf("您没有设置上下文 gin.context的值，将无法记录日志trace_id")
		} else {
			if tp, ok := r.GinContext.Get(app_obj.TRACE_ID); ok {
				s = fmt.Sprintf("%v", tp)
			}
		}
		r.Db.InstantSet(app_obj.TRACE_ID, s)
	}
	if r.CacheClient == nil {
		r.CacheClient = app_obj.GetRedisClient()
	}
	return r
}
