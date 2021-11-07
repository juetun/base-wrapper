// Package signencrypt
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package signencrypt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
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

	var bb bytes.Buffer
	bb.WriteString(argJoin)

	// 第三步：使用MD5/HMAC加密
	b := make([]byte, 0)

	h := hmac.New(sha1.New, []byte(secret))
	h.Write(bb.Bytes())
	b = h.Sum(nil)
	b = []byte(base64.StdEncoding.EncodeToString(b))

	// 返回签名完成的字符串
	res = strings.ToLower(string(b))
	if listenHandlerStruct.MD5HMAC != nil {
		listenHandlerStruct.MD5HMAC(string(b))
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

type GetSecretHandler func(appName string) (secret string, err error)

// SignGinRequest http请求加密算法
// c *gin.Context,
func (s *SignUtils) SignGinRequest(c *gin.Context, getSecret GetSecretHandler) (validateResult bool, signResult string, err error) {
	var appName, secret string
	if appName, err = s.getHeaderAppName(c); err != nil {
		return
	}
	if secret, err = getSecret(appName); err != nil {
		return
	}
	var bt bytes.Buffer
	var encryptionCode bytes.Buffer
	bt.WriteString(c.Request.Method)
	bt.WriteString(c.Request.URL.Path)

	var t int
	// 判断签名是否传递了时间
	if headerT := c.GetHeader("X-Timestamp"); headerT == "" {
		err = fmt.Errorf("the header must be include timestamp parameter(t)")
		return
	} else if t, err = strconv.Atoi(headerT); err != nil {
		err = fmt.Errorf("格式不不正确(时间戳:X-Timestamp)")
		return
	} else if app_obj.App.AppEnv != app_obj.EnvProd && int(time.Now().UnixNano()/1e6)-t > 86400000 { // 传递的时间格式必须大于当前时间-一天
		err = fmt.Errorf("the header of  parameter(t) must be more than now desc one days")
		return
	} else {
		bt.WriteString(headerT)
	}

	var body []byte
	// 如果传JSON 单独处理
	if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		bt.WriteString(secret)
		if body, err = ioutil.ReadAll(c.Request.Body); err != nil {
			return
		}
		// 读完body参数一定要回写，不然后边取不到参数
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	} else { // 如果是非JSON 传参
		// 如果不是JSON 则直接过去FORM表单参数
		if encryptionCode, err = s.sortParamsAndJoinData(s.getRequestParams(c), secret); err != nil {
			return
		}
		body = encryptionCode.Bytes()
	}
	bt.Write(body)
	encryptionString := strings.ToLower(bt.String())
	base64Code := base64.StdEncoding.EncodeToString([]byte(encryptionString))

	// 配置回调输出
	listenHandlerStruct := ListenHandlerStruct{}

	// 如果不是线上环境,可输出签名格式 (此处代码为调试 签名是否能正常使用准备)
	if app_obj.App.AppEnv != app_obj.EnvProd && c.GetHeader("debug") != "" {
		resp := c.Writer.Header()
		resp.Set("Sign-format", encryptionString)
		resp.Set("Sign-Base64Code", base64Code)
		listenHandlerStruct = ListenHandlerStruct{
			MD5HMAC:       func(s string) {},
			ByteTo16After: func(s string) { resp.Set("Sign-ByteTo16", s) },
			FinishHandler: func(s string) { resp.Set("Sign-f", s) },
		}
	}
	signResult = s.Encrypt(base64Code, secret, listenHandlerStruct)
	if signResult == c.GetHeader("X-Sign") {
		validateResult = true
	}
	return
}

// 加密字符串
func (s *SignUtils) sortParamsAndJoinData(data map[string]string, secret string) (res bytes.Buffer, err error) {
	if res, err = s.SignTopRequest(data, secret); err != nil {
		return
	}
	return
}

func (s *SignUtils) getRequestParams(c *gin.Context) (valueMap map[string]string) {
	valueMap = make(map[string]string, len(c.Request.PostForm))
	_ = c.Request.ParseMultipartForm(128) // 保存表单缓存的内存大小128M
	for k, v := range c.Request.Form {
		valueMap[k] = strings.Join(v, ";")
	}
	return
}

func (s *SignUtils) getHeaderAppName(c *gin.Context) (appName string, err error) {
	URI := strings.TrimPrefix(c.Request.URL.Path, "/")
	if URI == "" {
		err = fmt.Errorf("get app name failure")
		return
	}
	urlString := strings.Split(URI, "/")
	appName = urlString[0]
	return
}

// 默认utf8字符串
// func (s *SignUtils) GetUtf8Bytes(str string) []byte {
//	b := []byte(str)
//	return b
// }
