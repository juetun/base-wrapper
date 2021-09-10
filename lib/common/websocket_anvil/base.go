package websocket_anvil

import (
	"encoding/json"
)

type BaseHandler struct {
}

func (r *BaseHandler) GetSuccessWithData(data interface{}) Resp {
	return r.GetResult(Ok, CustomError[Ok], data)
}
func (r *BaseHandler) GetFailWithMsg(msg string) Resp {
	return r.GetResult(NotOk, msg, map[string]interface{}{})
}
func (r *BaseHandler) GetSuccess() Resp {
	return r.GetResult(Ok, CustomError[Ok], map[string]interface{}{})
}

func (r *BaseHandler) Struct2StructByJson(struct1 interface{}, struct2 interface{}) (err error) {

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
func (r *BaseHandler) GetResult(code int, msg string, data interface{}) Resp {
	return Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	}
}
