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
	"strings"
)

//生成32位MD5摘要
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

// 汉字截取
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

// 将字符串转换为数字
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

// 短信验证码字符串生成
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
