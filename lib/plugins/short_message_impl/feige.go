// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
//飞鸽传书 http://www.feige.ee/dev/dev20
package short_message_impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type FeiGe struct {
	Password string `json:"password"`
	SignId   string `json:"sign_id"`
	Url      string `json:"url"` //"http://api.feige.ee/SmsService/Send"
}

func (r *FeiGe) Send(param *app_obj.MessageArgument) (err error) {
	err = r.sendSMS(param.Mobile, param.Content)
	return
}

func NewFeiGe() (r app_obj.ShortMessageInter) {
	return &FeiGe{
		Url:      "http://api.feige.ee/SmsService/Send",
		SignId:   "",
		Password: "",
	}
}

//SMSRspJSON 飞鸽传书返回结构体
type SMSRspJSON struct {
	Code         int    `json:"Code"`
	Message      string `json:"Message"`
	SendId       string `json:"SendId"`
	InvalidCount int    `json:"InvalidCount"`
	SuccessCount int    `json:"SuccessCount"`
	BlackCount   int    `json:"BlackCount"`
}

//sendSMS 发送数据
func (r *FeiGe) sendSMS(mobile, content string) error {
	var formValues = url.Values{}
	formValues.Set("Mobile", mobile)
	formValues.Set("Content", content)
	formValues.Set("Account", mobile)
	formValues.Set("Pwd", r.Password)
	formValues.Set("SignId", r.SignId)
	rsp, err := http.PostForm(r.Url, formValues)
	if err != nil {
		return fmt.Errorf("请求失败")
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("返回body解析失败")
	}

	retJson := &SMSRspJSON{}
	err = json.Unmarshal(body, retJson)
	if err != nil {
		return fmt.Errorf("返回Json解析失败")
	}

	if retJson.Code != 0 {
		return fmt.Errorf("发送失败")
	}

	return err
}
