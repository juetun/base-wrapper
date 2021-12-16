package plugins

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins/short_message_impl"
	"gopkg.in/yaml.v2"
)

// PluginShortMessage 短信插件初始化
func PluginShortMessage(arg *app_start.PluginsOperate) (err error) {

	io.SystemOutPrintln("Load short message start")
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	// 数据库配置信息存储对象
	var configs = make(map[string]ShortMessageConfig)
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("message.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &configs); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load short message err(%#v) \n", err)
	}

	// 初始化短信通道
	shortHandle := map[string]app_obj.ShortMessageInter{}
	for nameSpace, value := range configs {
		shortHandle[nameSpace] = initShortMessage(nameSpace, &value)
	}

	app_obj.ShortMessageObj = app_obj.NewShortMessage(shortHandle)
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("ShortMessage load config finished \n"))
	return

}
func initShortMessage(nameSpace string, shortMessageConfig *ShortMessageConfig) (res app_obj.ShortMessageInter) {
	switch nameSpace { // 短信通道配置 结构体映射
	case "100sms":
		res = short_message_impl.NewSms100()
	case "feige":
		res = short_message_impl.NewFeiGe()
	default:
		panic(fmt.Sprintf("当前不支持此短信通道(%s)", nameSpace))
	}
	return
}

type ShortMessageConfig struct {
	Url       string `json:"url" yml:"Url"`
	AppKey    string `json:"app_key" yml:"AppKey"`
	AppSecret string `json:"app_secret" yml:"AppSecret"`
}
