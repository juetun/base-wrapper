package short_message_impl

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/juetun/base-wrapper/lib/base"
	"strings"
)

const (
	LogTypeErr   = "error"
	LogTypeInfo  = "info"
	LogTypeDebug = "debug"
	LogTypeFatal = "Fatal"
)

type AliYunSms struct {
	shortMessageConfig *ShortMessageConfig
	_client            *dysmsapi20170525.Client
}

func (r *AliYunSms) InitClient() (err error) {
	err = r.createClient()
	return
}

var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

func (r *AliYunSms) log(ctx *base.Context, logType, logMark string, logContent map[string]interface{}) {
	switch logType {
	case LogTypeErr:
		ctx.Error(logContent, logMark)
	case LogTypeInfo:
		ctx.Info(logContent, logMark)
	case LogTypeDebug:
		ctx.Debug(logContent, logMark)
	case LogTypeFatal:
		ctx.Fatal(logContent, logMark)
	}
}

func (r *AliYunSms) Send(ctx *base.Context, param *MessageArgument, logTypes ...string) (err error) {
	logContent := make(map[string]interface{}, 10)
	logMark := "AliYunSmsSend"
	var logType = ""
	if len(logTypes) > 0 {
		logType = logTypes[0]
	}

	//记录日志
	defer func() {
		if logType != "" {
			r.log(ctx, logType, logMark, logContent)
		} else if err != nil {
			r.log(ctx, LogTypeErr, logMark, logContent)
		} else {
			r.log(ctx, LogTypeInfo, logMark, logContent)
		}
	}()
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(param.SignName),
		TemplateCode:  tea.String(param.TemplateCode),
		PhoneNumbers:  tea.String(param.Mobile),
		TemplateParam: tea.String(param.Content),
	}
	logContent["sendSmsRequest"] = sendSmsRequest
	runtime := &util.RuntimeOptions{}
	var (
		sendSmsResponse *dysmsapi20170525.SendSmsResponse
	)
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				logContent["errorRecover"] = r.Error()
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		sendSmsResponse, err = r._client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			logContent["SendSmsWithOptions"] = err.Error()
			return
		}

		return
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		logContent["error_info"] = error.Message
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}

		if _, err = util.AssertAsString(error.Message); err != nil {
			return
		}
	}
	return
}

func NewAliYunSms(shortMessageConfig *ShortMessageConfig) (r ShortMessageInter) {
	return &AliYunSms{shortMessageConfig: shortMessageConfig,}
}

// Description:
//
// 使用AK&SK初始化账号Client
//
// @return Client
//
// @throws Exception
func (r *AliYunSms) createClient() (_err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(r.shortMessageConfig.AppKey),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(r.shortMessageConfig.AppSecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String(r.shortMessageConfig.Url)
	var _result = &dysmsapi20170525.Client{}
	if _result, _err = dysmsapi20170525.NewClient(config); _err != nil {
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("ShortMessage init client failure (err:%v)\n", _err.Error()))
		return
	}

	r._client = _result
	return
}

func (r *AliYunSms) GetShortMessageConfig() (shortMessageConfig *ShortMessageConfig) {
	return r.shortMessageConfig
}
