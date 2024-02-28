package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type (
	HeaderInfo struct {
		HApp      string      `json:"h_app,omitempty"`      //App类型
		HTerminal string      `json:"h_terminal,omitempty"` //终端
		HChannel  string      `json:"h_channel,omitempty"`  //渠道
		HVersion  string      `json:"h_version,omitempty"`  //版本
		HDebug    bool        `json:"h_debug,omitempty"`    //是否调试模式
		HLng      string      `json:"h_lng,omitempty"`      //经度
		HLat      string      `json:"h_lat,omitempty"`      //纬度
		HProvince string      `json:"h_pro,omitempty"`      //省
		HCityId   string      `json:"h_city,omitempty"`     //市
		HAreaId   string      `json:"h_area,omitempty"`     //
		HPk       string      `json:"h_pk,omitempty"`       //终端唯一号 如：设备ID 或其他区分设备唯一的号
		HExt      interface{} `json:"h_ext,omitempty"`
	}
)

func (r *HeaderInfo) setData(tmp *HeaderInfo) {
	r.HApp = tmp.HApp
	r.HTerminal = tmp.HTerminal
	r.HChannel = tmp.HChannel
	r.HVersion = tmp.HVersion
	r.HDebug = tmp.HDebug
	r.HLng = tmp.HLng
	r.HLat = tmp.HLat
	r.HProvince = tmp.HProvince
	r.HCityId = tmp.HCityId
	r.HAreaId = tmp.HAreaId
	r.HPk = tmp.HPk
	r.HExt = tmp.HExt
	return
}

func (r *HeaderInfo) InitHeaderInfo(ctx *gin.Context) (err error) {
	if r == nil {
		err = fmt.Errorf("HeaderInfo 对象未初始化")
		return
	}
	var (
		data interface{}
		ok   bool
	)
	if data, ok = ctx.Get(app_obj.HttpHeaderInfo); !ok || data == nil {
		err = fmt.Errorf("%v info is not exists", app_obj.HttpHeaderInfo)
		return
	}
	switch data.(type) {
	case HeaderInfo:
		tmp := data.(HeaderInfo)
		r.setData(&tmp)
	case *HeaderInfo:
		tmp := data.(*HeaderInfo)
		r.setData(tmp)
	default:
		err = fmt.Errorf("系统异常,%v信息错误", app_obj.HttpHeaderInfo)
	}
	return
}

func (r *HeaderInfo) Validate() (err error) {
	if r.HApp == "" {
		err = fmt.Errorf("please set app value")
		return
	}
	if r.HApp == "all" {
		err = fmt.Errorf("app can not is all")
		return
	}
	if r.HTerminal == "" {
		err = fmt.Errorf("please set terminal value")
		return
	}
	return
}
