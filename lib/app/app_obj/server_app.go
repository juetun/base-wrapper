package app_obj

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

const TmpSignKey = "signxxx"

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

func GetHeaderAppName(c *gin.Context) (appName, secret string, err error) {
	URI := strings.TrimPrefix(c.Request.URL.Path, "/")
	if URI == "" {
		err = fmt.Errorf("get app name failure")
		return
	}
	urlString := strings.Split(URI, "/")
	appName = urlString[0]
	secret = TmpSignKey
	// TODO 通过appName获取签名值
	return
}
