package plugins

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

type Mysql struct {
	NameSpace    string `json:"name_space" yaml:"name_space"`
	Addr         string `json:"addr" yaml:"addr"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"maxidleconns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"maxopenconns"`
}

func PluginMysql() (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	loadMysqlConfig()
	return
}

var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

func loadMysqlConfig() (err error) {

	io.SystemOutPrintln("Load database start")
	var mysqlConfig map[string]Mysql
	var yamlFile []byte
	filePath := common.GetConfigFilePath("database.yaml")
	if yamlFile, err = ioutil.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &mysqlConfig); err != nil {
		io.SystemOutFatalf("load database config err(%+v) \n", err)
	}
	io.SystemOutPrintf("load database config is:%+v \n", mysqlConfig)
	for key, value := range mysqlConfig {
		initMysql(key, &value)
	}

	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("Database load config finished \n"))
	return
}

// 数据库配置文件改变了加载动作
func databaseFileChange(e fsnotify.Event) { // 热加载
	fmt.Println("Database config file changed:", e.Name)
	// 重新加载数据库配置
	loadMysqlConfig()
}
func initMysql(nameSpace string, config *Mysql) {
	db := getMysql(nameSpace, config)

	// 开启 Logger, 以展示详细的日志
	db.LogMode(true)

	// mysql 日志处理
	db.SetLogger(common.NewGOrmLog(db))

	// 全局禁用表名复数
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(config.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Hour)
}
func getMysql(nameSpace string, addr *Mysql) *gorm.DB {
	io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("init mysql :%#v", addr)
	var db *gorm.DB
	var err error

	//数据库连接不可用会自动报错
	if db, err = gorm.Open("mysql", addr.Addr); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutPrintf(fmt.Sprintf("Fatal error database file: %v \n", err))
		panic(err)
	}
	app_obj.DbMysql[nameSpace] = db
	return app_obj.DbMysql[nameSpace]
}
func createTable(dbc *gorm.DB) error {
	var models = make([]interface{}, 0)
	for _, m := range models {
		if dbc.HasTable(m) {
			continue
		}
		if err := dbc.CreateTable(m).Error; err != nil {
			return err
		}
	}
	return nil
}
