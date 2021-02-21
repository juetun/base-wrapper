/**
* @Author:changjiang
* @Description:
* @File:identifying_code
* @Version: 1.0.0
* @Date 2021/2/22 12:16 上午
 */
package identifying_code_pkg

import (
	"fmt"

	"github.com/mojocn/base64Captcha"
)

// 使用实例
// var idKey, base64Captcha string
// var err error
// // 生成验证码逻辑
// idKey, base64Captcha, err = NewIdentifyingCode(Context(&CustomizeRdsStore{
// 	// 参数。..
// 	RedisClient: //验证码内容存储的位置用
// 	Context   : //记日志用
// }),
// 	CaptchaType(""), // 验证码类型
//
// ).CreateAndGetImgBase64Message()
//
// // 校验逻辑
// NewIdentifyingCode(Context(&CustomizeRdsStore{
// 	// 参数
// })).Context.Verify(idKey, "anwser", true)
// fmt.Println(idKey, base64Captcha, err)
// 验证码对象操作生成器
// 验证码生成逻辑
func NewIdentifyingCode(identifyingCodeHandler ...IdentifyingCodeHandler) (res *IdentifyingCode) {
	res = &IdentifyingCode{}
	for _, handler := range identifyingCodeHandler {
		handler(res)
	}
	return
}

type IdentifyingCode struct {
	Id            string
	CaptchaType   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit

	Context *CustomizeRdsStore
}

type IdentifyingCodeHandler func(arg *IdentifyingCode)

func Id(id string) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.Id = id
	}
}

// CaptchaType
func CaptchaType(captchaType string) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.CaptchaType = captchaType
		return
	}
}

// DriverAudio
func DriverAudio(driverAudio *base64Captcha.DriverAudio) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.DriverAudio = driverAudio
		return
	}
}

// DriverString
func DriverString(driverString *base64Captcha.DriverString) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.DriverString = driverString
		return
	}
}

// DriverChinese
func DriverChinese(driverChinese *base64Captcha.DriverChinese) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.DriverChinese = driverChinese
		return
	}
}

// DriverMath
func DriverMath(driverMath *base64Captcha.DriverMath) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.DriverMath = driverMath
		return
	}
}
func DriverDigit(driverDigit *base64Captcha.DriverDigit) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.DriverDigit = driverDigit
		return
	}
}

func Context(context *CustomizeRdsStore) (res IdentifyingCodeHandler) {
	return func(arg *IdentifyingCode) {
		arg.Context = context
		return
	}
}

// 创建验证码信息并获取base64格式
func (r *IdentifyingCode) CreateAndGetImgBase64Message() (idKeyD string,
	base64stringD string,
	err error) {

	// 判断类型是否支持
	if err = r.Context.FlagType(r.CaptchaType); err != nil {
		return
	}

	var driver base64Captcha.Driver
	switch r.CaptchaType {
	case "audio":
		driver = r.DriverAudio
	case "string":
		driver = r.DriverString.ConvertFonts()
	case "math":
		driver = r.DriverMath.ConvertFonts()
	case "chinese": // 中文
		driver = r.DriverChinese.ConvertFonts()
	default: // 默认数字验证码
		driver = r.DriverDigit
	}
	if driver == nil {
		err = fmt.Errorf("您没有设置验证码生成的必须参数")
		r.Context.Context.Error(map[string]interface{}{
			"IdentifyingCode": "IdentifyingCode.CreateAndGetImgBase64Message",
		}, err.Error())
		return
	}
	if idKeyD, base64stringD, err = base64Captcha.NewCaptcha(driver, r.Context).Generate(); err != nil {
		r.Context.Context.Error(map[string]interface{}{
			"err":             err,
			"IdentifyingCode": "IdentifyingCode.flagType",
		})
	}
	return

}
