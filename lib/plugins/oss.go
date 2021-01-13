/**
* @Author:changjiang
* @Description:
* @File:oss
* @Version: 1.0.0
* @Date 2020/4/7 6:48 下午
 */
package plugins

import (
	"io/ioutil"
	systemLog "log"
	"sync"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

type Oss struct {
	NameSpace       string `json:"namespace" yaml:"-"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `json:"accesskeyid" yaml:"accesskeyid"`
	AccessKeySecret string `json:"accesskeysecret" yaml:"accesskeysecret"`
	BucketName      string `json:"bucketname" yaml:"bucketname"`
	RoleArn         string `json:"rolearn" yaml:"rolearn"`
	SessionName     string `json:"sessionname" yaml:"sessionname"`
	BucketUrl       string `json:"bucketurl" yaml:"bucketurl"`
	DirName         string `json:"dirname" yaml:"dirname"`
	ExpiredTime     uint   `json:"expiredtime" yaml:"expiredtime"`
}

var oss = make(map[string]Oss)

func PluginOss() (err error) {

	systemLog.Printf("【INFO】Load oss config.")
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)
	io.SystemOutPrintln("Load common config start")
	defer func() {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("oss config is: '%v' \n", config.ToString())
		io.SystemOutPrintf("load oss config finished \n")
	}()

	var yamlFile []byte
	var v map[string]Oss
	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("oss.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &v); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load oss config config err(%#v) \n", err)
	}
	for k, d := range v {
		if k == "" {
			continue
		}
		d.NameSpace = k
		systemLog.Printf("【INFO】oss config:%v", d)
		oss[d.NameSpace] = d
	}
	return
}

func GetOssConfig(nameSpace ...string) *Oss {
	var space = "default"
	leng := len(nameSpace)
	if leng > 1 {
		panic("GetOssConfig 最多支持1个参数")
	} else if leng == 1 {
		space = nameSpace[0]
	}
	o := oss[space]
	return &o
}
