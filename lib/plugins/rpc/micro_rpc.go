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
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/utils"
)

const (
	RetryTimesOne      = 1                      // 尝试请求次数
	RetryTimesDuration = 500 * time.Millisecond // 尝试请求次数(500毫秒后尝试重新调用)
)

type CollectParam struct {
	PathUrl string `json:"path_url"`
}
type RequestOptions struct {
	Context            *base.Context `json:"-"`                // 上下文传参 操作日志对象
	NotMicro           bool          `json:"not_micro"`        // 不是微服务应用 (默认false -是微服务内调用)
	Method             string        `json:"method"`           // http请求方法
	AppName            string        `json:"app_name"`         // 应用名
	URI                string        `json:"uri"`              // 请求的地址
	ConnectTimeOut     time.Duration `json:"connect_time_out"` // 请求连接超时时长 默认300毫秒(建立HTTP请求的时长)
	RequestTimeOut     time.Duration `json:"request_time_out"` // 获取请求时长 默认5秒(获取数据的时长)
	Header             http.Header   `json:"header"`           // 请求的header请求
	Value              url.Values    `json:"value"`            // 请求参数
	BodyJson           []byte        `json:"body_json"`        // 请求的body信息
	PathVersion        string        `json:"path_version"`
	RetryTimes         int           `json:"retry_times"`          // 失败重试次数,1不重试（只发送一次请求） 2尝试再请求一次
	RetryTimesDuration time.Duration `json:"retry_times_duration"` // 重试时间间隔
	CollectParams      CollectParam  `json:"collect_params"`
}

// 请求操作结构体
type httpRpc struct {
	Request *RequestOptions `json:"request"` // 请求参数
	Error   error           `json:"error"`   //
	Body    []byte          `json:"-"`
	BaseUrl string          `json:"base_url"`
	Uri     string          `json:"uri"`
	resp    *http.Response
	client  *http.Client
}

// NewHttpRpc 请求入口
func NewHttpRpc(params *RequestOptions) (r *httpRpc) {
	r = &httpRpc{}
	r.Request = params
	if r.Request.RetryTimesDuration == 0 {
		r.Request.RetryTimesDuration = RetryTimesDuration
	}
	return
}

