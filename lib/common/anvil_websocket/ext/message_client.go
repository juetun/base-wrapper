package ext

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// MessageClient 消息客户端
type MessageClient struct {
	MessageAction MessageHandler `json:"-"`

	WebsocketBaseHandler `json:"-"`

	// 当前socket key
	Key string `json:"key"`

	// 当前socket连接实例
	Conn *websocket.Conn `json:"conn"`

	// 当前登录用户
	UserFunc UserHandler `json:"-"`

	// 当前登录用户ip地址
	Ip string `json:"ip"`

	// 发送消息通道
	SendChan *Chan `json:"-"`

	// 上次活跃时间
	LastActiveTime time.Time `json:"last_active_time"`

	// 重试次数
	RetryCount uint `json:"retry_count"`
}

func (r *MessageClient) Register() {

	HubMessage.lock.Lock()
	defer HubMessage.lock.Unlock()
	var err error

	userHid, userData, err := r.GetUserHid()
	t := time.Now()
	active, ok := HubMessage.UserLastActive[userHid]
	last := time.Unix(active, 0) // carbon.CreateFromTimestamp(active)

	HubMessage.Clients[r.Key] = r

	// 记录最后活跃时间戳
	HubMessage.UserLastActive[userHid] = t.Unix()

	if !ok || last.Add(lastActiveRegisterPeriod).Before(t) {
		if !r.Contains(HubMessage.UserIds, userHid) {
			HubMessage.UserIds = append(HubMessage.UserIds, userHid)
		}
		logContent := map[string]interface{}{
			"key":     r.Key,
			"userHid": userHid,
			"Ip":      r.Ip,
			"desc":    fmt.Sprintf("[消息中心][用户上线]"),
		}

		defer func() {
			if err != nil {
				logContent["err"] = err.Error()
				r.Context.Error(logContent, "MessageClientRegister")
			} else {
				r.Context.Info(logContent, "MessageClientRegister")
			}
		}()

		go func() {
			HubMessage.RefreshUserMessage.SafeSend([]string{userHid})
		}()

		// 广播当前用户上线
		// 通知除自己之外的人
		go HubMessage.Broadcast.SafeSend(MessageBroadcast{
			MessageWsResponseStruct: MessageWsResponseStruct{
				Type: MessageRespOnline,
				Detail: r.GetSuccessWithData(map[string]interface{}{
					"user": userData,
				}),
			},
			UserIds: r.ContainsThenRemove(HubMessage.UserIds, userHid),
		})
	}

}

func (r *MessageClient) receiveMessageAct(msg []byte, userHid string) (needBreak bool, err error) {
	logContent := map[string]interface{}{
		"msg":     string(msg),
		"userHid": userHid,
	}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent)
			return
		}
		r.Context.Info(logContent)
	}()

	// 数据转为json
	var req MessageWebSocketRequestStruct
	if err = json.Unmarshal(msg, &req); err != nil {
		logContent["desc"] = fmt.Sprintf("[Json2Struct]转换异常")
		return
	}
	detail := r.GetSuccess()
	switch req.Type {
	case MessageReqBreak:
		needBreak = true
	case MessageReqHeartBeat: // 心跳消息
		err = r.msgHeartBreak(userHid, detail, req.Data)
	case MessageReqPush: // 推送新消息
		err = r.msgSendNewMsg(userHid, detail, req.Data)
	case MessageReqBatchRead: // 批量已读
		err = r.msgHasRead(userHid, detail, req.Data)
	case MessageReqBatchDeleted: // 批量删除
		err = r.msgBranchDelete(userHid, detail, req.Data)
	case MessageReqAllRead: // 全部已读
		err = r.msgAllHasRead(userHid, detail, req.Data)
	case MessageReqAllDeleted: // 全部删除
		err = r.msgDelete(userHid, detail, req.Data)
	default:

		var resData interface{}
		// 如果是一般的消息 可以通过定义个函数，将参数传递过去处理
		if resData, err = r.MessageAction(userHid,req.Data); err != nil {
			detail = r.GetFailWithMsg(err.Error())
			// 发送响应
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		} else if resData != nil {
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: r.GetSuccessWithData(resData),
			})
		}
	}
	return
}
func (r *MessageClient) msgHeartBreak(userHid string, detail Resp, dataMessage interface{}) (err error) {
	_ = userHid
	switch dataMessage.(type) {
	case int64, int, float64:
		// 发送心跳
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespHeartBeat,
			Detail: detail,
		})
		return
	default:
		detail = r.GetFailWithMsg(fmt.Sprintf("心跳数据格式data必须为数字类型 (%#v %t)", dataMessage, dataMessage))
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
	}
	return
}
func (r *MessageClient) msgSendNewMsg(userHid string, detail Resp, dataMessage interface{}) (err error) {
	var data PushMessageRequestStruct
	if err = r.Struct2StructByJson(dataMessage, &data); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}

	data.FromUserId = userHid

	// mapColumn:=  data.FieldTrans()
	// if data
	// // 参数校验
	// err = global.NewValidatorError(global.Validate.Struct(data),)

	if !HubMessage.CheckIdempotenceTokenFunc(data.IdempotenceToken) {
		err = fmt.Errorf(IdempotenceTokenInvalidMsg)
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}

	if err = HubMessage.Service.Dao.CreateMessage(&data); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}

	// 刷新条数
	HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)

	// 发送响应
	r.SendChan.SafeSend(MessageWsResponseStruct{
		Type:   MessageRespNormal,
		Detail: detail,
	})
	return
}
func (r *MessageClient) msgHasRead(userHid string, detail Resp, dataMessage interface{}) (err error) {
	_ = userHid
	var data Req
	if err = r.Struct2StructByJson(dataMessage, &data); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}

	if err = HubMessage.Service.Dao.BatchUpdateMessageRead(data.GetIds()); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}
	// 刷新条数
	HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
	// 发送响应
	r.SendChan.SafeSend(MessageWsResponseStruct{
		Type:   MessageRespNormal,
		Detail: detail,
	})
	return
}

