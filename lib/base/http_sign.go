// Package signencrypt
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package base

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GolangCharset 签名的字符编码类型
type GolangCharset string

// 字符编码类型常量

// 当前类的指针
var sign *SignUtils

// SignUtils 签名类
type SignUtils struct {
	mapExtend *MapExtend
}

type MapExtend struct {
}

func (r *MapExtend) GetKeys(data map[string]string) (res []string, err error) {
	res = make([]string, 0, len(data))
	for key := range data {
		res = append(res, key)
	}
	return
}

// NewSign 实例化签名
func NewSign() *SignUtils {
	sign = &SignUtils{
		mapExtend: &MapExtend{},
	}
	return sign
}

// SignTopRequest 签名算法
// parameters 要签名的数据项
// secret 生成的publicKey
// signMethod 签名的字符编码
func (s *SignUtils) SignTopRequest(parameters map[string]string, secret string) (bb bytes.Buffer, err error) {

	/**
	  1、第一步：把字典按Key的字母顺序排序
	  2、第二步：把所有参数名和参数值串在一起
	  3、第三步：使用MD5/HMAC加密
	  4、第四步：把二进制转化为大写的十六进制
	*/

	// 第一步：把字典按Key的字母顺序排序
	var keys []string
	if keys, err = s.mapExtend.GetKeys(parameters); err != nil {
		return
	} else {
		sort.Strings(keys)
	}

	// 第二步：把所有参数名和参数值串在一起

	bb.WriteString(secret)

	for _, v := range keys {
		if val := parameters[v]; len(val) > 0 {
			bb.WriteString(v)
			bb.WriteString(val)
		}
	}
	return
}

type ListenHandler func(s string)
type ListenHandlerStruct struct {
	MD5HMAC       ListenHandler // 转换成 MD5后执行
	ByteTo16After ListenHandler // 把二进制转化为大写的十六进制
	FinishHandler ListenHandler // 返回签名完成的字符串
}

func (s *SignUtils) Encrypt(argJoin string, secret string, listenHandlerStruct ListenHandlerStruct) (res string) {

	// Crypto by HMAC-SHA1
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(argJoin))

	//进行base64编码
	res = hex.EncodeToString(mac.Sum([]byte(nil)))

	// 返回签名完成的字符串
	res = strings.ToLower(string(res))
	if listenHandlerStruct.MD5HMAC != nil {
		listenHandlerStruct.MD5HMAC(res)
	}
	// 第四步：把二进制转化为大写的十六进制
	if listenHandlerStruct.ByteTo16After != nil {
		listenHandlerStruct.ByteTo16After(res)
	}

	if listenHandlerStruct.FinishHandler != nil {
		listenHandlerStruct.FinishHandler(res)
	}
	return
}

func (s *SignUtils) getHeaderTimeStamp(header http.Header) (headerT string, err error) {
	var t int
	// 判断签名是否传递了时间
	if headerT = header.Get(app_obj.HttpTimestamp); headerT == "" {
		err = fmt.Errorf("the header must be include timestamp parameter(t)")
		return
	} else if t, err = strconv.Atoi(headerT); err != nil {
		err = fmt.Errorf("格式不不正确(时间戳:%s)", app_obj.HttpTimestamp)
		return
	} else if app_obj.App.AppEnv != app_obj.EnvProd && int(time.Now().UnixNano()/1e6)-t > 86400000 { // 传递的时间格式必须大于当前时间-一天
		err = fmt.Errorf("the header of  parameter(t) must be more than now desc one days")
		return
	}

	return

}

type SignParams struct {
	Method            string `json:"method"`
	Path              string `json:"path"`
	HeaderTimeStamp   string `json:"header_time_stamp"`
	HeaderContentType string `json:"header_content_type"`
	RequestBody       []byte `json:"request_body"`
}

func ConcatString(signParams *SignParams) (base64Code, encryptionString string, err error) {
	var (
		bt bytes.Buffer
	)
	bt.WriteString(signParams.Method)
	bt.WriteString(signParams.Path)
	if _, err = bt.WriteString(signParams.HeaderTimeStamp); err != nil {
		return
	}
	bt.Write(signParams.RequestBody)

	encryptionString = strings.ToLower(bt.String())

	base64Code = base64.StdEncoding.EncodeToString([]byte(encryptionString))
	if len([]byte(base64Code)) > 400 {
		base64Code = base64Code[0:400]
	}
	return
}

