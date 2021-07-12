// Package dao_impl
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package dao_impl

import (
	"encoding/json"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/daos"
	"github.com/juetun/base-wrapper/web/models"
)

type DaoSettingImpl struct {
	base.ServiceDao
}

func NewDaoSettingImpl(context ...*base.Context) (res daos.DaoSetting) {
	p := &DaoSettingImpl{}
	p.SetContext(context...)
	return p
}

// InitBasicField 初始化基本字段 邮件、slack等
func (r *DaoSettingImpl) InitBasicField() (err error) {
	setting := models.Setting{
		Code:  models.SlackCode,
		Key:   models.SlackUrlKey,
		Value: "",
	}
	if err = r.Context.Db.Create(&setting).Error; err != nil {
		return
	}

	setting.Id = 0

	setting = models.Setting{
		Code: models.SlackCode,
		Key:  models.SlackTemplateKey,
		Value: `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`,
	}

	if err = r.Context.Db.Create(&setting).Error; err != nil {
		return
	}
	setting = models.Setting{
		Code:  models.MailCode,
		Key:   models.MailServerKey,
		Value: "",
	}
	if err = r.Context.Db.Create(&setting).Error; err != nil {
		return
	}

	setting = models.Setting{
		Code: models.MailCode,
		Key:  models.MailTemplateKey,
		Value: `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`,
	}

	if err = r.Context.Db.Create(&setting).Error; err != nil {
		return
	}
	setting = models.Setting{
		Code: models.WebhookCode,
		Key:  models.WebhookTemplateKey,
		Value: `{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}"`,
	}
	if err = r.Context.Db.Create(&setting).Error; err != nil {
		return
	}
	setting = models.Setting{
		Code:  models.WebhookCode,
		Key:   models.WebhookUrlKey,
		Value: "",
	}
	if err = r.Context.Db.Create(&setting).Error; err != nil {
		return
	}
	return
}

func (r *DaoSettingImpl) Slack() (res models.Slack, err error) {
	list := make([]models.Setting, 0)
	err = r.Context.Db.
		Where("code = ?", models.SlackCode).
		Find(&list).
		Error
	slack := models.Slack{}
	if err != nil {
		return slack, err
	}

	r.formatSlack(list, &slack)

	return slack, err
}

func (r *DaoSettingImpl) formatSlack(list []models.Setting, slack *models.Slack) {
	for _, v := range list {
		switch v.Key {
		case models.SlackUrlKey:
			slack.Url = v.Value
		case models.SlackTemplateKey:
			slack.Template = v.Value
		default:
			slack.Channels = append(slack.Channels, models.Channel{
				Id: v.Id, Name: v.Value,
			})
		}
	}
}

func (r *DaoSettingImpl) UpdateSlack(url, template string) (err error) {

	data := map[string]interface{}{
		"code": models.SlackCode,
		"key":  models.SlackUrlKey,
	}
	if err = r.Context.Db.Where("value=?", url).
		Updates(data).Error; err != nil {
		return
	}
	data = map[string]interface{}{
		"code": models.SlackCode,
		"key":  models.SlackTemplateKey,
	}
	if err = r.Context.Db.Where("value=?", template).
		Updates(data).Error; err != nil {
		return
	}

	return nil
}

// CreateChannel 创建slack渠道
func (r *DaoSettingImpl) CreateChannel(channel string) (res int, err error) {
	setting := models.Setting{
		Code:  models.SlackCode,
		Key:   models.SlackChannelKey,
		Value: channel,
	}
	err = r.Context.Db.Create(&setting).Error
	if err == nil {
		res = setting.Id
	}
	return
}

func (r *DaoSettingImpl) IsChannelExist(channel string) (ok bool, err error) {

	setting := models.Setting{
		Code:  models.SlackCode,
		Key:   models.SlackChannelKey,
		Value: channel,
	}
	var count int64
	if err = r.Context.Db.Where(&setting).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		ok = true
	}
	return
}

// RemoveChannel 删除slack渠道
func (r *DaoSettingImpl) RemoveChannel(id int) (res int64, err error) {

	setting := models.Setting{
		Id:   id,
		Code: models.SlackCode,
		Key:  models.SlackChannelKey,
	}
	err = r.Context.Db.Delete(&setting).Error
	return
}

// Mail region 邮件配置
func (r *DaoSettingImpl) Mail() (mail models.Mail, err error) {
	list := make([]models.Setting, 0)
	err = r.Context.Db.Where("code = ?", models.MailCode).
		Find(&list).Error
	mail = models.Mail{MailUsers: make([]models.MailUser, 0)}
	if err != nil {
		return
	}
	r.formatMail(list, &mail)
	return
}

func (r *DaoSettingImpl) formatMail(list []models.Setting, mail *models.Mail) {
	mailUser := models.MailUser{}
	for _, v := range list {
		switch v.Key {
		case models.MailServerKey:
			_ = json.Unmarshal([]byte(v.Value), mail)
		case models.MailUserKey:
			_ = json.Unmarshal([]byte(v.Value), &mailUser)
			mailUser.Id = v.Id
			mail.MailUsers = append(mail.MailUsers, mailUser)
		case models.MailTemplateKey:
			mail.Template = v.Value
		}

	}
}

func (r *DaoSettingImpl) UpdateMail(config, template string) (err error) {
	data := map[string]interface{}{
		"code": models.MailCode,
		"key":  models.MailServerKey,
	}
	if err = r.Context.Db.Where("value=?", config).
		Updates(data).Error; err != nil {
		return
	}
	data = map[string]interface{}{
		"code": models.MailCode,
		"key":  models.MailTemplateKey,
	}
	if err = r.Context.Db.Where("value=?", template).
		Updates(data).Error; err != nil {
		return
	}
	return
}

func (r *DaoSettingImpl) CreateMailUser(username, email string) (res int64, err error) {
	mailUser := models.MailUser{Id: 0, Username: username, Email: email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		return 0, err
	}
	setting := models.Setting{
		Code:  models.MailCode,
		Key:   models.MailUserKey,
		Value: string(jsonByte),
	}

	err = r.Context.Db.Create(setting).Error
	return
}

func (r *DaoSettingImpl) RemoveMailUser(id int) (res int64, err error) {

	setting := models.Setting{
		Id:   id,
		Code: models.MailCode,
		Key:  models.MailUserKey,
	}
	err = r.Context.Db.Delete(&setting).Error
	return
}

func (r *DaoSettingImpl) WebHook() (webHook models.WebHook, err error) {
	list := make([]models.Setting, 0)
	err = r.Context.Db.Where("code = ?", models.WebhookCode).
		Find(&list).
		Error
	webHook = models.WebHook{}
	if err != nil {
		return
	}

	r.formatWebhook(list, &webHook)

	return
}

func (r *DaoSettingImpl) formatWebhook(list []models.Setting, webHook *models.WebHook) {
	for _, v := range list {
		switch v.Key {
		case models.WebhookUrlKey:
			webHook.Url = v.Value
		case models.WebhookTemplateKey:
			webHook.Template = v.Value
		}

	}
}

func (r *DaoSettingImpl) UpdateWebHook(url, template string) (err error) {

	data := map[string]interface{}{
		"code": models.WebhookCode,
		"key":  models.WebhookUrlKey,
	}
	if err = r.Context.Db.Where("value=?", url).
		Updates(data).Error; err != nil {
		return
	}
	data = map[string]interface{}{
		"code": models.WebhookCode,
		"key":  models.WebhookTemplateKey,
	}
	if err = r.Context.Db.Where("value=?", template).
		Updates(data).Error; err != nil {
		return
	}

	return
}
