/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-11
 * Time: 23:14
 */
package common

import (
	"strconv"
	"time"

	"github.com/juetun/base-wrapper/lib/utils"
)

func Offset(page string, limit string) (limitInt int, offset int) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limitInt, err = strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}

	return limitInt, (pageInt - 1) * limitInt
}

func GoMerge(arr1 []interface{}, arr2 []interface{}) []interface{} {
	for _, val := range arr2 {
		arr1 = append(arr1, val)
	}
	return arr1
}

func GoRepeat(str string, num int) string {
	var i int
	newStr := ""
	if num != 0 {
		for i = 0; i < num; i++ {
			newStr += str
		}
	}
	return newStr
}

func Rem(divisor int) bool {
	if (divisor+1)%4 == 0 {
		return true
	} else {
		return false
	}
}

func MDate(times time.Time) string {
	return times.Format(utils.DateTimeGeneral)
}

func MDate2(times time.Time) string {
	return times.Format("01-02")
}
