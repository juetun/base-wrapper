package short_message_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

type Sms100 struct {
	shortMessageConfig *app_obj.ShortMessageConfig
}

func (s *Sms100) InitClient() {
	return
}

func NewSms100(shortMessageConfig *app_obj.ShortMessageConfig) (r app_obj.ShortMessageInter) {
	return &Sms100{
		shortMessageConfig: shortMessageConfig,
	}
}

func (s *Sms100) Send(param *app_obj.MessageArgument) (err error) {
	fmt.Println("Sms100 发送短信")
	return
}
func (s *Sms100) GetShortMessageConfig(param *app_obj.MessageArgument) (shortMessageConfig *app_obj.ShortMessageConfig) {
	return s.shortMessageConfig
}
