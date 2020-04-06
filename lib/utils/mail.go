package utils

/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-04
 * Time: 21:58
 */
import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"

	"github.com/juetun/base-wrapper/lib/app_log"
)

const MAIlTYPE = "html"

type EmailType string

type EmailParam struct {
	User        EmailType `json:"user"`
	Password    EmailType `json:"password"`
	Host        EmailType `json:"host"`
	To          EmailType `json:"to"`
	Subject     EmailType `json:"subject"`
	Body        EmailType `json:"body"`
	MailType    EmailType `json:"mail_type"`
	Description EmailType `json:"description"`
	Attaches    map[string]string
}

var mailParam *EmailParam

var mailAddr string

type EM func(*EmailParam) (interface{}, error)

func (et EmailType) CheckIsNull() error {
	if string(et) == "" {
		app_log.GetLog().Errorln("message", "value can not be null")
		return errors.New("value can not be null")
	}
	return nil
}

func (ep *EmailParam) SetMailUser(user EmailType) EM {
	return func(e *EmailParam) (interface{}, error) {
		u := e.User
		err := user.CheckIsNull()
		if err != nil {
			return nil, err
		}
		e.User = user
		return u, nil
	}
}

func (ep *EmailParam) SetMailPwd(pwd EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		p := ep.Password
		err := pwd.CheckIsNull()
		if err != nil {
			return nil, err
		}
		ep.Password = pwd
		return p, nil
	}
}

func (et EmailType) IsRight() error {
	arr := strings.Split(string(et), ":")
	if len(arr) != 2 {
		app_log.GetLog().Errorln("may be is not semicolon")
		return errors.New("may be is not semicolon")
	}
	mailAddr = arr[0]
	return nil
}

func (ep *EmailParam) SetMailHost(host EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		h := ep.Host
		err := host.CheckIsNull()
		if err != nil {
			return nil, err
		}
		err = host.IsRight()
		if err != nil {
			return nil, err
		}
		ep.Host = host
		return h, nil
	}
}

func (ep *EmailParam) SetMailType(types EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		ty := ep.MailType
		err := types.CheckIsNull()
		if err != nil {
			return nil, err
		}
		ep.MailType = ty
		return ty, nil
	}
}

func (ep *EmailParam) MailInit(options ...EM) (*EmailParam, error) {
	q := &EmailParam{
		MailType: MAIlTYPE,
	}
	for _, option := range options {
		_, err := option(q)
		if err != nil {
			return nil, err
		}
	}
	mailParam = q
	return q, nil
}

func (ep *EmailParam) SetSubject(s EmailType) *EmailParam {
	ep.Subject = s
	return ep
}

func (ep *EmailParam) SetDescription(de EmailType) *EmailParam {
	ep.Description = de
	return ep
}

func (ep *EmailParam) SetAttaches(a map[string]string) *EmailParam {
	ep.Attaches = a
	return ep
}

func (ep *EmailParam) SetBody(b EmailType) *EmailParam {
	ep.Body = b
	return ep
}

func (ep *EmailParam) SetTo(to EmailType) *EmailParam {
	ep.To = to
	return ep
}

func (ep *EmailParam) SendMail2(to string) error {
	sendTo := strings.Split(to, ";")

	subject := ep.Subject
	boundary := "next message"

	mime := bytes.NewBuffer(nil)
	user := string(ep.User)
	password := string(ep.Password)
	host := string(ep.Host)
	// 设置邮件
	mime.WriteString(fmt.Sprintf("From: %s<%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\n", user, user, to, subject))

	mime.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	mime.WriteString("Content-Description: " + string(ep.Description) + "\r\n")

	// 邮件普通Text正文
	mime.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	mime.WriteString("Content-Type: text/plain; charset=utf-8\r\n")

	mime.WriteString(fmt.Sprintf("\n--%s\r\n", boundary))

	boundaryHtml := "boundaryHtml"
	mime.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", boundaryHtml))
	mime.WriteString("Content-Description: Message in alternative text and HTML forms\r\n")
	mime.WriteString(fmt.Sprintf("\n--%s\r\n", boundaryHtml))
	mime.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n", "text/html"))
	mime.WriteString(string(ep.Body))
	mime.WriteString(fmt.Sprintf("\n--%s--\r\n\r\n", boundaryHtml))

	fmt.Println(ep.Subject, ep.Attaches, ep.Description, ep.Body, sendTo, host)

	for k, v := range ep.Attaches {
		attaFile := v
		attaFileName := k
		mime.WriteString(fmt.Sprintf("\n--%s\r\n", boundary))
		mime.WriteString("Content-Type: application/octet-stream\r\n")
		mime.WriteString("Content-Description: 附件\r\n")
		mime.WriteString("Content-Transfer-Encoding: base64\r\n")
		mime.WriteString("Content-Disposition: attachment; filename=\"" + attaFileName + "\"\r\n\r\n")

		// 读取并编码文件内容
		attaData, err := ioutil.ReadFile(attaFile)
		if err != nil {
			return err
		}
		b := make([]byte, base64.StdEncoding.EncodedLen(len(attaData)))
		base64.StdEncoding.Encode(b, attaData)
		mime.Write(b)
	}

	app_log.GetLog().Infoln("message", "mail to the last")
	mime.WriteString("\r\n--" + boundary + "--\r\n\r\n")
	auth := smtp.PlainAuth("", user, password, mailAddr)
	err := smtp.SendMail(host, auth, user, sendTo, mime.Bytes())
	app_log.GetLog().Infoln("message", "mail to the last", "last", err)
	return err
}

func SendMail(to string, subject string, body string) error {
	user := string(mailParam.User)
	password := string(mailParam.Password)
	host := string(mailParam.Host)
	auth := smtp.PlainAuth("", user, password, mailAddr)
	var contentType string
	if mailParam.MailType == "html" {
		contentType = "Content-Type: text/html; charset=UTF-8"
		body = "<html><body>" + body + "</body></html>"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	msg = []byte(subject + contentType + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	app_log.GetLog().Infoln("message", "SendMail", "last", err)
	return err
}
