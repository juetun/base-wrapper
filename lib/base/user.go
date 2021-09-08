package base

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type RequestUser struct {
	UserHid           string     `json:"user_hid" form:"user_hid"`
	UserMobileIndex   string     `json:"user_mobile_index" form:"user_mobile_index"`
	UserEmailIndex    string     `json:"user_email_index" form:"user_email_index"`
	Portrait          string     `json:"portrait" form:"portrait"`
	NickName          string     `json:"nick_name" form:"nick_name"`
	UserName          string     `json:"user_name" form:"user_name"`
	Gender            int        `json:"gender" form:"gender"`
	Status            int        `json:"status" form:"status"`
	Score             int        `json:"score" form:"score"`
	RememberToken     string     `json:"remember_token" form:"remember_token"`
	MsgReadTimeCursor TimeNormal `json:"msg_read_time_cursor" form:"msg_read_time_cursor"`
}

func (r *RequestUser) InitRequestUser(c *gin.Context) (err error) {
	if r.UserHid == "" {
		r.UserHid = c.GetHeader(app_obj.HttpUserHid)
	}
	return
}
