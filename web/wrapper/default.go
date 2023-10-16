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

type (
	ArgWebSocket struct {
		UserHid int64  `json:"uid" form:"uid"`
		ToId    int64  `json:"to_id" form:"to_id"`
		MsgType string `json:"msg_type" form:"msg_type"`
		Pk      string `json:"pk"` //websocket的key
		base.ArgWebSocketBase
	}
	ArgumentDefault struct {
		IdKey string `json:"id_key" form:"id_key"`
	}
	ResultDefault struct {
		Users []models.User
	}
)

func (r *ArgumentDefault) SetPathParam(hid string) {

	return
}
