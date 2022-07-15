package redis_mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/common/redis_pkg/anvil_redis"
	"github.com/juetun/base-wrapper/lib/utils"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
)

const (
	DelayerDefaultLimit = 1000
)

//初始化消息队列方法
func NewRedisDelayMq(options ...RedisDelayMqOption) (res *RedisDelayMq) {
	res = &RedisDelayMq{}
	for _, v := range options {
		v(res)
	}

	res.DefaultOption()
	return
}

type (
	MqConsumerItem struct {
		Topic string      `json:"topic"`
		Timer *time.Timer `json:"-"`
	}
	RedisDelayMq struct {
		Ctx              context.Context     `json:"-"` // 上下文参数 用于停止监听动作
		Context          *base.Context       `json:"-"` // 上下文参数用于记录日志使用
		client           *redis.Client       `json:"-"`
		Topic            string              `json:"topic"`
		PersistentObject PersistentInterface `json:"-"`
		Config           RedisDelayMqConfig  `json:"config"`
		Ticker           *time.Ticker        `json:"-"`

		ConsumerMapGroupIdMessageHandler map[string]ConsumerHandler
	}
	RedisDelayMqConfig struct {
		Delayer RedisDelayMqConfigDelayer `json:"delayer"`
	}

	RedisDelayMqConfigDelayer struct {
		TimerInterval time.Duration `json:"timer_interval"`
		Limit         int64         `json:"-"` //每次读取条数
	}

	PersistentInterface interface {
		//将数据添加入队列
		AddData(data ...DelayMqData) (err error)

		//将数据移除
		RemoveData(dataKey ...string) (err error)
	}
	// RedisDelayMqOption RedisMQ赋值属性参数
	RedisDelayMqOption func(mq *RedisDelayMq)

	// ConsumerHandler 消费数据的逻辑
	// Param  string topic 主题
	// Param string msgBody 消息内容
	// Param string  messageId 消息ID
	ConsumerHandler func(topic, msgBody, messageId string) (err error)
	DelayMqData     struct {
		Timestamp time.Time   `json:"timestamp"`
		Data      interface{} `json:"data"`
	}
)

func (d DelayMqData) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(d)
	return
}

func (r *RedisDelayMq) initConsumerMapGroupIdMessageHandler() {
	if r.ConsumerMapGroupIdMessageHandler == nil {
		r.ConsumerMapGroupIdMessageHandler = map[string]ConsumerHandler{}
	}
}

//启动数据MQ消费逻辑
func (r *RedisDelayMq) Consumer(topic string, handler ConsumerHandler) {
	r.initConsumerMapGroupIdMessageHandler()
	var RedisDelayMqTopic = []MqConsumerItem{
		{
			Topic: topic,
		},
	}
	for _, item := range RedisDelayMqTopic {
		item.Timer = time.NewTimer(r.Config.Delayer.TimerInterval)
		go r.waitTimer(item, handler)
	}

}

func (r *RedisDelayMq) deferWaitTicker(t time.Time, mqConsumerItem MqConsumerItem, tickHandler ConsumerHandler) {
	log.Printf("tick触发 \n")

	defer func() {
		mqConsumerItem.Timer.Reset(r.Config.Delayer.TimerInterval)
	}()
	uk := fmt.Sprintf("rdslk:%s", mqConsumerItem.Topic)

	_ = anvil_redis.NewRedisLock(
		anvil_redis.RedisLockDuration(10*time.Second),
		anvil_redis.RedisLockContext(r.Context),
		anvil_redis.RedisLockCtx(r.Ctx),
		anvil_redis.RedisLockAttemptsTime(100), //尝试获取锁的次数
		anvil_redis.RedisLockAttemptsInterval(30*time.Millisecond),
		anvil_redis.RedisLockUniqueKey(utils.Guid(uk)),
		anvil_redis.RedisLockLockKey(uk),
		anvil_redis.RedisLockOkHandler(func(ctx context.Context) (err error) {
			r.tickHandler(t, mqConsumerItem.Topic, tickHandler)
			return
		}),
	).RunWithGetLock()
}

