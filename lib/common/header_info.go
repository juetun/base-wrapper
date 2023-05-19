package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type (
	HeaderInfo struct {
		App      string `json:"app,omitempty"`      //App类型
		Terminal string `json:"terminal,omitempty"` //终端
		Channel  string `json:"channel,omitempty"`  //渠道
		Version  string `json:"version,omitempty"`  //版本
		Debug    bool   `json:"debug,omitempty"`    //是否调试模式
		Lng      string `json:"lng,omitempty"`      //经度
		Lat      string `json:"lat,omitempty"`      //纬度
		Province string `json:"pro,omitempty"`      //省
		CityId   string `json:"city,omitempty"`     //市
		AreaId   string `json:"area,omitempty"`     //区
	}
)

func (r *HeaderInfo) InitHeaderInfo(ctx *gin.Context) (err error) {
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
		r = &tmp
	case *HeaderInfo:
		r = data.(*HeaderInfo)
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
