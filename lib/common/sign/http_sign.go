// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//签名的字符编码类型
type GOLANG_CHARSET string

//字符编码类型常量
const (
	CHARSET_ISO_2022_JP            GOLANG_CHARSET = "ISO-2022-JP"
	CHARSET_ISO_2022_CN                           = "ISO-2022-CN"
	CHARSET_ISO_2022_KR                           = "ISO-2022-KR"
	CHARSET_ISO_8859_5                            = "ISO-8859-5"
	CHARSET_ISO_8859_7                            = "ISO-8859-7"
	CHARSET_ISO_8859_8                            = "ISO-8859-8"
	CHARSET_BIG5                                  = "BIG5"
	CHARSET_GB18030                               = "GB18030"
	CHARSET_EUC_JP                                = "EUC-JP"
	CHARSET_EUC_KR                                = "EUC-KR"
	CHARSET_EUC_TW                                = "EUC-TW"
	CHARSET_SHIFT_JIS                             = "SHIFT_JIS"
	CHARSET_IBM855                                = "IBM855"
	CHARSET_IBM866                                = "IBM866"
	CHARSET_KOI8_R                                = "KOI8-R"
	CHARSET_MACCYRILLIC                           = "x-mac-cyrillic"
	CHARSET_WINDOWS_1251                          = "WINDOWS-1251"
	CHARSET_WINDOWS_1252                          = "WINDOWS-1252"
	CHARSET_WINDOWS_1253                          = "WINDOWS-1253"
	CHARSET_WINDOWS_1255                          = "WINDOWS-1255"
	CHARSET_UTF_8                                 = "UTF-8"
	CHARSET_UTF_16BE                              = "UTF-16BE"
	CHARSET_UTF_16LE                              = "UTF-16LE"
	CHARSET_UTF_32BE                              = "UTF-32BE"
	CHARSET_UTF_32LE                              = "UTF-32LE"
	CHARSET_TIS_620                               = "WINDOWS-874"
	CHARSET_HZ_GB_2312                            = "HZ-GB-2312"
	CHARSET_X_ISO_10646_UCS_4_3412                = "X-ISO-10646-UCS-4-3412"
	CHARSET_X_ISO_10646_UCS_4_2143                = "X-ISO-10646-UCS-4-2143"
)

//当前类的指针
var sign *signUtils

//同步锁
var signone sync.Once

//签名类
type signUtils struct {
	mapExtend *MapExtend
}

type MapExtend struct {
}

func (r *MapExtend) GetKeys(data map[string]string) (res []string, err error) {
	res = make([]string, 0, len(data))
	for key, _ := range data {
		res = append(res, key)
	}
	return
}

//实例化签名
func Sign() *signUtils {
	signone.Do(func() {
		sign = new(signUtils)
		sign.mapExtend = new(MapExtend)
	})
	return sign
}

/**
签名算法
parameters 要签名的数据项
secret 生成的publicKey
signMethod 签名的字符编码
*/
func (s *signUtils) SignTopRequest(parameters map[string]string, secret string, signMethod GOLANG_CHARSET) (res string, err error) {

	/**
	  1、第一步：把字典按Key的字母顺序排序
	  2、第二步：把所有参数名和参数值串在一起
	  3、第三步：使用MD5/HMAC加密
	  4、第四步：把二进制转化为大写的十六进制
	*/

	//第一步：把字典按Key的字母顺序排序
	var keys []string
	if keys, err = s.mapExtend.GetKeys(parameters); err != nil {
		return
	} else {
		sort.Strings(keys)
	}

	//第二步：把所有参数名和参数值串在一起
	var bb bytes.Buffer
	if CHARSET_UTF_8 == signMethod {
		bb.WriteString(secret)
	}
	for _, v := range keys {
		if val := parameters[v]; len(val) > 0 {
			bb.WriteString(v)
			bb.WriteString(val)
		}
	}

	fmt.Println(bb.String())

	//第三步：使用MD5/HMAC加密
	b := make([]byte, 0)
	if CHARSET_UTF_8 == signMethod {
		h := hmac.New(md5.New, s.GetUtf8Bytes(secret))
		h.Write(bb.Bytes())
		b = h.Sum(nil)
	} else {
		bb.WriteString(secret)
		h := md5.New()
		h.Write(bb.Bytes())
		b = h.Sum(nil)
	}

 	//第四步：把二进制转化为大写的十六进制
	var result bytes.Buffer
	for i := 0; i < len(b); i++ {
		s := strconv.FormatInt(int64(b[i]&0xff), 16)
		if len(s) == 1 {
			result.WriteString("0")
		}
		result.WriteString(s)
	}
	//返回签名完成的字符串
	res = strings.ToUpper(result.String())
	return
}

//默认utf8字符串
func (s *signUtils) GetUtf8Bytes(str string) []byte {
	b := []byte(str)
	return b
}
