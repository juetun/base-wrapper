package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"net/http"
)

// 请求头info信息处理
func HttpHeaderInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodOptions:
			c.Next()
			return
		case http.MethodHead:
			c.Next()
			return
		}
		var (
			err                      error
			secret, headerInfoString string
			HttpHeaderInformation    common.HeaderInfo //http请求头预埋的其他信息
			infoByte                 []byte
		)

		if headerInfoString = c.Request.Header.Get(app_obj.HttpHeaderInfo); headerInfoString == "" {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  fmt.Sprintf("%s is null", app_obj.HttpHeaderInfo),
			})
			return
		}

		//app_obj.HttpHeaderInformation
		if _, secret, err = app_obj.GetHeaderAppName(c); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  err.Error(),
			})
			return
		}
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  "Header info sign key is null!",
			})
		}

		if headerInfoString, err = common.NewAes().
			DecryptCtr(headerInfoString, secret); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  err.Error(),
			})
			return
		} else if headerInfoString == "" {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  fmt.Sprintf("%s is null", app_obj.HttpHeaderInfo),
			})
			return
		}

		if infoByte, err = base64.StdEncoding.DecodeString(headerInfoString); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  fmt.Sprintf("header (key:%v) format is error(%v) base64 decode error", app_obj.HttpHeaderInfo, err.Error()),
			})
			return
		}
		if err = json.Unmarshal(infoByte, &HttpHeaderInformation); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  fmt.Sprintf("header (key:%v) format is error(%v)", app_obj.HttpHeaderInfo, err.Error()),
			})
			return
		}
		if err = HttpHeaderInformation.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  fmt.Sprintf("header (key:%v) format is error(%v)", app_obj.HttpHeaderInfo, err.Error()),
			})
			return
		}
		c.Set(app_obj.HttpHeaderInfo, HttpHeaderInformation)
		c.Next()
	}
}
