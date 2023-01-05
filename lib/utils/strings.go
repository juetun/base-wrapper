// Package utils
/**
* Author:changjiang
* Description:
* File:strings
* Version: 1.0.0
* Date 2020/3/19 11:45 下午
 */
package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand" // 真随机
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Md5 生成32位MD5摘要
// 字符串加密 md5算法
func Md5(str string) string {
	c := md5.New()
	c.Write([]byte(str))
	cipherStr := c.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//获取字符串的utf8格式(中文格式)长度
func StringUtf8Length(s string) (res int) {
	res = utf8.RuneCountInString(s)
	return
}

func round(a int, b int) int {
	rem := a % b
	dis := a / b
	if rem > 0 {
		return dis + 1
	} else {
		return dis
	}
}

//URL转换为字符串格式(用于权限验证生成KEY)
func UrlParseKey(s string) (res string) {
	for i := 0; i < len(s); i++ {
		switch string(s[i]) {
		case "/":
			res += `#`
			continue
		default:
			res += string(s[i])
		}
	}
	return
}

// IsDigit 判断一个字符串是否为数字
// +0.00001是数字
// -0.1111是数字
// -0.11.11不是数字

func IsDigit(str string) (res bool) {
	dotNum := 0
	doUnicode := '.'
	doUnicodeAdd := '+'
	doUnicodeDesc := '-'
	runeString := []rune(str)
	for k, x := range runeString {
		if k == 0 && (x == doUnicodeAdd || x == doUnicodeDesc) { // 首字母为"+"或"-"
			continue
		}
		if !unicode.IsDigit(x) {
			if x == doUnicode { // 如果是小数点
				// 如果小数点在第一位或者最后一位，则不是数字
				if k == 0 || k == len(runeString)-1 {
					return
				}
				dotNum++
				continue
			}
			return
		}
	}
	if dotNum > 1 {
		return
	}
	res = true
	return
}

//手机号加星号
func HidTel(phone string) (res string, err error) {
	res = phone
	if phone == "" {
		return
	}
	var ok bool
	if ok, err = regexp.Match(`^1[0-9]{10}$`, []byte(phone)); err != nil {
		return
	} else if ok {
		res = phone[:3] + "****" + phone[7:]
		return
	}
	if ok, err = regexp.Match(`^\+861[0-9]{10}$`, []byte(phone)); err != nil {
		return
	} else if ok {
		res = phone[:6] + "****" + phone[10:]
		return
	}
	if ok, err = regexp.Match(`^0[0-9]{2,3}[\-]?[0-9]{3,4}[\-]?[0-9]{4,8}$`, []byte(phone)); err != nil {//^0[0-9]{2,3}[\-]?[0-9]{7,8}[\-]?[0-9]?$
		return
	} else if ok {
		phoneArr := strings.Split(phone, "-")

		switch len(phoneArr) {
		case 1:
			res = phone[:6] + "***" + phone[10:]
			return
		case 2:
			phoneArr[1]=phoneArr[1][:2] + "***" + phoneArr[1][6:]
			res = strings.Join(phoneArr, "-")
			return
		case 3:
			phoneArr[2]=phoneArr[2][:2] + "***" + phoneArr[2][6:]
			res = strings.Join(phoneArr, "-")
			return
		}


	}

	return
}
func ReplaceStringByRegex(str, rule, replace string) (string, error) {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return "", fmt.Errorf("正则MustCompile错误:%s", err.Error())
	}
	return reg.ReplaceAllString(str, replace), nil
}

// SubString 汉字截取
func SubString(strValue string, num int, suffix ...string) (res string) {
	var b []int32
	var i = 0
	for _, value := range strValue {
		if i >= num {
			break
		}
		b = append(b, value)
		i++
	}
	res = string(b) + strings.Join(suffix, "")
	return
}

// GetStringUniqueNumber 将字符串转换为数字
// param strValue
// return int64
func GetStringUniqueNumber(strValue string) (num int64) {
	sp := []rune(strValue)
	for _, value := range sp {
		num += int64(value)
	}
	return
}

// 短信验证码字符串生成

// RandomString 短信验证码字符串生成
func RandomString(length ...int) (authCode string, err error) {
	var lengthNumber = 6
	if len(length) > 0 {
		lengthNumber = length[0]
		return
	}
	var buff bytes.Buffer

	for i := 0; i < lengthNumber; i++ {
		result, _ := rand.Int(rand.Reader, big.NewInt(10))
		buff.WriteString(result.String())
	}
	authCode = buff.String()
	return
}

// IsIdCard 判断身份证号是否合法
func IsIdCard(idCard string) (ok bool, err error) {
	if ok, err = regexp.Match(`^[1-9]\d{17}$`, []byte(idCard)); err != nil {
		return
	}
	if ok {
		return
	}
	if ok, err = regexp.Match(`^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{4}$`, []byte(idCard)); err != nil {
		return
	}
	return
}
