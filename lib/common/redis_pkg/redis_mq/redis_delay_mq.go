package redis_mq

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
)

const (
	DelayerDefaultLimit = 1000
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
		Limit         int64         `json:"-"` //每次读取条数
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
	DelayMqData     struct {
	}
)

//启动数据MQ消费逻辑
func (r *RedisDelayMq) Consumer(topic string, handler ConsumerHandler) {

	for _, topic := range RedisDelayMqTopic {
		ticker := time.NewTicker(r.Config.Delayer.TimerInterval)
		go r.waitTicker(ticker, topic, handler)
	}
}

func (r *RedisDelayMq) waitTicker(timer *time.Ticker, topic string, tickHandler ConsumerHandler) {
	for {
		select {
		case t := <-timer.C:
			r.tickHandler(t, topic)
		}
	}
	//r.Ticker = ticker
}

func (r *RedisDelayMq) getFromTopic(topic string) (bucketItem []string, err error) {
	err = r.client.ZRangeByScore(r.Ctx, topic, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", time.Now().UnixNano()),
		Offset: 0,
		Count:  r.Config.Delayer.Limit,
	}).ScanSlice(&bucketItem)
	return
}

func (r *RedisDelayMq) getDataWithKeys(topic string, bucketItem []string) (err error) {
	var data interface{}
	for _, item := range bucketItem {
		if err = r.client.HGetAll(r.Ctx, fmt.Sprintf("%s_%s", topic, item)).
			Scan(data); err != nil {
			return
		}
			r.client.zd
	}
	return
}

// 扫描bucket, 取出延迟时间小于当前时间的Job
func (r *RedisDelayMq) tickHandler(t time.Time, topic string) {
	var (
		err        error
		bucketItem []string
	)
	for {
		//从Redis中读取指定条数数据
		if bucketItem, err = r.getFromTopic(topic); err != nil {
			r.Context.Error(map[string]interface{}{
				"desc": fmt.Sprintf("扫描bucket错误#bucket-%s#%s", topic, err.Error()),
			}, "RedisDelayMqTickHandler")
			break
		}

		// 集合为空
		if bucketItem == nil || len(bucketItem) == 0 {
			break
		}

		r.getDataWithKeys(topic, bucketItem)

		//// 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
		//job, err := getJob(bucketItem.jobId)
		//if err != nil {
		//	r.Context.Error(map[string]interface{}{
		//		"desc":  "获取Job元信息失败#bucket-%s#%s",
		//		"topic": topic,
		//		"err":   err.Error(),
		//	}, "")
		//	continue
		//}
		//
		//// job元信息不存在, 从bucket中删除
		//if job == nil {
		//	removeFromBucket(bucketName, bucketItem.jobId)
		//	continue
		//}
		//
		//// 再次确认元信息中delay是否小于等于当前时间
		//if job.Delay > t.Unix() {
		//	// 从bucket中删除旧的jobId
		//	removeFromBucket(bucketName, bucketItem.jobId)
		//	// 重新计算delay时间并放入bucket中
		//	pushToBucket(<-bucketNameChan, job.Delay, bucketItem.jobId)
		//	continue
		//}
		//
		//err = pushToReadyQueue(job.Topic, bucketItem.jobId)
		//if err != nil {
		//	log.Printf("JobId放入ready queue失败#bucket-%s#job-%+v#%s",
		//		bucketName, job, err.Error())
		//	continue
		//}
		//
		//// 从bucket中删除
		//removeFromBucket(bucketName, bucketItem.jobId)
	}
}

//消费逻辑
func (r *RedisDelayMq) consumerRun(topic string, handler ConsumerHandler) {
	// ZREMRANGEBYSCORE key min max
	r.client.ZRemRangeByScore(context.TODO(), topic, `0`,
		fmt.Sprintf("%d", time.Now().UnixNano()),
	).Args()
}

func (r *RedisDelayMq) DefaultOption() {
	if r.Config.Delayer.TimerInterval == 0 { //默认一分钟一跳
		r.Config.Delayer.TimerInterval = 1 * time.Second
	}
	if r.Config.Delayer.Limit == 0 {
		r.Config.Delayer.Limit = DelayerDefaultLimit
	}
	return
}
func NewRedisDelayMq(options ...RedisDelayMqOption) (res *RedisDelayMq) {
	res = &RedisDelayMq{}

	for _, v := range options {
		v(res)
	}

	res.DefaultOption()

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