func (r *RedisDelayMq) waitTimer(mqConsumerItem MqConsumerItem, tickHandler ConsumerHandler) {
	for {
		select {
		case t := <-mqConsumerItem.Timer.C:
			r.deferWaitTicker(t, mqConsumerItem, tickHandler)
		default:
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
	}
	return
}

//处理数据
func (r *RedisDelayMq) dealData(tickHandler ConsumerHandler, topic string, bucketItem ...string) (err error) {

	for _, data := range bucketItem {
		if err = r.getCacheListByKey(data); err != nil {
			return
		}
	}
	return
}

func (r *RedisDelayMq) getCacheListByKey(key string) (err error) {

	return
}

func (r *RedisDelayMq) readData(topic string, tickHandler ConsumerHandler) (exitFlag bool) {
	var (
		err        error
		bucketItem []string
	)
	//从Redis中读取指定条数数据
	if bucketItem, err = r.getFromTopic(topic); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":   err.Error(),
			"topic": topic,
			"desc":  "扫描bucket错误#bucket",
		}, "RedisDelayMqTickHandler")
		exitFlag = true
		return
	}

	// 集合为空
	if bucketItem == nil || len(bucketItem) == 0 {
		exitFlag = true
		return
	}

	if err = r.dealData(tickHandler, topic, bucketItem...); err != nil {
		exitFlag = true
		return
	}

	// 从Redis和数据库中删除数据
	if err = r.removeFromBucket(topic, bucketItem...); err != nil {
		exitFlag = true
	}
	return
}

// 扫描bucket, 取出延迟时间小于当前时间的Job
func (r *RedisDelayMq) tickHandler(t time.Time, topic string, tickHandler ConsumerHandler) {
	var exitF bool
	for {
		if exitF = r.readData(topic, tickHandler); exitF {
			break
		}
	}
	return
}

func (r *RedisDelayMq) removeFromBucket(topic string, bucketItem ...string) (err error) {
	var msg = make([]string, 0, len(bucketItem))
	desc := ""
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"err":        err.Error(),
			"topic":      topic,
			"desc":       desc,
			"msg":        msg,
			"bucketItem": bucketItem,
		}, "RedisDelayMqRemoveFromBucket")
	}()

	if r.PersistentObject != nil {
		if err = r.PersistentObject.RemoveData(bucketItem...); err != nil {
			desc = "PersistentObjectRemoveData"
			return
		}
	}

	var e error
	members := make([]interface{}, 0, len(bucketItem))
	for _, item := range bucketItem {
		members = append(members, item)
		hashKey := r.getDataSaveKey(topic, item)
		e = r.client.Del(r.Ctx, hashKey).Err()
		if e != nil {
			desc = "删除缓存数据异常"
			msg = append(msg, fmt.Sprintf("hashKey:%s,err:%s", hashKey, e.Error()))
		}
	}

	//redis3.0.2 版本支持 (NX不更新存在的成员。只添加新成员)
	err = r.client.ZRem(r.Ctx, topic, members...).Err()
	if err != nil {
		desc = "删除缓存(有序集合)数据异常"
	}
	return
}

func (r *RedisDelayMq) getValueKey(topic string, timestamp time.Time) (res string, score int64) {
	score = timestamp.UnixNano() / 1e6 //以着当前时间的毫秒值作为存储格式
	res = r.getDataSaveKey(topic, fmt.Sprintf("%d", score))
	return
}

func (r *RedisDelayMq) getDataSaveKey(topic, sc string) (res string) {
	res = fmt.Sprintf("%s:%s", topic, sc)
	return
}

//往延迟队列添加数据
//Param topic string
//Param data interface
//Param executeTimestamp 延迟队列数据执行时刻
//Return error
func (r *RedisDelayMq) Add(topic string, data DelayMqData) (err error) {

	hashKey, sc := r.getValueKey(topic, data.Timestamp)

	desc := ""
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(map[string]interface{}{
			"err": err.Error(), "topic": topic, "data": data, "desc": desc, "sc": sc, "score": data.Timestamp.Format(utils.DateTimeGeneral),
		})
	}()

	if r.PersistentObject != nil {
		if err = r.PersistentObject.
			AddData(data); err != nil {
			desc = "PersistentObjectAddData"
			return
		}
	}

	//用list存储数据的值
	if err = r.client.LPush(r.Ctx,
		hashKey,
		data).
		Err(); err != nil {
		desc = "RedisLPush"
		return
	}
	//redis3.0.2 版本支持 (NX不更新存在的成员。只添加新成员)
	err = r.client.ZAddNX(r.Ctx, topic,
		&redis.Z{
			Member: fmt.Sprintf("%d", sc),
			Score:  float64(sc),
		},
	).Err()
	if err != nil {
		desc = "RedisZAddNX"
	}
	return
}

//消费逻辑
func (r *RedisDelayMq) consumerRun(topic string, handler ConsumerHandler) {
	// ZREMRANGEBYSCORE key min max
	r.client.ZRemRangeByScore(context.TODO(), topic, `0`,
		fmt.Sprintf("%d", time.Now().UnixNano()/1e6),
	).Args()
}

func (r *RedisDelayMq) DefaultOption() {
	if r.Config.Delayer.TimerInterval == 0 { //默认一分钟一跳
		r.Config.Delayer.TimerInterval = 1 * time.Second
	}
	if r.Config.Delayer.Limit == 0 {
		r.Config.Delayer.Limit = DelayerDefaultLimit
	}
	r.initConsumerMapGroupIdMessageHandler()
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
