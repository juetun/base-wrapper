package redis_mq

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
)

var (
	RedisDelayMqTopic = []string{}
)

type (
	RedisDelayMq struct {
		Ctx              context.Context     `json:'-'` // 上下文参数 用于停止监听动作
		Context          *base.Context       `json:'-'` // 上下文参数用于记录日志使用
		client           *redis.Client       `json:'-'`
		Topic            string              `json:"topic"`
		PersistentObject PersistentInterface `json:'-'`
		Config           RedisDelayMqConfig  `json:"config"`
		Ticker           *time.Ticker        `json:"-"`
	}
	RedisDelayMqConfig struct {
		Delayer RedisDelayMqConfigDelayer `json:"delayer"`
	}

	RedisDelayMqConfigDelayer struct {
		TimerInterval time.Duration `json:"timer_interval"`
	}

	PersistentInterface interface {
		AddData(data ...int)
	}
	// RedisDelayMqOption RedisMQ赋值属性参数
	RedisDelayMqOption func(mq *RedisDelayMq)

	// ConsumerHandler 消费数据的逻辑
	// Param  string topic 主题
	// Param string msgBody 消息内容
	// Param string  messageId 消息ID
	ConsumerHandler func(topic, msgBody, messageId string) (err error)
)
 //启动数据MQ消费逻辑
func (r *RedisDelayMq) Consumer(topic string, handler ConsumerHandler) {

	ticker := time.NewTicker(time.Duration(r.Config.Delayer.TimerInterval) * time.Millisecond)

 		var bucketName string
		for i := 0; i < config.Setting.BucketSize; i++ {
			timers[i] = time.NewTicker(1 * time.Second)
			bucketName = fmt.Sprintf(config.Setting.BucketName, i+1)
			go waitTicker(timers[i], bucketName)
		}
	 
	
	func waitTicker(timer *time.Ticker, bucketName string) {
		for {
			select {
			case t := <-timer.C:
				tickHandler(t, bucketName)
			}
		}
	}


	go func() {
		for range ticker.C {
			//消费逻辑
			r.consumerRun(topic, handler)
		}
	}()

	r.Ticker = ticker
}

//消费逻辑
func (r *RedisDelayMq) consumerRun(topic string, handler ConsumerHandler) {
	// ZREMRANGEBYSCORE key min max
	r.client.ZRemRangeByScore(context.TODO(), topic, `0`,
		fmt.Sprintf("%d", time.Now().UnixNano()),
	).Args()
}

func NewRedisDelayMq(options ...RedisDelayMqOption) (res *RedisDelayMq) {
	res = &RedisDelayMq{}

	for _, v := range options {
		v(res)
	}

	if res.Config.Delayer.TimerInterval == 0 { //默认一分钟一跳
		res.Config.Delayer.TimerInterval = 1 * time.Second
	}
	return
}

// RedisDelayClient 设置参数方法
func RedisDelayOptionClient(client *redis.Client) RedisDelayMqOption {
	return func(mq *RedisDelayMq) {
		mq.client = client
	}
}

// RedisDelayOptionContext 设置参数方法
func RedisDelayOptionContext(context *base.Context) RedisDelayMqOption {
	return func(mq *RedisDelayMq) {
		mq.Context = context
	}
}

// RedisDelayOptionCtx 设置参数方法
func RedisDelayOptionCtx(Ctx context.Context) RedisDelayMqOption {
	return func(mq *RedisDelayMq) {
		mq.Ctx = Ctx
	}
}
