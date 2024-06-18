package short_message_impl

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type Sms100 struct {
	shortMessageConfig *ShortMessageConfig
}

func (s *Sms100) InitClient() (err error) {
	return
}

func NewSms100(shortMessageConfig *ShortMessageConfig) (r ShortMessageInter) {
	return &Sms100{
		shortMessageConfig: shortMessageConfig,
	}
}

func (s *Sms100) Send(ctx *base.Context, param *MessageArgument, logTypes ...string) (err error) {
	fmt.Println("Sms100 发送短信")
	return
}
func (s *Sms100) GetShortMessageConfig() (shortMessageConfig *ShortMessageConfig) {
	return s.shortMessageConfig
}
