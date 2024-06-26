package short_message_impl

import (
	"fmt"
	"os"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

const (
	ShortMessageSmsAliYun = "aliyunsms"
	ShortMessageSmsFeiGe  = "feige"
	ShortMessageSms100Sms = "100sms"
)

var (
	ShortMessageSmsOptions = base.ModelItemOptions{
		{Label: "阿里云短信", Value: ShortMessageSmsAliYun},
		{Label: "飞鸽短信", Value: ShortMessageSmsFeiGe},
		{Label: "短信100", Value: ShortMessageSms100Sms},
	}
)

type (
	// 渠道发送需要实现的接口
	ShortMessageInter interface {
		Send(ctx *base.Context, param *MessageArgument, logTypes ...string) (err error)
		InitClient() (err error)
		GetShortMessageConfig() (shortMessageConfig *ShortMessageConfig)
	}
	MessageArgument struct {
		Mobile             string              `json:"mobile"`         // 手机号
		AreaCode           string              `json:"area"`           // 地区号 默认 86
		Content            string              `json:"content"`        // 短信内容
		ExceptChannel      []string            `json:"except_channel"` // 排除渠道，（此字段主要为当某一渠道发送不成功后，重试发送切换渠道使用）
		Channel            string              `json:"channel"`        // 短信渠道号 不设置使用默认规则 不传值调用app_obj.ShortMessageObj.GetSendChannel()方法可用随机获取一个短信发送渠道
		TemplateCode       string              `json:"template_code"`  // 短信模版CODE （阿里云短信用）
		SignName           string              `json:"sign_name"`      // 签名名称 （阿里云短信用）
		Type               int                 `json:"type"`           // 验证码发送的位置的KEY
		ShortMessageConfig *ShortMessageConfig `json:"short_message_config"`
	}
	ShortMessageConfig struct {
		Url                       string `json:"url" yml:"url"`        //请求地址
		AppKey                    string `json:"app_key" yml:"appkey"` //
		AppSecret                 string `json:"app_secret" yml:"appsecret"`
		AliYunAuthName            string `json:"aliyun_auth_name" yml:"aliyunauthname"`
		AliYunAuthTemplateCode    string `json:"aliyun_auth_template_code" yml:"aliyunauthtemplatecode"`
		AliYunGeneralName         string `json:"aliyun_general_name" yml:"aliyungeneralname"`
		AliYunGeneralTemplateCode string `json:"aliyun_general_template_code" yml:"aliyungeneraltemplatecode"`
	}
)

// PluginShortMessage 短信插件初始化
func PluginShortMessage(arg *app_start.PluginsOperate) (err error) {
	_ = arg
	io.SystemOutPrintln("Load short message start")
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	// 数据库配置信息存储对象
	var (
		conShortMsgConfigs    = make(map[string]ShortMessageConfig)
		itemConfig            ShortMessageConfig
		configs               ShortMessageAppConfig
		filePath, connectName string
		yamlFile              []byte
		ok                    bool
	)

	filePath = common.GetConfigFilePath("message.yaml")
	if yamlFile, err = os.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &configs); err != nil {
		io.SystemOutFatalf("load short message config file(%v) err(%+v) \n", filePath, err)
		return
	}
	app_obj.DistributedShortMessageConnects = append(app_obj.DistributedShortMessageConnects, configs.DistributedConnects...)
	//读取common_config配置文件中的信息
	filePath = common.GetCommonConfigFilePath("message.yaml", true)
	if yamlFile, err = os.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
		return
	}
	if err = yaml.Unmarshal(yamlFile, &conShortMsgConfigs); err != nil {
		io.SystemOutFatalf("load short message config file(%v) err(%+v) \n", filePath, err)
		return
	}

	// 初始化短信通道
	shortHandle := map[string]ShortMessageInter{}

	for _, connectName = range configs.Connects {
		if connectName == "" {
			continue
		}
		if itemConfig, ok = conShortMsgConfigs[connectName]; !ok {
			err = fmt.Errorf("当前common_config中不支持您要使用的短信通道连接(%v)", connectName)
			io.SystemOutFatalf("load short message  config err(%+v) \n", err)
			return
		}
		io.SystemOutPrintf("【 %s 】%+v \n", connectName, itemConfig.ToString())
		shortHandle[connectName], err = initShortMessage(connectName, &itemConfig)
	}

	ShortMessageObj = NewShortMessage(shortHandle)
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("ShortMessage load config finished \n"))
	return

}

func initShortMessage(nameSpace string, shortMessageConfig *ShortMessageConfig) (res ShortMessageInter, err error) {
	switch nameSpace { // 短信通道配置 结构体映射
	case ShortMessageSms100Sms:
		res = NewSms100(shortMessageConfig)
		err = res.InitClient()
	case ShortMessageSmsFeiGe:
		res = NewFeiGe(shortMessageConfig)
		err = res.InitClient()
	case ShortMessageSmsAliYun:
		res = NewAliYunSms(shortMessageConfig)
		err = res.InitClient()
	default:
		err = fmt.Errorf("当前不支持此短信通道(%s)", nameSpace)
	}
	return
}

func (r *ShortMessageConfig) ToString() (res string) {
	res = fmt.Sprintf("Url:%s ,AppKey:%s,AppSecret:%v", r.Url, r.AppKey, r.AppSecret)
	return
}

