package redis_mq

import (
	"context"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	. "github.com/juetun/base-wrapper/lib/plugins" // 组件目录
	"log"
	"os"
	"testing"
	"time"
)

type TestBaseController struct {
	base.ControllerBase
}

func newTestBaseController() (res *TestBaseController) {
	return &TestBaseController{}
}
func TestNewRedisDelayMq(t *testing.T) {
	initReady()
	con := newTestBaseController()
	ctxp, _ := base.CreateCrontabContext(con.ControllerBase)

	type args struct {
		options []RedisDelayMqOption
	}
	tests := []struct {
		name    string
		args    args
		wantRes *RedisDelayMq
	}{
		{
			args: args{options: []RedisDelayMqOption{
				RedisDelayOptionClient(ctxp.CacheClient),
				RedisDelayOptionContext(ctxp),
				RedisDelayOptionCtx(context.TODO()),
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := NewRedisDelayMq(tt.args.options...)
			for i := 0; i < 100; i++ {
				gotRes.Add("topic_delay_mq", DelayMqData{
					Timestamp: time.Now(),
					Data:      fmt.Sprintf("adfasdf_%d", i),
				})
			}

			//gotRes.Consumer("", "", consumerData)
			//if gotRes := NewRedisDelayMq(tt.args.options...); !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("NewRedisDelayMq() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}
func initReady() {
	dir, _ := os.Getwd()
	common.ExecPath = fmt.Sprintf("%s/../../../..", dir)
	app_start.NewPlugins().Use(
		PluginRegistry,
		PluginClickHouse,
		PluginOss,
		PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		// PluginElasticSearchV7,
		PluginShortMessage,
		PluginAppMap,
		//PluginAuthorization,
		// func(arg *app_start.PluginsOperate) (err error) {
		// 	// 启动websocket
		// 	go anvil_websocket.WebsocketStart()
		// 	return
		// },
		// plugins.PluginOss,
	).LoadPlugins() // 加载插件动作

}
func TestNewRedisDelayMqConsumer(t *testing.T) {
	initReady()
	con := newTestBaseController()
	ctxp, _ := base.CreateCrontabContext(con.ControllerBase)

	type args struct {
		options []RedisDelayMqOption
	}
	tests := []struct {
		name    string
		args    args
		wantRes *RedisDelayMq
	}{
		{
			args: args{options: []RedisDelayMqOption{
				RedisDelayOptionClient(ctxp.CacheClient),
				RedisDelayOptionContext(ctxp),
				RedisDelayOptionCtx(context.TODO()),
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := NewRedisDelayMq(tt.args.options...)
			gotRes.Consumer("topic_delay_mq", "gid", func(topic, msgBody, messageId string) (err error) {
				log.Printf("消费日志:topic:%s msgBody:%s messageId:%s \n", topic, msgBody, messageId)
				return
			})

			masterRuntime()
			//gotRes.Consumer("", "", consumerData)
			//if gotRes := NewRedisDelayMq(tt.args.options...); !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("NewRedisDelayMq() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}
func masterRuntime() {
	var i int
	for {
		if i > 1000 {
			break
		}
		log.Printf("延迟处理 \n")
		time.Sleep(1 * time.Second)
		i++
	}
}
