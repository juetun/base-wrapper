package short_message_impl

import "github.com/juetun/base-wrapper/lib/app_obj"

type Sms100 struct {
}

func NewSms100() (r app_obj.ShortMessageInter) {
	return &Sms100{}
}

func (s Sms100) Send(param *app_obj.ShortMessageConfig) (err error) {
	panic("implement me")
}
