/**
* @Author:changjiang
* @Description:
* @File:micro_rpc
* @Version: 1.0.0
* @Date 2020/10/18 12:56 下午
 */
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
)

type RequestOptions struct {
	Method         string        `json:"method"`
	AppName        string        `json:"app_name"`
	URI            string        `json:"uri"`
	Header         http.Header   `json:"header"`
	Value          url.Values    `json:"value"`
	PathVersion    string        `json:"path_version"`
	NotMicro       bool          `json:"not_micro"`        // 不是微服务应用
	Context        *base.Context `json:"context"`          // 上下文传参 操作日志对象
	ConnectTimeOut time.Duration `json:"connect_time_out"` // 请求连接超时时长 默认300毫秒(建立HTTP请求的时长)
	RequestTimeOut time.Duration `json:"request_time_out"` // 获取请求时长 默认5秒(获取数据的时长)
}

// 请求操作结构体
type httpRpc struct {
	Request *RequestOptions `json:"request"` // 请求参数
	Error   error           `json:"error"`   //
	body    []byte          `json:"-"`
	BaseUrl string          `json:"base_url"`
	resp    *http.Response
	client  *http.Client
}

// 请求入口
func NewHttpRpc(params *RequestOptions) (r *httpRpc) {
	r = &httpRpc{}
	r.Error = params.validateParams()
	if r.Error != nil {
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
		r.ConnectTimeOut = 300 * time.Millisecond
	}
	if r.RequestTimeOut == 0 {
		r.RequestTimeOut = 5 * time.Second
	}
	traceId := ""
	if nil != r.Context {
		r.Context.GinContext.GetHeader(app_obj.HTTP_TRACE_ID)
	}
	if r.Header == nil {
		r.Header = http.Header{
			app_obj.HTTP_TRACE_ID: []string{traceId},
		}
	} else {
		r.Header.Add(app_obj.HTTP_TRACE_ID, traceId)
	}
	if !r.NotMicro && r.PathVersion == "" {
		r.PathVersion = "v1"
	}
}

// 校验参数
func (r *RequestOptions) validateParams() (err error) {
	if r.AppName == "" {
		err = fmt.Errorf("您没有选择的应用名(%s)", r.AppName)
		return
	}
	if r.URI == "" {
		err = fmt.Errorf("您没有输入的请求路径")
		return
	}
	return
}

// 发送请求
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
		r.initClient().get()
	case "POST":
		r.initClient().post()
	case "PUT":
		r.initClient().put()
	case "DELETE":
		r.initClient().delete()
	case "PATCH":
		r.initClient().patch()
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
				conn, err = net.DialTimeout(network, addr, r.Request.ConnectTimeOut)
				if err != nil {
					return
				}
				conn.SetDeadline(deadline)
				return
			},
		},
	}
	return
}
func (r *httpRpc) DefaultBaseUrl() {
	if _, ok := app_obj.AppMap[r.Request.AppName]; !ok {
		r.Error = fmt.Errorf("base url config is not exists(%s)", r.Request.AppName)
		return
	}
	if !r.Request.NotMicro { // 如果不是微服务应用
		r.BaseUrl = fmt.Sprintf("%s/%s/%s", app_obj.AppMap[r.Request.AppName], r.Request.AppName, r.Request.PathVersion)
	}
}

func (r *httpRpc) getUrl() (res string) {
	return fmt.Sprintf("%s%s", r.BaseUrl, r.Request.URI)
}

func (r *httpRpc) get() {
	var request *http.Request
	url := r.getUrl()
	request, r.Error = http.NewRequest("GET", url, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) delete() {
	var request *http.Request
	url := r.getUrl()
	request, r.Error = http.NewRequest("DELETE", url, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) put() {
	var request *http.Request
	url := r.getUrl()
	request, r.Error = http.NewRequest("PUT", url, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) patch() {
	var request *http.Request
	url := r.getUrl()
	request, r.Error = http.NewRequest("PATCH", url, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}
func (r *httpRpc) post() {
	var request *http.Request
	url := r.getUrl()
	request, r.Error = http.NewRequest("POST", url, nil)
	request.Header = r.Request.Header
	r.resp, r.Error = r.client.Do(request)
}

// 生成GET URL结构
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
	if len(r.body) > 0 {
		r.Error = json.Unmarshal(r.body, obj)
	}
	return
}

func (r *httpRpc) GetBody() (res *httpRpc) {
	res = r
	if r.Error != nil {
		return
	}
	// 保证I/O正常关闭
	defer r.resp.Body.Close()
	// 判断请求状态
	if r.resp.StatusCode != 200 {
		r.Error = fmt.Errorf("请求失败(%d)", r.resp.StatusCode)
		return
	}
	// 失败，返回状态
	r.body, r.Error = ioutil.ReadAll(r.resp.Body)
	if r.Error != nil {
		// 读取错误,返回异常
		r.Error = fmt.Errorf("读取请求返回失败(%s)", r.Error.Error())
		return
	}
	// 成功，返回数据及状态
	return

}