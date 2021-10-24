// Package rpc
/**
* @Author:ChangJiang
* @Description:
* @File:micro_rpc
* @Version: 1.0.0
* @Date 2020/10/18 12:56 下午
 */
package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
)

type CollectParam struct {
	PathUrl string `json:"path_url"`
}
type RequestOptions struct {
	Method         string        `json:"method"`   // http请求方法
	AppName        string        `json:"app_name"` // 应用名
	URI            string        `json:"uri"`
	Header         http.Header   `json:"header"`
	Value          url.Values    `json:"value"`
	BodyJson       []byte        `json:"body_json"` // json数据传递
	PathVersion    string        `json:"path_version"`
	NotMicro       bool          `json:"not_micro"`        // 不是微服务应用
	Context        *base.Context `json:"-"`                // 上下文传参 操作日志对象
	ConnectTimeOut time.Duration `json:"connect_time_out"` // 请求连接超时时长 默认300毫秒(建立HTTP请求的时长)
	RequestTimeOut time.Duration `json:"request_time_out"` // 获取请求时长 默认5秒(获取数据的时长)
	CollectParams  CollectParam  `json:"collect_params"`
}

// 请求操作结构体
type httpRpc struct {
	Request *RequestOptions `json:"request"` // 请求参数
	Error   error           `json:"error"`   //
	Body    []byte          `json:"-"`
	BaseUrl string          `json:"base_url"`
	resp    *http.Response
	client  *http.Client
}

// NewHttpRpc 请求入口
func NewHttpRpc(params *RequestOptions) (r *httpRpc) {
	r = &httpRpc{}
	if r.Error = params.validateParams(); r.Error != nil {
		return
	}
	params.initDefault()
	r.Request = params
	return
}

// 初始化默认参数
func (r *RequestOptions) initDefault() {
	if r.Method == "" {
		r.Method = "GET"
	}
	if r.ConnectTimeOut == 0 {
		r.ConnectTimeOut = 1 * time.Second
	}
	if r.RequestTimeOut == 0 {
		r.RequestTimeOut = 5 * time.Second
	}
	// 不是访问内部服务
	if r.NotMicro {
		return
	}
	traceId := ""
	if nil != r.Context {
		traceId = r.Context.GinContext.GetHeader(app_obj.HttpTraceId)
	}
	if r.Header == nil {
		r.Header = http.Header{}
	}
	r.Header.Add(app_obj.HttpTraceId, traceId)

	// if !r.NotMicro && r.PathVersion == "" {
	// 	r.PathVersion = "v1"
	// }

}

// 校验参数
func (r *RequestOptions) validateParams() (err error) {

	if !r.NotMicro && r.AppName == "" {
		err = fmt.Errorf("您没有选择的应用名(%s)", r.AppName)
		return
	}
	if r.URI == "" {
		err = fmt.Errorf("您没有输入的请求路径")
		return
	}
	return
}

// Send 发送请求
func (r *httpRpc) Send() (res *httpRpc) {
	res = r
	if r.Error != nil {
		return
	}
	r.DefaultBaseUrl()
	if r.Error != nil {
		return
	}
	switch strings.ToUpper(r.Request.Method) {
	case "GET":
		r.initClient().
			get()
	case "POST":
		r.initClient().
			post()
	case "PUT":
		r.initClient().
			put()
	case "DELETE":
		r.initClient().
			delete()
	case "PATCH":
		r.initClient().
			patch()
	default:
		r.Error = fmt.Errorf("当前不支您输入的请求方法(%s)", r.Request.Method)
	}
	return
}
func (r *httpRpc) initClient() (res *httpRpc) {
	res = r
	r.client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
				deadline := time.Now().Add(r.Request.RequestTimeOut)
				if conn, err = net.DialTimeout(network, addr, r.Request.ConnectTimeOut); err != nil {
					r.Request.Context.Error(map[string]interface{}{
						"err":      err.Error(),
						"network":  network,
						"addr":     addr,
						"deadline": deadline.Format(utils.DateTimeGeneral),
					}, "rpcHttpRpcInitClient")
					return
				}
				if err = conn.SetDeadline(deadline); err != nil {
					r.Request.Context.Error(map[string]interface{}{
						"err":      err.Error(),
						"network":  network,
						"addr":     addr,
						"deadline": deadline.Format(utils.DateTimeGeneral),
					}, "rpcHttpRpcInitClient")
					return
				}
				return
			},
		},
	}
	return
}

