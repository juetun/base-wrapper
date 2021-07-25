// Package utils
/**
* @Author:ChangJiang
* @Description:
* @File:date
* @Version: 1.0.0
* @Date 2021/4/18 6:13 下午
 */
package utils

import (
	"fmt"
	"math"
	"time"
)

const (
	DateTimeGeneral = "2006-01-02 15:04:05"
	DateGeneral     = "2006-01-02"
	TimeDay         = 86400 * time.Second
	TimeWeek        = 7 * 86400 * time.Second
)

func DateTime(t time.Time, format ...string) (res string) {
	var f = DateTimeGeneral
	if len(format) > 0 {
		f = format[0]
	}
	return t.Format(f)
}

// DateDuration 指定时间离当前时间的差额
func DateDuration(value string, baseTime *time.Time, format ...string) (res string, err error) {
	var t time.Time
	if t, err = DateParse(value, format...); err != nil {
		return
	}
	if baseTime == nil {
		*baseTime = time.Now()
	}
	dif := baseTime.Unix() - t.Unix()
	if dif < 60 {
		res = fmt.Sprintf("%d秒前", int(dif))
		return
	}
	if difM := math.Floor(float64(dif / 60)); difM < 60 {
		res = fmt.Sprintf("%d分前", int(difM))
		return
	}
	if difH := math.Floor(float64(dif / 3600)); difH < 24 {
		res = fmt.Sprintf("%d小时前", int(difH))
		return
	}
	if difD := math.Floor(float64(dif / 86400)); difD < 7 {
		res = fmt.Sprintf("%d天前", int(difD))
		return
	}
	if difW := math.Floor(float64(dif / (86400 * 7))); difW < 4 {
		res = fmt.Sprintf("%d周前", int(difW))
		return
	}
	res = t.Format(DateGeneral)
	return
}
func Date(t time.Time) (res string) {
	return DateTime(t, DateGeneral)
}

// DateParse 解析时间格式
func DateParse(value string, format ...string) (stamp time.Time, err error) {
	var f = DateTimeGeneral
	if len(format) > 0 {
		f = format[0]
	}

	// ======= 将时间字符串转换为时间戳 =======
	stamp, err = time.ParseInLocation(f, value, time.Local)

	return
}

// DateStandard 标准的golang格式
func DateStandard(value string) (t time.Time, err error) {
	return DateParse(value, time.RFC3339)
}
