/**
* @Author:changjiang
* @Description:
* @File:appmap
* @Version: 1.0.0
* @Date 2020/10/21 11:00 下午
 */
package plugins

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/spf13/viper"
)

func PluginAppMap() (err error) {
	io.SystemOutPrintln("Load AppMap start")
	defer io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("AppMap load config finished \n"))

	configSource := viper.New()
	configSource.SetConfigName("appmap") // name of config file (without extension)
	configSource.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	dir := common.GetConfigFileDirectory()

	configSource.AddConfigPath(dir)   // path to look for the config file in
	err = configSource.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error  AppMap file: %v \n", err))
		return
	}
	// 数据库配置信息存储对象
	app_obj.AppMap = make(map[string]string)

	if err = configSource.Unmarshal(&app_obj.AppMap); err != nil {
		io.SetInfoType(base.LogLevelInfo).
			SystemOutPrintf("Load  AppMap config failure  '%v' ", app_obj.AppMap)
		panic(err)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("app map config :%+v", app_obj.AppMap)
	// 监听配置变动
	viper.WatchConfig()
	viper.OnConfigChange(databaseFileChange)
	return

}
