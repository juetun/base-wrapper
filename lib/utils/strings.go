/**
* @Author:changjiang
* @Description:
* @File:strings
* @Version: 1.0.0
* @Date 2020/3/19 11:45 下午
 */
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(s string) string {
	c := md5.New()
	c.Write([]byte(s))
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

//汉字截取
func SubString(s string, num int, suffix ...string) (res string) {
	var b []int32
	var i = 0
	for _, value := range s {
		if i >= num {
			break
		}
		b = append(b, value)
		i++
	}
	res = string(b) + strings.Join(suffix, "")
	return
}