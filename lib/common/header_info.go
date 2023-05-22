package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type (
	HeaderInfo struct {
		App      string      `json:"app,omitempty"`      //App类型
		Terminal string      `json:"terminal,omitempty"` //终端
		Channel  string      `json:"channel,omitempty"`  //渠道
		Version  string      `json:"version,omitempty"`  //版本
		Debug    bool        `json:"debug,omitempty"`    //是否调试模式
		Lng      string      `json:"lng,omitempty"`      //经度
		Lat      string      `json:"lat,omitempty"`      //纬度
		Province string      `json:"pro,omitempty"`      //省
		CityId   string      `json:"city,omitempty"`     //市
		AreaId   string      `json:"area,omitempty"`     //
		Ext      interface{} `json:"ext,omitempty"`
	}
)

func (r *HeaderInfo) setData(tmp *HeaderInfo) {
	r.App = tmp.App
	r.Terminal = tmp.Terminal
	r.Channel = tmp.Channel
	r.Version = tmp.Version
	r.Debug = tmp.Debug
	r.Lng = tmp.Lng
	r.Lat = tmp.Lat
	r.Province = tmp.Province
	r.CityId = tmp.CityId
	r.AreaId = tmp.AreaId
	r.Ext = tmp.Ext
	return
}

func (r *HeaderInfo) InitHeaderInfo(ctx *gin.Context) (err error) {
	if r == nil {
		r = &HeaderInfo{}
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
	if r.App == "" {
		err = fmt.Errorf("please set app value")
		return
	}
	if r.App == "all" {
		err = fmt.Errorf("app can not is all")
		return
	}
	if r.Terminal == "" {
		err = fmt.Errorf("please set terminal value")
		return
	}
	return
}
