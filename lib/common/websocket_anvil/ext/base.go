package ext

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

// WebsocketBaseHandler 公共操作父结构体
type WebsocketBaseHandler struct {
	Context   *base.Context `json:"-"`
	RequestId string        `json:"request_id"`
}

// Contains 判断uint数组是否包含item元素
func (r *WebsocketBaseHandler) Contains(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

type ErrorHandlerCallBack func(e interface{}) (res interface{}, needExit bool, err error)

// ErrorHandler
// return needExit 是否需要接触当前操作
func (r *WebsocketBaseHandler) ErrorHandler(callBack ErrorHandlerCallBack) (needExit bool) {
	var err error
	var e interface{}
	if e = recover(); e == nil {
		return
	}

	var res interface{}
	res, needExit, _ = callBack(e)
	_, file, line, _ := runtime.Caller(1)
	logContent := map[string]interface{}{
		"callBackRes": res,
		"e":           e,
		"loc":         fmt.Sprintf("%s(l:%d)", file, line),
		"desc":        fmt.Sprintf("[消息中心][接收端]连接可能已断开"),
	}
	if err != nil {
		logContent["err"] = err.Error()
	}
	r.Context.Error(logContent, "MessageClientErrorHandler")

	return
}

// ContainsIndex 判断uint数组是否包含item元素, 返回index
func (r *WebsocketBaseHandler) ContainsIndex(arr []string, item string) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}

// ContainsThenRemove 判断uint数组是否包含item元素, 并移除
func (r *WebsocketBaseHandler) ContainsThenRemove(arr []string, item string) []string {
	index := r.ContainsIndex(arr, item)
	if index >= 0 {
		arr = append(arr[:index], arr[index+1:]...)
	}
	return arr
}
func (r *WebsocketBaseHandler) GetSuccessWithData(data interface{}) Resp {
	return r.GetResult(Ok, CustomError[Ok], data)
}
func (r *WebsocketBaseHandler) GetFailWithMsg(msg string) Resp {
	return r.GetResult(NotOk, msg, map[string]interface{}{})
}
func (r *WebsocketBaseHandler) GetSuccess(msg ...string) Resp {
	var data = map[string]interface{}{}
	if len(msg) > 0 {
		data["content"] = msg[0]
	} else {
		data["content"] = CustomError[Ok]
	}
	return r.GetResult(Ok, CustomError[Ok], data)
}

func (r *WebsocketBaseHandler) Struct2StructByJson(struct1 interface{}, struct2 interface{}) (err error) {

	var bt []byte
	// 结构体转结构体, json为中间桥梁, struct2必须以指针方式传递, 否则可能获取到空数据
	// 转换为响应结构体, 隐藏部分字段
	if bt, err = json.Marshal(struct1); err != nil {
		return
	}
	if err = json.Unmarshal(bt, struct2); err != nil {
		return
	}

	return
}
func (r *WebsocketBaseHandler) GetResult(code int, msg string, data interface{}) Resp {
	return Resp{
		Code:      code,
		Data:      data,
		Msg:       msg,
		RequestId: r.RequestId,
		Time:      fmt.Sprintf("%d", time.Now().UnixNano()/1e6),
	}
}
