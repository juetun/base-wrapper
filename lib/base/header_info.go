package base

import (
	"fmt"
)

type (
	HeaderInfo struct {
		DeviceId   string `json:"device_id" form:"device_id"`     // 用户设备号
		ClientType string `json:"client_type" form:"client_type"` // 终端类型 "m ,android,iso,weixin,alipay"
		Channel    string `json:"channel" form:"channel"`         // APP发布的渠道
		App        string `json:"app" form:"app"`                 // APP名称
		AppVersion string `json:"app_version" form:"app_version"` // APP版本号
	}
)

func NewHeaderInfo() (res *HeaderInfo) {
	res = &HeaderInfo{}
	return
}

func (r *HeaderInfo) StringToToken(ctx *Context) (headerInfoString string, err error) {
	if r == nil {
		r = &HeaderInfo{}
	}
	if headerInfoString, err = CreateTokenFromObject(r, ctx); err != nil {
		err = fmt.Errorf("系统异常，请刷新或稍后重试")
		return
	}
	return
}

func (r *HeaderInfo) ParseFromString(ctx *Context, headerInfoString string) (err error) {
	if headerInfoString == "" {
		return
	}
	if err = ParseJwtKey(headerInfoString, ctx, r); err != nil {
		err = fmt.Errorf("系统异常，请刷新或稍后重试")
		return
	}
	return
}
