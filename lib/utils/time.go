// @Copyright (c) 2021. 
// @Author ${USER}
// @Date ${DATE}
package utils

import "time"

//转换日期格式
func ParseDate(t time.Time) (res string) {
	return t.Format("2006-01-02")
}

//转换时间格式
func ParseDateTime(t time.Time) (res string) {
	return t.Format("2006-01-02 15:04:05")
}
