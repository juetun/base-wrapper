/**
* @Author:changjiang
* @Description:
* @File:oss
* @Version: 1.0.0
* @Date 2020/4/7 6:48 下午
 */
package plugins

import (
	systemLog "log"
	"sync"

	"github.com/juetun/base-wrapper/lib/common"
	"github.com/spf13/viper"
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

	configSource := viper.New()
	configSource.AddConfigPath(common.GetConfigFileDirectory())
	configSource.SetConfigName("oss")
	configSource.SetConfigType("yml")
	if err := configSource.ReadInConfig(); err != nil {
		panic(err)
	}

	var v map[string]Oss
	// 直接反序列化为Struct
	if err := configSource.Unmarshal(&v); err != nil {
		panic(err)
	}
	for k, d := range v {
		if k == "" {
			continue
		}
		d.NameSpace = k
		systemLog.Printf("【INFO】oss config:%v", d)
		oss[d.NameSpace] = d
	}
	systemLog.Printf("【INFO】Load oss config finished")
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
