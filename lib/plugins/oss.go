// Package plugins
// /**
package plugins

import (
	"encoding/json"
	systemLog "log"
	"os"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

type Oss struct {
	NameSpace           string                  `json:"namespace" yaml:"-"`
	Endpoint            string                  `json:"endpoint" yaml:"endpoint"`        //Cdn地址
	EndSrcPoint         string                  `json:"endsrcpoint"  yaml:"endsrcpoint"` //源站地址
	AccessKeyId         string                  `json:"accesskeyid" yaml:"accesskeyid"`
	AccessKeySecret     string                  `json:"accesskeysecret" yaml:"accesskeysecret"`
	BucketName          string                  `json:"bucketname" yaml:"bucketname"`
	RoleArn             string                  `json:"rolearn" yaml:"rolearn"`
	SessionName         string                  `json:"sessionname" yaml:"sessionname"`
	BucketUrl           string                  `json:"bucketurl" yaml:"bucketurl"`
	DirName             string                  `json:"dirname" yaml:"dirname"`
	ExpiredTime         uint                    `json:"expiredtime" yaml:"expiredtime"`
	PipelineIdVideo     string                  `json:"pipelineidvideo" yaml:"pipelineidvideo"`         //70b91995b94846b7a2f8c4f7df362d04
	VideoBucketLocation string                  `json:"videobucketlocation" yaml:"videobucketlocation"` //例:oss-cn-beijing
	ParseCodeTemp       []ParseCodeTemplateItem `json:"parsecodetemp" yaml:"parsecodetemp"`             //转码模板配置
	CDNExpiredTime      int64                   `json:"cdn_expired_time" yaml:"cdnexpiredtime"`
	CDNAuthKey          string                  `json:"cdn_auth_key" yaml:"cdnauthkey"`
}
type ParseCodeTemplateItem struct {
	TemplateId string `json:"templateid"  yaml:"templateid"`
	ExtName    string `json:"extname" yaml:"extname"`
	Label      string `json:"label" yaml:"label"`
	Key        string `json:"key" yaml:"key"`
}

var oss = make(map[string]Oss)

func PluginOss(arg *app_start.PluginsOperate) (err error) {

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
	if yamlFile, err = os.ReadFile(common.GetCommonConfigFilePath("oss.yaml", true)); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   %#v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &v); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load oss config config err(%#v) \n", err)
	}
	for k, d := range v {
		if k == "" {
			continue
		}
		d.NameSpace = k
		systemLog.Printf("【INFO】oss config(%s)", k)
		printConfig(&d)
		oss[d.NameSpace] = d
	}
	return
}
func printConfig(d *Oss) (err error) {
	var bt []byte
	if bt, err = json.Marshal(d); err != nil {
		panic(err)
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bt, &data); err != nil {
		panic(err)
		return
	}
	for key, value := range data {
		systemLog.Printf("【%s】%+v", key, value)
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
