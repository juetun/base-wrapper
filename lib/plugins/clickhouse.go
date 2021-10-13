package plugins

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginClickHouse(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	_ = loadClickHouse()
	return
}
func loadClickHouse() (err error) {

	io.SystemOutPrintln("Load clickHouse start")
	var configClickHouse map[string]ClickHouse
	var yamlFile []byte
	filePath := common.GetConfigFilePath("clickhouse.yaml")
	if yamlFile, err = ioutil.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &configClickHouse)
		err != nil {
		io.SystemOutFatalf("load clickHouse config err(%+v) \n", err)
	}
	io.SystemOutPrintf("load clickHouse config is:%+v \n", configClickHouse)
	for key, value := range configClickHouse {
		_ = initClickHouse(key, &value)
	}

	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("Database ClickHouse load config finished \n"))
	return
}

// 初始化ClickHouse链接句柄
// 操作文档 https://github.com/ClickHouse/clickhouse-go
// 官方文档 https://clickhouse.com/docs/en/interfaces/third-party/client-libraries/
func initClickHouse(key string, clickHouseConfig *ClickHouse) (err error) {

	var connect *sql.DB

	if connect, err = sql.Open("clickhouse", clickHouseConfig.Addr); err != nil {
		io.SystemOutFatalf("load clickHouse config err(%+v) \n", err)
	}
	if err = connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			io.SystemOutFatalf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			return
		}
		io.SystemOutFatalf("load clickHouse config err(%+v) \n", err)
		return
	}
	app_obj.DbClickHouse[key] = &app_obj.ClickHouseClient{
		DB: connect,
	}

	return
}

type ClickHouse struct {
	NameSpace string `json:"name_space" yaml:"name_space"`
	Addr      string `json:"addr" yaml:"addr"` // clickhouse连接地址 "tcp://127.0.0.1:9000?debug=true"
}
