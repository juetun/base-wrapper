// Package wrapper
/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:18 下午
 */
package wrapper

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/models"
)

const (
	ChatMsgTypeSingleChart = "single" // 单聊
	ChatMsgTypeRoomChart   = "room"   // 群聊
)

var (
	SliceChatMsgType = base.ModelItemOptions{
		{
			Label: "单聊",
			Value: ChatMsgTypeSingleChart,
		},
		{
			Label: "群聊",
			Value: ChatMsgTypeRoomChart,
		},
	}
)

type (
	ArgWebSocket struct {
		App        string `json:"app" form:"app"` //来源APP
		FromType   uint8  `json:"from_type" form:"from_type"`
		FromId     int64  `json:"from_id" form:"from_id"`
		ToId       int64  `json:"to_id" form:"to_id"`
		ToType     uint8  `json:"msg_type" form:"msg_type"`
		Pk         string `json:"pk"` //websocket的key
		XAuthToken string `json:"x_auth_token" form:"x_auth_token"`
		base.ArgWebSocketBase
	}
	ArgumentDefault struct {
		IdKey string `json:"id_key" form:"id_key"`
	}
	ResultDefault struct {
		Users []models.User
	}
)

func (r *ArgumentDefault) Default(ctx *base.Context) (err error) {

	return
}

func (r *ArgumentDefault) SetPathParam(hid string) {

	return
}
