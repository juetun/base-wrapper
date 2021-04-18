/**
* @Author:changjiang
* @Description:
* @File:date
* @Version: 1.0.0
* @Date 2021/4/18 6:13 下午
 */
package utils

import (
	"time"
)

func DateTime(t time.Time) (res string) {
	return t.Format("2006-01-02 15:04:05")
}
func Date(t time.Time) (res string) {
	return t.Format("2006-01-02")
}
