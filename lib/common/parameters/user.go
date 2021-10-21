package parameters

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
)

// 获取用户信息对应的表
const (
	UserDataTypeEmail  = "user_email"
	UserDataTypeMain   = "user_main"
	UserDataTypeInfo   = "user_info"
	UserDataTypeMobile = "user_mobile"
)

// 获取用户信息的响应参数结构
type (
	ResultUser struct {
		List map[string]ResultUserItem `json:"list"`
	}
	ResultUserItem struct {
		UserHid          string     `json:"user_hid,omitempty"`  // 用户ID
		Portrait         string     `json:"portrait,omitempty"`  // 头像
		NickName         string     `json:"nick_name,omitempty"` // 昵称
		UserName         string     `json:"user_name,omitempty"` // 用户名
		Gender           int        `json:"gender,omitempty"`    //
		Status           int        `json:"status,omitempty"`    //
		Score            int        `json:"score,omitempty"`     //
		AuthDesc         string     `json:"auth_desc,omitempty"` // 认证描述
		IsV              int        `json:"is_v,omitempty"`      // 用户头像加V
		Signature        string     `json:"signature,omitempty"`
		RegisterChannel  string     `json:"register_channel,omitempty"`
		CountryCode      string     `json:"country_code,omitempty"`
		Mobile           string     `json:"mobile,omitempty"`
		MobileVerifiedAt *time.Time `json:"mobile_verified_at,omitempty"`
		Email            string     `json:"email,omitempty"`
		EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty"`
		ShopId           string     `json:"shop_id"`
	}
	RequestUser struct {
		UserHid           string          `json:"user_hid" form:"user_hid"`
		UserMobileIndex   string          `json:"user_mobile_index" form:"user_mobile_index"`
		UserEmailIndex    string          `json:"user_email_index" form:"user_email_index"`
		Portrait          string          `json:"portrait" form:"portrait"`
		NickName          string          `json:"nick_name" form:"nick_name"`
		UserName          string          `json:"user_name" form:"user_name"`
		Gender            int             `json:"gender" form:"gender"`
		Status            int             `json:"status" form:"status"`
		Score             int             `json:"score" form:"score"`
		RememberToken     string          `json:"remember_token" form:"remember_token"`
		MsgReadTimeCursor base.TimeNormal `json:"msg_read_time_cursor" form:"msg_read_time_cursor"`
		ShopId            string          `json:"shop_id"`
	}
)


func (r *RequestUser) InitRequestUser(c *gin.Context) (err error) {
	if r.UserHid == "" {
		r.UserHid = c.GetHeader(app_obj.HttpUserHid)
	}
	return
}
