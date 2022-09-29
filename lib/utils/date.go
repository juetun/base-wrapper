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
	DateTimeGeneral     = "2006-01-02 15:04:05"
	DateTimeGeneralNano = "2006-01-02 15:04:05.999999999" //纳秒时间格式
	DateTimeDashboard   = "2006.01.02 15:04"
	DateGeneral         = "2006-01-02"
	DateTimeChat        = "01月02 15:04" //聊天信息展示的时间格式
	TimeDay             = 24 * time.Hour
	TimeWeek            = 7 * 24 * time.Hour
)

// GetMondayDateStamp 获取指定时间的星期一凌晨0时0分0秒
func GetMondayDateStamp(t time.Time) (mondayStamp time.Time) {
	var offset int
	if offset = int(time.Monday - t.Weekday()); offset > 0 {
		offset = -6
	}
	mondayStamp = time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0,
		time.Local).
		AddDate(0, 0, offset)
	return
}

// DateTime 时间格式转换
func DateTime(t time.Time, format ...string) (res string) {

	var f = DateTimeGeneral
	if len(format) > 0 {
		f = format[0]
	}
	return t.Format(f)
}

func DateTimeDiff(valueTime, baseTime time.Time, format ...string) (res string, difTime time.Duration, err error) {
	if baseTime.IsZero() {
		baseTime = time.Now()
	}
	dif := baseTime.Unix() - valueTime.Unix()
	difTime = time.Duration(baseTime.UnixNano() - valueTime.UnixNano())
	if dif < 60 {
		res = "刚刚"
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
	res = valueTime.Format(DateGeneral)
	return
}

// DateDuration 指定时间离当前时间的差额
func DateDuration(value string, baseTime time.Time, format ...string) (res string, difTime time.Duration, err error) {
	var t time.Time
	if t, err = DateParse(value, format...); err != nil {
		return
	}
	return DateTimeDiff(t, baseTime, format...)

}

// DateDuration 指定时间离当前时间的差额
func DateDurationV1(value string, baseTime time.Time, format ...string) (res string, difTime time.Duration, err error) {
	var t time.Time
	if t, err = DateParse(value, format...); err != nil {
		return
	}
	return DateTimeDiff(t, baseTime, format...)

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

//将unix时间戳(单位:纳秒)转换为时间格式
func TimeUnixFromNano(unixTimeStampNano int64) (res time.Time) {
	res = time.Unix(
		int64(math.Floor(float64(unixTimeStampNano)/1e9)),
		unixTimeStampNano%1e9,
	)
	return
}

//将unix时间戳(单位:秒)转换为时间格式
func TimeUnixFrom(unixTimeStamp int64) (res time.Time) {
	res = time.Unix(
		unixTimeStamp,
		0,
	)
	return
}

// ParseDate 转换日期格式
func ParseDate(t time.Time) (res string) {
	return t.Format(DateGeneral)
}

// ParseDateTime 转换时间格式
func ParseDateTime(t time.Time) (res string) {
	return t.Format(DateTimeGeneral)
}
