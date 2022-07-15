package redis_mq

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
)

// 需要版本 3.0
// Redis实现的订阅与发布

const (
	RedisMQBaseValue = "close"
)

type (
	// RedisMQ redis mq操作对象
	RedisMQ struct {
		Ctx         context.Context // 上下文参数 用于停止监听动作
		cancel      context.CancelFunc
		Context     *base.Context // 上下文参数用于记录日志使用
		ChannelList []string `json:"channel_list"`

		client *redis.Client
		pb     *redis.PubSub
	}


	// RedisMqOption RedisMQ赋值属性参数
	RedisMqOption func(mq *RedisMQ)

	// ActionHandler 接收消息处理句柄
	ActionHandler func(message string, channel string) (err error)
)

// NewRedisMQ 初始化redisMQ对象
func NewRedisMQ(options ...RedisMqOption) (res *RedisMQ) {
	res = &RedisMQ{}
	for _, option := range options {
		option(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	res.Ctx, res.cancel = context.WithCancel(res.Ctx)
	return res
}

// RedisOptionClient 设置参数方法
func RedisOptionClient(client *redis.Client) RedisMqOption {
	return func(mq *RedisMQ) {
		mq.client = client
	}
}

// RedisOptionContext 设置参数方法
func RedisOptionContext(context *base.Context) RedisMqOption {
	return func(mq *RedisMQ) {
		mq.Context = context
	}
}

// RedisOptionCtx 设置参数方法
func RedisOptionCtx(Ctx context.Context) RedisMqOption {
	return func(mq *RedisMQ) {
		mq.Ctx = Ctx
	}
}

// RedisOptionChannelList 设置参数方法
func RedisOptionChannelList(ChannelList []string) RedisMqOption {
	return func(mq *RedisMQ) {
		mq.ChannelList = ChannelList
	}
}

// Subscribe 开启订阅的渠道
func (r *RedisMQ) Subscribe() {
	r.pb = r.client.Subscribe(r.Ctx, r.ChannelList...)
}

// PUBLISH 发布消息
func (r *RedisMQ) PUBLISH(channel, message string) {
	res := r.client.Publish(r.Ctx, channel, message)
	var err error
	if err = res.Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":     err.Error(),
			"channel": channel,
			"message": message,
		}, "RedisMQPUBLISH")
	}
	return
}

// AddSubscribe 添加订阅渠道
func (r *RedisMQ) AddSubscribe(channelList ...string) {

	if len(channelList) == 0 {
		return
	}
	if r.pb == nil {
		r.pb = r.client.Subscribe(r.Ctx, channelList...)
		return
	}
	r.client.Subscribe(r.Ctx, channelList...)
}

// Unsubscribe 退订渠道
func (r *RedisMQ) Unsubscribe(channelList ...string) (err error) {
	if len(channelList) == 0 {
		return
	}
	if err = r.pb.Unsubscribe(r.Ctx, channelList...); err != nil {
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
		}, "RedisMQUnsubscribe")
	}
	return
}

// Close 停止监听
func (r *RedisMQ) Close() {
	r.ChannelList = nil
	r.cancel()
	r.pb = nil
}

// AcceptMsg 监听数据
func (r *RedisMQ) AcceptMsg(actionHandler ActionHandler) {
	for {
		select {
		case msg := <-r.pb.Channel():
			// 等待从 channel 中发布 close 关闭服务
			switch msg.Payload {
			case RedisMQBaseValue: // 如果监听到了退出指令
				r.Context.Info(map[string]interface{}{
					"desc": "收到监听消息指令",
				}, "RedisMQAcceptMsg")
				r.Close()
				break
			default:
				_ = actionHandler(msg.Payload, msg.Channel)
			}
		default:

		}
	}
}