// 添加渠道
// channelName string 渠道名称
// channel 渠道实现的调用的结构体
func AddMessageChannel(channelName string, channel ShortMessageInter) {
	var syc sync.Mutex
	syc.Lock()
	defer syc.Unlock()
	ShortMessageObj.channelListHandler[channelName] = channel
}

// 发送短信调用接口
func (r *shortMessage) GetSendChannel(param *MessageArgument) (channelName string, err error) {

	if len(r.channelListHandler) == 0 {
		err = fmt.Errorf("当前没有可发送短信的通道")
		return
	}
	_, channelName, err = r.initChannel(param)
	return
}

// 发送短信调用接口
func (r *shortMessage) SendMsg(ctx *base.Context, param *MessageArgument, logTypes ...string) (channelName string, err error) {

	if len(r.channelListHandler) == 0 {
		err = fmt.Errorf("当前没有可发送短信的通道")
		return
	}
	channelData, channelName, err := r.initChannel(param)
	if err != nil {
		return
	}
	if channelData == nil {
		err = fmt.Errorf("未获取到发送短信的通道信息")
		return
	}
	err = channelData.Send(ctx, param, logTypes...)
	return
}

var ShortMessageObj = &shortMessage{channelListHandler: map[string]ShortMessageInter{},}

func NewShortMessage(channelMap ...map[string]ShortMessageInter) (res *shortMessage) {
	res = &shortMessage{
		channelListHandler: map[string]ShortMessageInter{},
	}
	for _, item := range channelMap {
		for key, value := range item {
			res.channelListHandler[key] = value
		}
	}
	return
}

// 短息发送调用的公共动作
type shortMessage struct {
	channelListHandler map[string]ShortMessageInter // 当前支持的短息通道列表
	shortMessageIndex  int                          // 当前发送短信的序号
	syncMutex          sync.Mutex
}

// 短信发送的参数
type (
	ShortMessageAppConfig struct {
		Connects            []string `json:"connects" yaml:"connects"`                         //当前应用使用了的数据库连接
		Default             string   `json:"default"  yaml:"default"`                          //默认数据库
		DistributedConnects []string `json:"distributed_connects" yaml:"distributed_connects"` //需要使用的分布式数据库连接
	}
)

// 获取短信渠道列表
func (r *shortMessage) GetChannelKey() (res []string) {
	res = make([]string, 0, len(r.channelListHandler))
	for key := range r.channelListHandler {
		res = append(res, key)
	}
	return
}

// 获取当前可选的短信通道
func (r *shortMessage) getChannelListHandler(param *MessageArgument) (channelListHandler map[string]ShortMessageInter) {
	channelListHandler = make(map[string]ShortMessageInter, len(r.channelListHandler))

	// 如果没有排除的通道，说明按照系统默认的算法选择通道发送
	if len(param.ExceptChannel) <= 0 {
		channelListHandler = r.channelListHandler
		return
	}
	// 将黑名单短信通道排除
	for key, value := range r.channelListHandler {
		if !r.flagExceptChannel(param.ExceptChannel, key) {
			channelListHandler[key] = value
		}
	}
	return
}

// 判断指定通道是否为 黑名单通道
func (r *shortMessage) flagExceptChannel(exceptChannel []string, channelName string) (res bool) {
	for _, channel := range exceptChannel {
		if channel == channelName {
			res = true
			break
		}
	}
	return
}

func (r *shortMessage) upIndex() {
	r.syncMutex.Lock()
	// 多个短信通道轮流发
	r.shortMessageIndex++
	if r.shortMessageIndex > 1000 {
		r.shortMessageIndex = 0
	}
	r.syncMutex.Unlock()
}

func (r *shortMessage) hasSetChannel(param *MessageArgument) (channelData ShortMessageInter, name string, err error) {
	if _, ok := r.channelListHandler[param.Channel]; !ok {
		err = fmt.Errorf("当前不支持你选择的短信发送通道(%s)", param.Channel)
		return
	}
	channelListHandler := r.getChannelListHandler(param)
	for chanelName, value := range channelListHandler {

		//如果已经设置了chanelName
		if param.Channel != "" {
			if param.Channel == chanelName {
				channelData = value
				config := value.GetShortMessageConfig()
				param.ShortMessageConfig = config
				name = chanelName
			}
			continue
		}
	}
	return
}

func (r *shortMessage) hasNotSetChannel(param *MessageArgument) (channelData ShortMessageInter, name string, err error) {
	var (
		i, ind int
	)

	// 更新轮询条件
	r.upIndex()
	ind = r.shortMessageIndex % len(r.channelListHandler)
	i = 0

	channelListHandler := r.getChannelListHandler(param)
	for chanelName, value := range channelListHandler {
		if ind == i {
			channelData = value
			config := value.GetShortMessageConfig()
			param.ShortMessageConfig = config
			name = chanelName
		}
		i++

	}
	return
}

func (r *shortMessage) initChannel(param *MessageArgument) (channelData ShortMessageInter, name string, err error) {

	if param.Channel != "" {
		channelData, name, err = r.hasSetChannel(param)
		return
	}
	channelData, name, err = r.hasNotSetChannel(param)
	return
}