// 初始化默认参数
func (r *RequestOptions) initDefault() {
	if r.Method == "" {
		r.Method = http.MethodGet
	}
	if r.ConnectTimeOut == 0 {
		r.ConnectTimeOut = 1 * time.Second
	}
	if r.RequestTimeOut == 0 {
		r.RequestTimeOut = 5 * time.Second
	}
	if r.RetryTimes == 0 {
		r.RetryTimes = RetryTimesOne
	}

	// 不是访问内部服务
	if r.NotMicro {
		return
	}
	if r.Header == nil {
		r.Header = http.Header{}
	}
	traceId := ""
	if nil != r.Context {
		traceId = r.Context.GinContext.GetHeader(app_obj.HttpTraceId)
	}
	r.Header.Add(app_obj.HttpTraceId, traceId)
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

func (r *httpRpc) beforeSend() {
	if r.Error = r.Request.validateParams(); r.Error != nil {
		return
	}
	if r.Request.initDefault(); r.Error != nil {
		return
	}
	if r.DefaultBaseUrl(); r.Error != nil {
		return
	}
	r.Request.Method = strings.ToUpper(r.Request.Method)
}

func (r *httpRpc) GetResp() (res *http.Response) {
	return r.resp
}

// Send 发送请求
func (r *httpRpc) Send() (res *httpRpc) {
	res = r
	if r.beforeSend(); r.Error != nil {
		return
	}

	// 多次尝试发送请求(默认一次)
	var needBreak bool
	var i = 0
	for {
		if i < r.Request.RetryTimes {
			r.Error = nil
			if needBreak = r.sendAct(); needBreak {
				break
			}
			if r.Request.RetryTimesDuration > 0 {
				time.Sleep(r.Request.RetryTimesDuration)
			}
		} else {
			break
		}
		i++
	}
	return
}

func (r *httpRpc) sendAct() (needBreak bool) {
	switch r.Request.Method {
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
		needBreak = true
		return
	}
	// 判断请求状态
	if r.resp != nil && r.resp.StatusCode == http.StatusOK {
		needBreak = true
		return
	}
	return
}

func (r *httpRpc) dialContextHandlerfunc(ctx context.Context, network, addr string) (conn net.Conn, err error) {
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
}

func (r *httpRpc) initClient() (res *httpRpc) {
	res = r
	r.client = &http.Client{Transport: &http.Transport{DialContext: r.dialContextHandlerfunc,},}
	return
}

func (r *httpRpc) DefaultBaseUrl() {

	if !r.Request.NotMicro { // 如果不是微服务应用
		if _, ok := app_obj.AppMap[r.Request.AppName]; !ok {
			r.Error = fmt.Errorf("base url config is not exists(%s)", r.Request.AppName)
			return
		}
		if r.Request.PathVersion == "" {
			r.BaseUrl = fmt.Sprintf("%s/%s", app_obj.AppMap[r.Request.AppName], r.Request.AppName)
			return
		}
		if strings.Index(r.Request.PathVersion, "/") == -1 {
			r.BaseUrl = fmt.Sprintf("%s/%s/%s", app_obj.AppMap[r.Request.AppName], r.Request.AppName, r.Request.PathVersion)
			return
		}
		r.BaseUrl = fmt.Sprintf("%s/%s%s", app_obj.AppMap[r.Request.AppName], r.Request.AppName, r.Request.PathVersion)
		return
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
	r.Uri = r.getUrl(r.Request.Value.Encode())

	request, r.Error = http.NewRequest(http.MethodGet, r.Uri, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) delete() {
	var request *http.Request
	r.Uri = r.getUrl(r.Request.Value.Encode())
	request, r.Error = http.NewRequest(http.MethodDelete, r.Uri, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) put() {
	var request *http.Request
	r.Uri = r.getUrl()
	if len(r.Request.BodyJson) > 0 {
		r.Request.Header.Add("Content-Type", "application/json")
		request, r.Error = http.NewRequest(http.MethodPut, r.Uri, bytes.NewReader(r.Request.BodyJson))
	} else {
		request, r.Error = http.NewRequest(http.MethodPut, r.Uri, nil)
	}
	request.Header = r.Request.Header
	if len(r.Request.Value) > 0 {
		request.PostForm = r.Request.Value
	}
	r.resp, r.Error = r.client.Do(request)
}

func (r *httpRpc) patch() {
	var request *http.Request
	r.Uri = r.getUrl()
	if len(r.Request.BodyJson) > 0 {
		r.Request.Header.Add("Content-Type", "application/json")
		request, r.Error = http.NewRequest(http.MethodPatch, r.Uri, bytes.NewReader(r.Request.BodyJson))
	} else {
		request, r.Error = http.NewRequest(http.MethodPatch, r.Uri, nil)
	}
	request.Header = r.Request.Header
	if len(r.Request.Value) > 0 {
		request.PostForm = r.Request.Value
	}
	r.resp, r.Error = r.client.Do(request)
}

func (r *httpRpc) post() {
	r.Uri = r.getUrl()
	if len(r.Request.BodyJson) > 0 {
		if r.postJson(); r.Error != nil {
			return
		}
		return
	}
	if r.postGeneral(); r.Error != nil {
		return
	}
}

// 发送请求
func (r *httpRpc) sendDo(request *http.Request) {
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}

func (r *httpRpc) postGeneral() {
	var request *http.Request
	r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request, r.Error = http.NewRequest(http.MethodPost, r.Uri, strings.NewReader(r.Request.Value.Encode()))
	if r.Error != nil {
		return
	}
	r.sendDo(request)
	return
}

func (r *httpRpc) postJson() {
	var request *http.Request
	r.Request.Header.Add("Content-Type", "application/json")
	request, r.Error = http.NewRequest(http.MethodPost, r.Uri, bytes.NewReader(r.Request.BodyJson))
	if r.Error != nil {
		return
	}
	r.sendDo(request)
	return
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
	res = r
	logContent := map[string]interface{}{
		"request": r.Request,
		"uri":     r.Uri,
	}
	if r.Request.BodyJson != nil {
		logContent["reqBody"] = string(r.Request.BodyJson)
	}
	defer func() {
		if r.Request == nil || r.Request.Context == nil {
			return
		}
		if r.Error != nil {
			logContent["err"] = r.Error.Error()
			r.Request.Context.Error(logContent, "httpRpcGetBody")
		} else {
			r.Request.Context.Info(logContent, "httpRpcGetBody")
		}

	}()

	if r.Error != nil {
		return
	}
	// 保证I/O正常关闭
	defer func() {
		if r.resp != nil && r.resp.Body != nil {
			_ = r.resp.Body.Close()
		}
	}()

	if r.resp == nil {
		r.Error = fmt.Errorf("请求发送失败")
		return
	}
	// 判断请求状态
	if r.resp.StatusCode != 200 {
		r.Error = fmt.Errorf("请求失败(%d)", r.resp.StatusCode)
		return
	}
	// 失败，返回状态
	if r.Body, r.Error = io.ReadAll(r.resp.Body); r.Error != nil {
		// 读取错误,返回异常
		r.Error = fmt.Errorf("读取请求返回失败(%s)", r.Error.Error())
		return
	}

	logContent["body"] = string(r.Body)

	// 成功，返回数据及状态
	return

}