// SignGinRequest http请求加密算法
// c *gin.Context,
func (s *SignUtils) SignGinRequest(c *gin.Context) (validateResult bool, signResult string, err error) {

	//如果是内网访问接口
	if ok := InterPath(c); ok {
		signResult = c.Request.Header.Get(app_obj.HttpSign)
		validateResult = true
		return
	}
	var (
		headerT                      string
		encryptionCode               bytes.Buffer
		secret                       string
		requestBody                  []byte
		base64Code, encryptionString string
	)
	if _, secret, err = app_obj.GetHeaderAppName(c); err != nil {
		return
	}

	if headerT, err = s.getHeaderTimeStamp(c.Request.Header); err != nil {
		return
	}

	// 如果传JSON 单独处理
	if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		if requestBody, err = io.ReadAll(c.Request.Body); err != nil {
			return
		}
		// 读完body参数一定要回写，不然后边取不到参数
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		requestBody = append([]byte(secret), requestBody...)
	} else { // 如果是非JSON 传参
		// 如果不是JSON 则直接过去FORM表单参数
		if encryptionCode, err = s.SortParamsAndJoinData(s.GetRequestParams(c), secret); err != nil {
			return
		}
		requestBody = encryptionCode.Bytes()
	}
	var (
		signParams = &SignParams{
			Method:            c.Request.Method,
			Path:              c.Request.URL.Path,
			HeaderTimeStamp:   headerT,
			HeaderContentType: c.GetHeader("Content-Type"),
			RequestBody:       requestBody,
		}
		listenHandlerStruct ListenHandlerStruct
	)
	if base64Code, encryptionString, err = ConcatString(signParams); err != nil {
		return
	}
	listenHandlerStruct = ListenHandlerStruct{}
	switch app_obj.App.AppEnv {
	case app_obj.EnvProd:
		if c.GetBool(app_obj.DebugFlag) {
			listenHandlerStruct = s.testAddHeader(c, base64Code, encryptionString)
		}
	default: // 如果不是线上环境,可输出签名格式 (此处代码为调试 签名是否能正常使用准备)
		listenHandlerStruct = s.testAddHeader(c, base64Code, encryptionString)
	}

	signResult = s.Encrypt(base64Code, secret, listenHandlerStruct)
	if signResult == c.Request.Header.Get(app_obj.HttpSign) {
		validateResult = true
	}
	return
}

func (s *SignUtils) testAddHeader(c *gin.Context, base64Code, encryptionString string) (listenHandlerStruct ListenHandlerStruct) {
	// 配置回调输出
	listenHandlerStruct = ListenHandlerStruct{}
	resp := c.Writer.Header()
	resp.Set("Sign-format", encryptionString)
	resp.Set("Sign-Base64Code", base64Code)
	listenHandlerStruct = ListenHandlerStruct{
		MD5HMAC:       func(s string) {},
		ByteTo16After: func(s string) { resp.Set("Sign-ByteTo16", s) },
		FinishHandler: func(s string) { resp.Set("Sign-f", s) },
	}
	return
}

// 加密字符串
func (s *SignUtils) SortParamsAndJoinData(data map[string]string, secret string) (res bytes.Buffer, err error) {
	if res, err = s.SignTopRequest(data, secret); err != nil {
		return
	}
	return
}
func (s *SignUtils) ParamsConcat(params url.Values) (valueMap map[string]string) {
	valueMap = make(map[string]string, len(params))
	for k, v := range params {
		valueMap[k] = strings.Join(v, ";")
	}
	return
}
func (s *SignUtils) GetRequestParams(c *gin.Context) (valueMap map[string]string) {
	valueMap = make(map[string]string, len(c.Request.PostForm))
	_ = c.Request.ParseMultipartForm(128) // 保存表单缓存的内存大小128M
	valueMap = s.ParamsConcat(c.Request.Form)
	return
}
