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
	"strings"
)

func GetRequestURIAndMethod(request *http.Request) (uri, method string) {
	if request == nil {
		return
	}
	uri = request.URL.Path
	method = request.Method
	return
}

func INConfigUrl(uri, method string, uriList []app_obj.UrlFormat) (res bool) {
	if uriList == nil {
		return
	}
	var (
		ok bool
	)
	for _, item := range uriList {
		if item.Uri == uri {
			if _, ok = item.Method[method]; ok {
				res = true
				return
			}
		}
		if item.IsPrefix { //如果是以前缀开始匹配
			if strings.Index(uri, item.Uri) == 0 {
				res = true
				return
			}
		}

	}
	return
}

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

		//websocket的参数验证跳过
		if c.Request.Header.Get("Upgrade") == "websocket" {
			c.Next()
			return
		}

		uri, method := GetRequestURIAndMethod(c.Request)
		//如果当前请求在配置的请求范围内
		if INConfigUrl(uri, method, app_obj.App.NotSendHeader) {
			c.Next()
			return
		}

		//如果网页请求 则不调用签名验证
		//if app_obj.App != nil && app_obj.App.AppRouterPrefix.Page != "" {
		//	//如果网页请求 则不调用签名验证
		//	if strings.Index(uri, "/"+app_obj.App.AppRouterPrefix.Page) == 0 {
		//		c.Next()
		//		return
		//	}
		//}

		var (
			err                   error
			secret                string
			headerInfoString      = c.Request.Header.Get(app_obj.HttpHeaderInfo)
			HttpHeaderInformation common.HeaderInfo //http请求头预埋的其他信息
			infoByte              []byte
		)

		if headerInfoString == "" {
			if ok := base.InterPath(c); !ok { //如果不是内网访问，切没传headerinfo 报错
				c.AbortWithStatusJSON(http.StatusOK, base.Result{
					Code: http.StatusUnauthorized,
					Msg:  fmt.Sprintf("%s is null", app_obj.HttpHeaderInfo),
				})
				return
			} else {
				c.Next()
				return
			}
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

		//如果支持签名调试
		if app_obj.App.AppSignDebug && !c.GetBool(app_obj.DebugFlag) { //DebugFlag只要不为true 可重复赋值
			c.Set(app_obj.DebugFlag, HttpHeaderInformation.HDebug)
		}
		c.Set(app_obj.HttpHeaderInfo, HttpHeaderInformation)
		c.Next()
	}
}
