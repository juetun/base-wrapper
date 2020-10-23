package utils

/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-04
 * Time: 22:29
 */

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/juetun/base-wrapper/lib/app_log"
)

// Define AlarmType to string
// for to check the params is right
type AlarmType string

type AlarmMailReceive string

// this are some const params what i defined
// only this can be to input
const (
	AlarmTypeOne   AlarmType = "mail"
	AlarmTypeTwo   AlarmType = "wechat"
	AlarmTypeThree AlarmType = "message"
)

type AlarmParam struct {
	Types  AlarmType
	MailTo AlarmMailReceive
	Log    *app_log.AppLog
}

var alarmParam *AlarmParam

// Define a closure type to next
type ap func(*AlarmParam) (interface{}, error)

// can use this function to set a new value
// but to check it is a right type
func (alarm *AlarmParam) SetType(t AlarmType) ap {
	return func(alarm *AlarmParam) (interface{}, error) {
		str := strings.Split(string(t), ",")
		if len(str) == 0 {
			alarm.Log.Error(nil, map[string]interface{}{"content": "you must input a value"})
			return nil, errors.New("you must input a value")
		}
		for _, types := range str {
			s := AlarmType(types)
			_, err := s.IsCurrentType()
			if err != nil {
				alarm.Log.Error(nil, map[string]interface{}{"content": "the value type is error"})
				return nil, err
			}
		}
		ty := alarm.Types
		alarm.Types = t
		return ty, nil
	}
}

func (alarm *AlarmParam) SetMailTo(t AlarmMailReceive) ap {
	return func(alarm *AlarmParam) (interface{}, error) {
		to := alarm.MailTo
		_, err := t.CheckIsNull()
		if err != nil {
			return nil, err
		}
		_, err = t.MustMailFormat()
		if err != nil {
			return nil, err
		}
		alarm.MailTo = t
		return to, nil
	}
}

// alarm receive account can not null
func (t AlarmMailReceive) CheckIsNull() (AlarmMailReceive, error) {
	if len(t) == 0 {
		app_log.GetLog().Logger.Errorln("content", "value can not be null")
		return "", errors.New("value can not be null")
	}
	return t, nil
}

// alarm receive account must be mail format
func (t AlarmMailReceive) MustMailFormat() (AlarmMailReceive, error) {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", string(t)); !m {
		app_log.GetLog().Logger.Errorln("content", "value format is not right")
		return "", errors.New("value format is not right")
	}
	return t, nil
}

// judge it is a right type what i need
// if is it a wrong type, i must return a panic to above
func (at AlarmType) IsCurrentType() (res AlarmType, err error) {
	res = at
	switch at {
	case AlarmTypeOne:
	case AlarmTypeTwo:
	case AlarmTypeThree:
	default:
		app_log.GetLog().Logger.Errorln("content", "the alarm type is error")
		err = fmt.Errorf("the alarm type is error")
 	}
	return
}

// implementation value
func (alarm *AlarmParam) AlarmInit(options ...ap) error {
	q := &AlarmParam{
	}
	for _, option := range options {
		_, err := option(q)
		if err != nil {
			return err
		}
	}
	alarmParam = q
	return nil
}

func Alarm(content string) {
	types := strings.Split(string(alarmParam.Types), ",")
	var err error
	for _, a := range types {
		switch AlarmType(a) {
		case AlarmTypeOne:
			if alarmParam.MailTo == "" {
				app_log.GetLog().Logger.Errorln("content", "邮件接收者不能为空")
				break
			}
			err = SendMail(string(alarmParam.MailTo), "报警", content)
			break
		case AlarmTypeTwo:
			break
		case AlarmTypeThree:
			break
		}
		if err != nil {
			app_log.GetLog().Logger.Errorln("content", "alarm is error,err:"+err.Error())
		}
	}
}
