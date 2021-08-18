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
)

// Md5 生成32位MD5摘要
// 字符串加密 md5算法
func Md5(str string) string {
	c := md5.New()
	c.Write([]byte(str))
	cipherStr := c.Sum(nil)
	return hex.EncodeToString(cipherStr)
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

type Paginate struct {
	Limit   int `json:"limit"`
	Count   int `json:"count"`
	Total   int `json:"total"`
	Last    int `json:"last"`
	Current int `json:"current"`
	Next    int `json:"next"`
}

func MyPaginate(count int64, limit int, page int) Paginate {
	res := round(int(count), limit)
	totalPage := res

	lastPage := 0
	if page-1 <= 0 {
		lastPage = 1
	} else {
		lastPage = page - 1
	}

	currentPage := 0
	if page >= res {
		currentPage = res
	} else {
		currentPage = page
	}

	nextPage := 0
	if page+1 >= res {
		nextPage = res
	} else {
		nextPage = page + 1
	}

	return Paginate{
		Limit:   limit,
		Count:   int(count),
		Total:   totalPage,
		Last:    lastPage,
		Current: currentPage,
		Next:    nextPage,
	}
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
	len := len(length)
	if len > 1 {
		err = fmt.Errorf("length 最多1个数字")
		return
	} else if len == 0 {
		len = 6
	}
	var buff bytes.Buffer

	for i := 0; i < length[0]; i++ {
		result, _ := rand.Int(rand.Reader, big.NewInt(10))
		buff.WriteString(result.String())
	}
	authCode = buff.String()
	return
}

// IsIdCard 判断身份证号是否合法
func IsIdCard(idCard string) (ok bool, err error) {
	if ok, err = regexp.Match("/^([1-6][1-9]|50)\\d{4}\\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\\d{3}$/", []byte(idCard)); err != nil {
		return
	}
	if ok {
		return
	}
	if ok, err = regexp.Match("/^([1-6][1-9]|50)\\d{4}(18|19|20)\\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$/", []byte(idCard)); err != nil {
		return
	}
	return
}
