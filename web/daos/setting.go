//Package daos
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package daos

import "github.com/juetun/base-wrapper/web/models"

type DaoSetting interface {

	//
	InitBasicField() (err error)

	//
	Slack() (res models.Slack, err error)

	UpdateSlack(url, template string) (err error)

	CreateChannel(channel string) (res int, err error)

	IsChannelExist(channel string) (ok bool, err error)

	RemoveChannel(id int) (res int64, err error)

	Mail() (mail models.Mail, err error)

	UpdateMail(config, template string) (err error)

	CreateMailUser(username, email string) (res int64, err error)

	RemoveMailUser(id int) (res int64, err error)

	WebHook() (webHook models.WebHook, err error)

	UpdateWebHook(url, template string) (err error)
}