func (r *MessageClient) msgBranchDelete(userHid string, detail Resp, dataMessage interface{}) (err error) {
	_ = userHid
	var data Req
	if err = r.Struct2StructByJson(dataMessage, &data); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}

	if err = HubMessage.Service.Dao.BatchUpdateMessageDeleted(data.GetIds()); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}
	// 刷新条数
	HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
	// 发送响应
	r.SendChan.SafeSend(MessageWsResponseStruct{
		Type:   MessageRespNormal,
		Detail: detail,
	})
	return
}
func (r *MessageClient) msgAllHasRead(userHid string, detail Resp, dataMessage interface{}) (err error) {
	_ = dataMessage
	err = HubMessage.Service.Dao.UpdateAllMessageRead(userHid)
	if err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}
	// 刷新条数
	HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
	// 发送响应
	r.SendChan.SafeSend(MessageWsResponseStruct{
		Type:   MessageRespNormal,
		Detail: detail,
	})
	return
}
func (r *MessageClient) msgDelete(userHid string, detail Resp, dataMessage interface{}) (err error) {
	_ = dataMessage
	if err = HubMessage.Service.Dao.UpdateAllMessageDeleted(userHid); err != nil {
		detail = r.GetFailWithMsg(err.Error())
		// 发送响应
		r.SendChan.SafeSend(MessageWsResponseStruct{
			Type:   MessageRespNormal,
			Detail: detail,
		})
		return
	}
	// 刷新条数
	HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
	// 发送响应
	r.SendChan.SafeSend(MessageWsResponseStruct{
		Type:   MessageRespNormal,
		Detail: detail,
	})
	return
}
func (r *MessageClient) Receive() {
	var (
		needBreak bool
		err       error
	)

	defer func() {
		if e := recover(); e != nil {
			if err = r.close(); err != nil {
				return
			}
			_, file, line, _ := runtime.Caller(1)
			logContent := map[string]interface{}{
				"r":    r,
				"loc":  fmt.Sprintf("%s(l:%d)", file, line),
				"desc": fmt.Sprintf("[消息中心][接收端]连接可能已断开"),
			}
			if err != nil {
				logContent["err"] = err.Error()
			}
			r.Context.Error(logContent, "MessageClientErrorHandler")
			return
		}
	}()

	userHid, _, _ := r.GetUserHid()

	for {

		logContent := map[string]interface{}{
			"Key":     r.Key,
			"userHid": userHid,
		}
		var msg []byte
		if _, msg, err = r.Conn.ReadMessage(); err != nil {
			switch err.Error() {
			case "websocket: close 1001 (going away)":
				panic(err)
			}

			logContent["err"] = err.Error()
			logContent["desc"] = fmt.Sprintf("[消息中心][接收端]接收消息异常")
			r.Context.Error(logContent, "MessageClientReceive1")
			continue
		}

		logContent["msgV"] = string(msg)

		// 记录活跃时间
		r.LastActiveTime = time.Now()

		r.RetryCount = 0

		logContent["desc"] = fmt.Sprintf("[消息中心][接收端][%s]接收数据成功", r.Key)

		// 解压数据
		// data := utils.DeCompressStrByZlib(string(msg))
		if needBreak, err = r.receiveMessageAct(msg, userHid); err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "MessageClientReceive1")
			continue
		}

		if !needBreak {
			r.Context.Info(logContent, "MessageClientReceive1")
			continue
		}

		// if err = r.close(); err != nil {
		// 	logContent["closeErr"] = err.Error()
		// 	r.Context.Error(logContent, "MessageClientReceive1")
		// 	return
		// }
		r.Context.Info(logContent, "MessageClientReceive1")

	}
}