func (r *httpRpc) DefaultBaseUrl() {

	if !r.Request.NotMicro { // 如果不是微服务应用
		if _, ok := app_obj.AppMap[r.Request.AppName]; !ok {
			r.Error = fmt.Errorf("base url config is not exists(%s)", r.Request.AppName)
			return
		}
		r.BaseUrl = fmt.Sprintf("%s/%s%s", app_obj.AppMap[r.Request.AppName], r.Request.AppName, r.Request.PathVersion)
	}
}

func (r *httpRpc) getUrl(paramString ...string) (res string) {
	if len(paramString) > 0 {
		res = fmt.Sprintf("%s%s?%s", r.BaseUrl, r.Request.URI, strings.Join(paramString, "&"))
	} else {
		res = fmt.Sprintf("%s%s", r.BaseUrl, r.Request.URI)
	}
	r.Request.CollectParams.PathUrl = res
	return
}

func (r *httpRpc) get() {
	var request *http.Request
	urlPath := r.getUrl(r.Request.Value.Encode())

	request, r.Error = http.NewRequest(http.MethodGet, urlPath, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) delete() {
	var request *http.Request
	urlPath := r.getUrl(r.Request.Value.Encode())
	request, r.Error = http.NewRequest(http.MethodDelete, urlPath, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) put() {
	var request *http.Request
	urlPath := r.getUrl()
	if len(r.Request.BodyJson) > 0 {
		r.Request.Header.Add("Content-Type", "application/json")
		request, r.Error = http.NewRequest(http.MethodPut, urlPath, bytes.NewReader(r.Request.BodyJson))
	} else {
		request, r.Error = http.NewRequest(http.MethodPut, urlPath, nil)
	}
	request.Header = r.Request.Header
	if len(r.Request.Value) > 0 {
		request.PostForm = r.Request.Value
	}
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) patch() {
	var request *http.Request
	urlPath := r.getUrl()
	if len(r.Request.BodyJson) > 0 {
		r.Request.Header.Add("Content-Type", "application/json")
		request, r.Error = http.NewRequest(http.MethodPatch, urlPath, bytes.NewReader(r.Request.BodyJson))
	} else {
		request, r.Error = http.NewRequest(http.MethodPatch, urlPath, nil)
	}
	request.Header = r.Request.Header
	if len(r.Request.Value) > 0 {
		request.PostForm = r.Request.Value
	}
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) post() {
	var request *http.Request
	urlPath := r.getUrl()
	if len(r.Request.BodyJson) > 0 {
		r.Request.Header.Add("Content-Type", "application/json")
		request, r.Error = http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(r.Request.BodyJson))
	} else {
		r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request, r.Error = http.NewRequest(http.MethodPost, urlPath, strings.NewReader(r.Request.Value.Encode()))
	}
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}

// SetURLParams 生成GET URL结构
func (r *httpRpc) SetURLParams(data map[string]string) (res *httpRpc) {
	res = r
	u, _ := url.Parse(r.Request.URI)
	q := u.Query()
	for k, v := range data {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	r.Request.URI = u.String()
	return
}

func (r *httpRpc) Bind(obj interface{}) (res *httpRpc) {
	res = r
	if r.Error != nil {
		return
	}
	if len(r.Body) > 0 {
		r.Error = json.Unmarshal(r.Body, obj)
	}
	return
}

func (r *httpRpc) GetBodyAsString() (res string) {
	if len(r.Body) > 0 {
		res = string(r.Body)
	}
	return
}

func (r *httpRpc) GetBody() (res *httpRpc) {
	logContent := map[string]interface{}{
		"request": r.Request,
	}
	defer func() {
		if r.Error != nil {
			logContent["err"] = r.Error.Error()
			r.Request.Context.Error(logContent, "httpRpcGetBody")
		} else {
			r.Request.Context.Info(logContent, "httpRpcGetBody")
		}

	}()
	res = r
	if r.Error != nil {
		return
	}
	// 保证I/O正常关闭
	defer func() {

		if r.resp != nil && r.resp.Body != nil {
			_ = r.resp.Body.Close()
		}
	}()
	// 判断请求状态
	if r.resp.StatusCode != 200 {
		r.Error = fmt.Errorf("请求失败(%d)", r.resp.StatusCode)
		return
	}
	// 失败，返回状态
	r.Body, r.Error = ioutil.ReadAll(r.resp.Body)
	if r.Error != nil {
		// 读取错误,返回异常
		r.Error = fmt.Errorf("读取请求返回失败(%s)", r.Error.Error())
		return
	}

	logContent["body"] = string(r.Body)

	// 成功，返回数据及状态
	return

}
