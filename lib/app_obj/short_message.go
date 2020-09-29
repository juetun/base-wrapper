package app_obj

import (
	"fmt"
	"sync"
)

var ShortMessageObj *shortMessage

//短信发送的参数
type ShortMessageConfig struct {
	Mobile  string `json:"mobile"`  //手机号
	Area    string `json:"area"`    //地区号 默认 86
	Content string `json:"content"` //短信内容

	ExceptChannel []string `json:"except_channel"` //排除渠道，（此字段主要为当某一渠道发送不成功后，重试发送切换渠道使用）
	Channel       string   `json:"channel"`        //短信渠道号 不设置使用默认规则
}

//渠道发送需要实现的接口
type ShortMessageInter interface {
	Send(param *ShortMessageConfig) (err error)
}

func NewShortMessage(channelMap ...map[string]ShortMessageInter) (res *shortMessage) {
	res = &shortMessage{
		channelListHandler: map[string]ShortMessageInter{},
	}
	for _, item := range channelMap {
		for key, value := range item {
			res.channelListHandler[key] = value
		}
	}
	return
}

//添加渠道
//channelName string 渠道名称
//channel 渠道实现的调用的结构体
func AddMessageChannel(channelName string, channel ShortMessageInter) {
	var syc sync.Mutex
	syc.Lock()
	defer syc.Unlock()
	ShortMessageObj.channelListHandler[channelName] = channel
}

//短息发送调用的公共动作
type shortMessage struct {
	channelListHandler map[string]ShortMessageInter //当前支持的短息通道列表
	shortMessageIndex  int                          //当前发送短信的序号
}

//发送短信调用接口
func (r *shortMessage) SendMsg(param *ShortMessageConfig) (channelName string, err error) {

	if len(r.channelListHandler) == 0 {
		err = fmt.Errorf("当前没有可发送短信的通道")
		return
	}
	channelData, channelName, err := r.initChannel(param)
	if err != nil {
		return
	}
	err = channelData.Send(param)
	return
}

//获取渠道列表
func (r *shortMessage) GetChannelKey() (res []string) {
	res = make([]string, 0, len(r.channelListHandler))
	for key, _ := range r.channelListHandler {
		res = append(res, key)
	}
	return
}

//获取当前可选的短信通道
func (r *shortMessage) getChannelListHandler(param *ShortMessageConfig) (channelListHandler map[string]ShortMessageInter) {
	channelListHandler = make(map[string]ShortMessageInter, len(r.channelListHandler))
	if len(param.ExceptChannel) > 0 {
		var f bool
		for key, value := range r.channelListHandler {
			for _, channel := range param.ExceptChannel {
				if channel == key {
					f = true
					break
				}
			}
			if f {
				continue
			}
			channelListHandler[key] = value
		}
	} else {
		channelListHandler = r.channelListHandler
	}
	return
}

func (r *shortMessage) initChannel(param *ShortMessageConfig) (channelData ShortMessageInter, name string, err error) {

	if param.Channel == "" { //多个短信通道轮流发
		r.shortMessageIndex++
		if r.shortMessageIndex > 1000 {
			r.shortMessageIndex = 0
		}
		ind := r.shortMessageIndex % len(r.channelListHandler)

		channelListHandler := r.getChannelListHandler(param)
		i := 0
		for chanelName, value := range channelListHandler {

			if ind == i {
				channelData = value
				name = chanelName
			}
			i++
		}
		return
	}

	if _, ok := r.channelListHandler[param.Channel]; !ok {
		err = fmt.Errorf("当前不支持你选择的短信发送通道(%s)", param.Channel)
	}

	return
}
