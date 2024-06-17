package plugins

import (
	"fmt"
	"os"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins/short_message_impl"
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

// PluginShortMessage 短信插件初始化
func PluginShortMessage(arg *app_start.PluginsOperate) (err error) {

	io.SystemOutPrintln("Load short message start")
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	// 数据库配置信息存储对象
	var (
		conShortMsgConfigs    = make(map[string]app_obj.ShortMessageConfig)
		itemConfig            app_obj.ShortMessageConfig
		configs               app_obj.ShortMessageAppConfig
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
	shortHandle := map[string]app_obj.ShortMessageInter{}

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
		shortHandle[connectName] = initShortMessage(connectName, &itemConfig)
	}

	app_obj.ShortMessageObj = app_obj.NewShortMessage(shortHandle)
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("ShortMessage load config finished \n"))
	return

}

func initShortMessage(nameSpace string, shortMessageConfig *app_obj.ShortMessageConfig) (res app_obj.ShortMessageInter) {
	switch nameSpace { // 短信通道配置 结构体映射
	case ShortMessageSms100Sms:
		res = short_message_impl.NewSms100(shortMessageConfig)
		res.InitClient()
	case ShortMessageSmsFeiGe:
		res = short_message_impl.NewFeiGe(shortMessageConfig)
		res.InitClient()
	case ShortMessageSmsAliYun:
		res = short_message_impl.NewAliYunSms(shortMessageConfig)
		res.InitClient()
	default:
		var err = fmt.Errorf("当前不支持此短信通道(%s)", nameSpace)
		panic(err)
	}
	return
}
