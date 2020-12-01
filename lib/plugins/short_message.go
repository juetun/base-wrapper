package plugins

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/plugins/short_message_impl"
	"github.com/spf13/viper"
)

// 短信插件初始化
func PluginShortMessage() (err error) {

	io.SystemOutPrintln("Load short message start")
	configSource := viper.New()
	configSource.SetConfigName("message") // name of config file (without extension)
	configSource.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	dir := common.GetConfigFileDirectory()

	configSource.AddConfigPath(dir)   // path to look for the config file in
	err = configSource.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error  short message file: %v \n", err))
		return
	}
	// 数据库配置信息存储对象
	var configs = make(map[string]ShortMessageConfig)

	if err = configSource.Unmarshal(&configs); err != nil {
		io.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("Load  short message config failure  '%v' ", configs)
		panic(err)
	}
	// 初始化短信通道
	shortHandle := map[string]app_obj.ShortMessageInter{}
	for nameSpace, value := range configs {
		shortHandle[nameSpace] = initShortMessage(nameSpace, &value)
	}

	app_obj.ShortMessageObj = app_obj.NewShortMessage(shortHandle)

	// 监听配置变动
	viper.WatchConfig()
	viper.OnConfigChange(databaseFileChange)
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
