package ext

import (
	"time"
)

const (
	// ClientConnectMax 最大的socket连接数（当前服务器连接数小于此值时，性能较优）
	ClientConnectMax = 3000

	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// HeartBeatPeriod 心跳间隔
	HeartBeatPeriod = 3 * time.Second

	// 最后一次活跃上线通知时间间隔
	lastActiveRegisterPeriod = 10 * time.Minute

	// HeartBeatMaxRetryCount 心跳最大重试次数
	HeartBeatMaxRetryCount = 3

	// MessageReqHeartBeat 消息请求类型
	// 第1个数字: 1请求, 2响应
	// 第2个数字: 消息种类(请求和响应的消息种类没有直接联系)
	// 第3个数字: 种类排序
	// 心跳消息
	MessageReqHeartBeat string = "1-1-1"
	// MessageReqPush 推送新消息
	MessageReqPush string = "1-2-1"
	// MessageReqBatchRead 批量已读
	MessageReqBatchRead string = "1-2-2"
	// MessageReqBatchDeleted 批量删除
	MessageReqBatchDeleted string = "1-2-3"
	// MessageReqAllRead 全部已读
	MessageReqAllRead string = "1-2-4"
	// MessageReqAllDeleted 全部删除
	MessageReqAllDeleted string = "1-2-5"

	// MessageRespHeartBeat 消息响应类型(首字符为2)
	// 心跳消息
	MessageRespHeartBeat string = "2-1-1"
	// MessageRespNormal 普通消息
	MessageRespNormal string = "2-2-1"
	// MessageRespUnRead 未读数
	MessageRespUnRead string = "2-3-1"
	// MessageRespOnline 用户上线
	MessageRespOnline string = "2-4-1"
)
const (
	Ok                  = 201
	NotOk               = 405
	Unauthorized        = 401
	Forbidden           = 403
	InternalServerError = 500
)
const (
	OkMsg                      = "操作成功"
	NotOkMsg                   = "操作失败"
	UnauthorizedMsg            = "登录过期, 需要重新登录"
	LoginCheckErrorMsg         = "用户名或密码错误"
	ForbiddenMsg               = "无权访问该资源, 请联系网站管理员授权"
	InternalServerErrorMsg     = "服务器内部错误"
	IdempotenceTokenEmptyMsg   = "幂等性token为空"
	IdempotenceTokenInvalidMsg = "幂等性token失效, 重复提交"
	UserDisabledMsg            = "账户已被禁用, 请联系网站管理员"
)

var CustomError = map[int]string{
	Ok:                  OkMsg,
	NotOk:               NotOkMsg,
	Unauthorized:        UnauthorizedMsg,
	Forbidden:           ForbiddenMsg,
	InternalServerError: InternalServerErrorMsg,
}
