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
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"io/ioutil"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

type CommonConfig struct {
	Domain     string `json:"domain" yaml:"domain"`
	AppDomain  string `json:"app_domain" yaml:"app_domain"`
	DomainUser string `json:"domain_user" yaml:"domain_user"`
	DomainApi  string `json:"domain_api" yaml:"domain_api"`
}

func (r *CommonConfig) ToString() string {
	v, _ := json.Marshal(r)
	return string(v)
}

var config CommonConfig

func GetCommonConfig() *CommonConfig {
	return &config
}
func PluginLoadCommonParams(arg  *app_start.PluginsOperate) (err error) {
	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintln("Load common config start")
	defer func() {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("common config is: '%v' \n", config.ToString())
		io.SystemOutPrintf("load common config finished \n")
	}()

	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("common.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load common config config err(%#v) \n", err)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("Load common config is : '%#v' ", config)
	return
}
