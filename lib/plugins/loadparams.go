/**
* @Author:changjiang
* @Description:
* @File:loadparams
* @Version: 1.0.0
* @Date 2020/5/6 8:51 下午
 */
package plugins

import (
	"encoding/json"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/spf13/viper"
)

type CommonConfig struct {
	Domain     string `json:"domain" `
	AppDomain  string `json:"app_domain" `
	DomainUser string `json:"domain_user"`
	DomainApi  string `json:"domain_api"`
}

func (r *CommonConfig) ToString() string {
	v, _ := json.Marshal(r)
	return string(v)
}

var config CommonConfig

func GetCommonConfig() *CommonConfig {
	return &config
}
func PluginLoadCommonParams() (err error) {
	app := common.NewApplication()
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintln("Load common config start")
	defer func() {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("common config is: '%v' \n", config.ToString())
		io.SystemOutPrintf("load common config finished \n")
	}()
	viper.SetConfigName("common") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	dir := common.GetConfigFileDirectory()
	io.SystemOutPrintf("config directory is : '%s' ", dir)
	viper.AddConfigPath(dir + "/../" + app.AppEnv + "/") // path to look for the config file in
	err = viper.ReadInConfig()                           // Find and read the config file
	if err != nil {                                      // Handle errors reading the config file
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) { // 热加载
		fmt.Println("Config file changed:", e.Name)
	})

	config.Domain = viper.GetString("common.domain")
	config.AppDomain = viper.GetString("common.app_domain")
	config.DomainUser = viper.GetString("common.domain_user")
	config.DomainApi = viper.GetString("common.domain_api")
	return
}
