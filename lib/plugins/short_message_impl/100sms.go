package short_message_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/plugins"
)

type Sms100 struct {
	shortMessageConfig *plugins.ShortMessageConfig
}

func (s *Sms100) InitClient() (err error) {
	return
}

func NewSms100(shortMessageConfig *plugins.ShortMessageConfig) (r plugins.ShortMessageInter) {
	return &Sms100{
		shortMessageConfig: shortMessageConfig,
	}
}

func (s *Sms100) Send(ctx *base.Context, param *plugins.MessageArgument, logTypes ...string) (err error) {
	fmt.Println("Sms100 发送短信")
	return
}
func (s *Sms100) GetShortMessageConfig() (shortMessageConfig *plugins.ShortMessageConfig) {
	return s.shortMessageConfig
}
