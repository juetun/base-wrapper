// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
// 二维码生成封装
package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/skip2/go-qrcode"
)

type QrCodeParams struct {
	Width         int    `json:"width"`   //二维码的宽
	Content       string `json:"content"` //二维码的内容
	TargetImgPath string `json:"target_img_path"`
}

func NewQrCodeParams() (res *QrCodeParams) {
	res = &QrCodeParams{}
	return
}

//生成二维码信息
func (r *QrCodeParams) CreateQrCodeToFile() (err error) {
	if r.TargetImgPath == "" {
		err = fmt.Errorf("the params target_img_path is nil")
		return
	}
	err = qrcode.WriteFile(r.Content, qrcode.Medium, r.Width, r.TargetImgPath)
	//err := qrcode.WriteColorFile(arg.Content, qrcode.Medium, arg.Width, color.Black, color.White, "qr.png")
	return
}

//生成二维码的base64格式
func (r *QrCodeParams) CreateQrCodeAsBase64Code() (res string, err error) {
	var png []byte
	if png, err = qrcode.Encode(r.Content, qrcode.Medium, r.Width); err != nil {
		return
	}
	res = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(png))
	return
}
