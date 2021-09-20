package ext

import (
	"strings"
)

type (

	// MessageHandler 普通消息处理逻辑
	MessageHandler func(userHid string, data interface{}) (res interface{}, err error)

	// MessageBroadcast 消息广播
	MessageBroadcast struct {
		MessageWsResponseStruct
		UserIds []string `json:"-"`
	}

	// MessageWsResponseStruct websocket消息响应
	MessageWsResponseStruct struct {
		// 消息类型, 见const
		Type string `json:"type"`
		// 消息详情
		Detail Resp `json:"detail"`
	}

	// Resp http请求响应封装
	Resp struct {
		Code      int         `json:"code"`      // 错误代码
		Data      interface{} `json:"data"`      // 数据内容
		Msg       string      `json:"msg"`       // 消息提示
		RequestId string      `json:"requestId"` // 请求id
		Time      string      `json:"time"`      // 时间戳
	}

	// PushMessageRequestStruct 推送消息结构体
	PushMessageRequestStruct struct {
		FromUserId       string   `json:"fromUserId" form:"fromUserId"`
		Type             *ReqUint `json:"type" form:"type" validate:"required"`
		ToUserIds        []string `json:"toUserIds" form:"toUserIds"`
		ToRoleIds        []string `json:"toRoleIds" form:"toRoleIds"`
		Title            string   `json:"title" form:"title" validate:"required"`
		Content          string   `json:"content" form:"content" validate:"required"`
		IdempotenceToken string   `json:"idempotenceToken" form:"idempotenceToken"`
	}
	// MessageWebSocketRequestStruct websocket消息请求
	MessageWebSocketRequestStruct struct {
		// 消息类型, 见const
		Type string `json:"type"`
		// 数据内容
		Data interface{} `json:"data"`
	}

	// ReqUint 请求uint类型
	ReqUint string

	// Req 适用于大多数场景的请求参数绑定
	Req struct {
		Ids string `json:"ids" form:"ids"` // 传多个id
	}

	// UserInterface 用户信息
	UserInterface interface {
		// GetUserHid 获取用户信息
		GetUserHid() (userHid string, err error)
	}

	WebSocketAnvilOption func(arg *WebSocketAnvil)

	UserHandler func() (userHid UserInterface, err error)
)

// FieldTrans 翻译需要校验的字段名称
func (r *PushMessageRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Type"] = "消息类型"
	m["Title"] = "消息标题"
	m["Content"] = "消息内容"
	return m
}

// GetIds 获取
func (r *Req) GetIds() (ids []string) {
	idArr := strings.Split(r.Ids, ",")
	for _, v := range idArr {
		if v == "" {
			continue
		}
		ids = append(ids, v)
	}
	return
}
