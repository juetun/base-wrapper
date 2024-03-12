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
	DateTimeGeneral       = "2006-01-02 15:04:05"
	DateTimeGeneralNano   = "2006-01-02 15:04:05.999999999" //纳秒时间格式
	DateTimeDashboard     = "2006.01.02 15:04"
	DateTimeDashboardShow = "2006.01.02 15:04:05"
	DateGeneral           = "2006-01-02"
	DateTimeChat          = "01月02 15:04" //聊天信息展示的时间格式
	TimeDay               = 24 * time.Hour
	TimeWeek              = 7 * 24 * time.Hour
)

const ( //表示时间为空的字符串
	DateNullString1       = "0001-01-01 00:00:00"
	DateNullString2       = "0001-01-02 00:00:00"
	DateNullStringDefault = "2000-01-01 00:00:00"
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

func getNextTimeDesc(abs float64, valueTime time.Time, format string) (res string) {

	if difM := math.Ceil(abs / 60); difM < 60 {
		res = fmt.Sprintf("%d分后", int(difM))
		return
	}
	if difH := math.Ceil(abs / 3600); difH < 24 {
		res = fmt.Sprintf("%d小时后", int(difH))
		return
	}
	if difD := math.Ceil(abs / 86400); difD < 7 {
		res = fmt.Sprintf("%d天后", int(difD))
		return
	}
	if difW := math.Ceil(abs / (86400 * 7)); difW < 4 {
		res = fmt.Sprintf("%d周后", int(difW))
		return
	}
	res = valueTime.Format(format)
	return
}

func DateTimeDiff(valueTime, baseTime time.Time, formats ...string) (res string, difTime time.Duration, err error) {
	if baseTime.IsZero() {
		baseTime = time.Now()
	}
	var format = DateGeneral
	if len(formats) > 0 {
		format = formats[0]
	}
	dif := baseTime.Unix() - valueTime.Unix()
	difTime = time.Duration(baseTime.UnixNano() - valueTime.UnixNano())
	if dif < 60 && dif > 0 {
		res = "刚刚"
		return
	}
	var (
		abs = float64(dif)
	)
	if dif < 0 {
		res = getNextTimeDesc(math.Abs(float64(dif)), valueTime, format)
		return
	}
	if difM := math.Floor(abs / 60); difM < 60 {
		res = fmt.Sprintf("%d分前", int(difM))
		return
	}
	if difH := math.Floor(abs / 3600); difH < 24 {
		res = fmt.Sprintf("%d小时前", int(difH))
		return
	}
	if difD := math.Floor(abs / 86400); difD < 7 {
		res = fmt.Sprintf("%d天前", int(difD))
		return
	}
	if difW := math.Floor(abs / (86400 * 7)); difW < 4 {
		res = fmt.Sprintf("%d周前", int(difW))
		return
	}
	res = valueTime.Format(format)
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

//时间时间格式转换
func TimeZeroToString(timeStamp *time.Time, formats ...string) (res string) {
	if timeStamp == nil || timeStamp.IsZero() {
		return
	}
	var formatString = DateTimeDashboard
	if len(formats) > 0 {
		formatString = formats[0]
	}
	res = timeStamp.Format(formatString)
	if res == "0001.01.01 00:00" || res == "0001.01.01" {
		res = ""
		return
	}
	return
}