func (r *MessageClient) GetUserHid() (res string, user UserInterface, err error) {
	if user, err = r.UserFunc(); err != nil {
		return
	}
	if res, err = user.GetUserHid(); err != nil {
		return
	}
	return
}

// Send 发送数据
func (r *MessageClient) Send() {
	var err error
	// 创建定时器, 超出指定时间间隔, 向前端发送ping消息心跳
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err = r.close()
		if rv := recover(); rv != nil {
			r.Context.Error(map[string]interface{}{
				"key":  r.Key,
				"rv":   rv,
				"desc": fmt.Sprintf("[消息中心][发送]---------------连接可能已断开"),
			}, "MessageClientSend0")
		}
	}()

	for {
		select {

		// 发送通道
		case msg, ok := <-r.SendChan.C:

			// 设定回写超时时间 10 S
			err = r.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// send通道已关闭
				if err = r.writeMessage(websocket.CloseMessage, "closed"); err != nil {
					r.Context.Error(map[string]interface{}{
						"key":  r.Key,
						"ip":   r.Ip,
						"msg":  msg,
						"ok":   ok,
						"desc": fmt.Sprintf("[消息中心][发送]---------------send通道已关闭"),
					}, "MessageClientSend1")
				}
				panic("connection closed")
			}
			var bt []byte
			if msg != nil {
				if bt, err = json.Marshal(msg); err != nil {
					r.Context.Error(map[string]interface{}{
						"key":  r.Key,
						"ip":   r.Ip,
						"err":  err.Error(),
						"msg":  msg,
						"desc": fmt.Sprintf("[消息中心][发送]---------------生成JSON错误"),
					}, "MessageClientSend2")
				}
			}
			// 发送文本消息
			if err = r.writeMessage(websocket.TextMessage, string(bt)); err != nil {
				r.Context.Error(map[string]interface{}{
					"key":  r.Key,
					"ip":   r.Ip,
					"err":  err.Error(),
					"msg":  msg,
					"desc": fmt.Sprintf("[消息中心][发送]---------------发送文本消息异常"),
				}, "MessageClientSend3")
				panic(err)
			}
		// 长时间无新消息
		case <-ticker.C:
			if err = r.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				r.Context.Error(map[string]interface{}{
					"key":  r.Key,
					"ip":   r.Ip,
					"err":  err.Error(),
					"desc": fmt.Sprintf("[消息中心][发送]---------------长时间无新消息错误"),
				}, "MessageClientSend4")
				panic(err)
			}
			// 发送ping消息
			if err = r.writeMessage(websocket.PingMessage, "ping"); err != nil {
				r.Context.Error(map[string]interface{}{
					"key":  r.Key,
					"ip":   r.Ip,
					"err":  err.Error(),
					"desc": fmt.Sprintf("[消息中心][发送]---------------发送ping消息"),
				}, "MessageClientSend5")
				panic(err)
			}
		}

	}
}

// 回写消息
func (r *MessageClient) writeMessage(messageType int, data string) (err error) {

	// 字符串压缩
	// s, _ := utils.CompressStrByZlib(data)

	r.Context.Info(map[string]interface{}{
		"key":         r.Key,
		"data":        data,
		"messageType": messageType,
		"desc":        fmt.Sprintf("[消息中心][发送端]"),
	}, "MessageClientWriteMessage")

	err = r.Conn.WriteMessage(messageType, []byte(data))

	return
}

// 关闭连接
func (r *MessageClient) close() (err error) {
	HubMessage.lock.Lock()
	defer HubMessage.lock.Unlock()

	if _, ok := HubMessage.Clients[r.Key]; ok {

		delete(HubMessage.Clients, r.Key)

		// 关闭发送通道
		r.SendChan.SafeClose()

		userHid, _ := r.UserFunc()

		r.Context.Error(map[string]interface{}{
			"userHid": userHid,
			"ip":      r.Ip,
			"key":     r.Key,
			"desc":    fmt.Sprintf("[消息中心][用户下线]"),
		}, "MessageWriteMessageClose")

	}

	err = r.Conn.Close()
	return
}
