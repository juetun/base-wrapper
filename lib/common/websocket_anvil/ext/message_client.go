package ext

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juetun/base-wrapper/lib/base"
)

// MessageClient 消息客户端
type MessageClient struct {
	WebsocketBaseHandler `json:"-"`

	Context *base.Context `json:"-"`

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

// Contains 判断uint数组是否包含item元素
func (r *MessageClient) Contains(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
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

		_ = userData
		// go func() {
		// 	HubMessage.RefreshUserMessage.SafeSend([]string{userHid})
		// }()
		//
		// // 广播当前用户上线
		// // 通知除自己之外的人
		// go HubMessage.Broadcast.SafeSend(MessageBroadcast{
		// 	MessageWsResponseStruct: MessageWsResponseStruct{
		// 		Type: MessageRespOnline,
		// 		Detail: r.GetSuccessWithData(map[string]interface{}{
		// 			"user": userData,
		// 		}),
		// 	},
		// 	UserIds: r.ContainsThenRemove(HubMessage.UserIds, userHid),
		// })

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

func (r *MessageClient) Receive() {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = r.close()
			r.Context.Error(map[string]interface{}{
				"key":  r.Key,
				"e":    e,
				"desc": fmt.Sprintf("[消息中心][接收端]连接可能已断开"),
			}, "MessageClientReceive0")
		}
	}()
	userHid, _, _ := r.GetUserHid()
	for {

		var msg []byte
		if _, msg, err = r.Conn.ReadMessage(); err != nil {
			r.Context.Error(map[string]interface{}{
				"msgV": msg,
				"err":  err.Error(),
				"key":  r.Key,
				"desc": fmt.Sprintf("[消息中心][接收端]接收消息异常"),
			}, "MessageClientReceive1")
			return
		}

		// 记录活跃时间
		r.LastActiveTime = time.Now()
		r.RetryCount = 0
		if err != nil {
			panic(err)
		}
		// 解压数据
		// data := utils.DeCompressStrByZlib(string(msg))

		r.Context.Info(map[string]interface{}{
			"msgV":    string(msg),
			"userHid": userHid,
			"desc":    fmt.Sprintf("[消息中心][接收端][%s]接收数据成功", r.Key),
		}, "MessageClientReceive1")

		// 数据转为json
		var req MessageWebSocketRequestStruct
		if err = json.Unmarshal(msg, &req); err != nil {
			r.Context.Error(map[string]interface{}{
				"msg":  msg,
				"err":  err.Error(),
				"desc": fmt.Sprintf("[Json2Struct]转换异常"),
			}, "MessageClientReceive2")
		}
		switch req.Type {
		case MessageReqHeartBeat: // 心跳消息
			if _, ok := req.Data.(float64); ok {
				// 发送心跳
				r.SendChan.SafeSend(MessageWsResponseStruct{
					Type:   MessageRespHeartBeat,
					Detail: r.GetSuccess(),
				})
			}
		case MessageReqPush: // 推送新消息
			var data PushMessageRequestStruct
			err = r.Struct2StructByJson(req.Data, &data)

			// mapColumn:=  data.FieldTrans()
			// if data
			// // 参数校验
			// err = global.NewValidatorError(global.Validate.Struct(data),)
			detail := r.GetSuccess()
			if err == nil {
				if !HubMessage.CheckIdempotenceTokenFunc(data.IdempotenceToken) {
					err = fmt.Errorf(IdempotenceTokenInvalidMsg)
				} else {
					data.FromUserId = userHid
					err = HubMessage.Service.Dao.CreateMessage(&data)
				}
			}
			if err != nil {
				detail = r.GetFailWithMsg(err.Error())
			} else {
				// 刷新条数
				HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
			}
			// 发送响应
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqBatchRead: // 批量已读
			var data Req
			err = r.Struct2StructByJson(req.Data, &data)
			err = HubMessage.Service.Dao.BatchUpdateMessageRead(data.GetIds())
			detail := r.GetSuccess()
			if err != nil {
				detail = r.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
			// 发送响应
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqBatchDeleted: // 批量删除
			var data Req
			err = r.Struct2StructByJson(req.Data, &data)
			err = HubMessage.Service.Dao.BatchUpdateMessageDeleted(data.GetIds())
			detail := r.GetSuccess()
			if err != nil {
				detail = r.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
			// 发送响应
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqAllRead: // 全部已读
			err = HubMessage.Service.Dao.UpdateAllMessageRead(userHid)
			detail := r.GetSuccess()
			if err != nil {
				detail = r.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
			// 发送响应
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqAllDeleted: // 全部删除
			err = HubMessage.Service.Dao.UpdateAllMessageDeleted(userHid)
			detail := r.GetSuccess()
			if err != nil {
				detail = r.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			HubMessage.RefreshUserMessage.SafeSend(HubMessage.UserIds)
			// 发送响应
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		default:
			r.SendChan.SafeSend(MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: r.GetFailWithMsg(fmt.Sprintf("当前不支持您选择的消息类型（%s）", req.Type)),
			})
		}
	}
}

// ContainsIndex 判断uint数组是否包含item元素, 返回index
func (r *MessageClient) ContainsIndex(arr []string, item string) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}

// ContainsThenRemove 判断uint数组是否包含item元素, 并移除
func (r *MessageClient) ContainsThenRemove(arr []string, item string) []string {
	index := r.ContainsIndex(arr, item)
	if index >= 0 {
		arr = append(arr[:index], arr[index+1:]...)
	}
	return arr
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
				"desc": fmt.Sprintf("[消息中心][发送端]连接可能已断开"),
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
						"desc": fmt.Sprintf("send通道已关闭"),
					}, "MessageClientSend1")
				}
				panic("connection closed")
			}
			var bt []byte
			if msg != nil {
				if bt, err = json.Marshal(msg); err != nil {
					r.Context.Error(map[string]interface{}{
						"key": r.Key,
						"ip":  r.Ip,
						"err": err.Error(),
						"msg": msg,
					}, "MessageClientSend2")
				}
			}
			// 发送文本消息
			if err = r.writeMessage(websocket.TextMessage, string(bt)); err != nil {
				r.Context.Error(map[string]interface{}{
					"key": r.Key,
					"ip":  r.Ip,
					"err": err.Error(),
					"msg": msg,
				}, "MessageClientSend3")
				panic(err)
			}
		// 长时间无新消息
		case <-ticker.C:
			if err = r.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				r.Context.Error(map[string]interface{}{
					"key": r.Key,
					"ip":  r.Ip,
					"err": err.Error(),
				}, "MessageClientSend4")
				panic(err)
			}
			// 发送ping消息
			if err = r.writeMessage(websocket.PingMessage, "ping"); err != nil {
				r.Context.Error(map[string]interface{}{
					"key": r.Key,
					"ip":  r.Ip,
					"err": err.Error(),
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
