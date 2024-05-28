package plugins

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	MysqlAppConfig struct {
		Connects            []string `json:"connects" yaml:"connects"`                         //当前应用使用了的数据库连接
		Default             string   `json:"default"  yaml:"default"`                          //默认数据库
		DistributedConnects []string `json:"distributed_connects" yaml:"distributed_connects"` //需要使用的分布式数据库连接
	}
	Mysql struct {
		NameSpace    string `json:"name_space" yaml:"name_space"`
		Addr         string `json:"addr" yaml:"addr"`
		MaxIdleConns int    `json:"max_idle_conns" yaml:"maxidleconns"`
		MaxOpenConns int    `json:"max_open_conns" yaml:"maxopenconns"`
	}
)

func (r *Mysql) ToString() (res string) {
	res = fmt.Sprintf("name_space:%s ,addr:%s ,max_idle_conns:%d ,max_open_conns:%d", r.NameSpace, r.getHiddenPwd(), r.MaxIdleConns, r.MaxOpenConns)
	return
}

func (r *Mysql) getHiddenPwd() (res string) {
	addr := strings.Split(r.Addr, ":")
	if len(addr) >= 2 {
		sArr := make([]int32, 0, len(addr[1]))
		for k, it := range addr[1] {
			if k < 4 {
				sArr = append(sArr, '*')
			} else {
				sArr = append(sArr, it)
			}
		}
		addr[1] = string(sArr)
	}
	res = strings.Join(addr, ":")
	return
}

func PluginMysql(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	loadMysqlConfig()
	return
}

var io = base.NewSystemOut().SetInfoType(base.LogLevelInfo)

func loadMysqlConfig() (err error) {

	io.SystemOutPrintln("Load database start")
	var (
		yamlFile              []byte
		mysqlConfig           MysqlAppConfig
		itemMysqlConfig       Mysql
		mapMysqlConfig        map[string]Mysql
		filePath, connectName string
		ok                    bool
	)
	filePath = common.GetConfigFilePath("database.yaml")
	if yamlFile, err = os.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &mysqlConfig); err != nil {
		io.SystemOutFatalf("load database config err(%+v) \n", err)
	}

	//读取common_config配置文件中的信息
	filePath = common.GetCommonConfigFilePath("database.yaml")
	if yamlFile, err = os.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}
	if err = yaml.Unmarshal(yamlFile, &mysqlConfig); err != nil {
		io.SystemOutFatalf("load database config err(%+v) \n", err)
	}

	for _, connectName = range mysqlConfig.Connects {
		if connectName == "" {
			continue
		}
		if itemMysqlConfig, ok = mapMysqlConfig[connectName]; !ok {
			err = fmt.Errorf("当前common_config中不支持您要使用的数据库连接(%v)", connectName)
			io.SystemOutFatalf("load database config err(%+v) \n", err)
			return
		}
		io.SystemOutPrintf("【 %s 】%+v \n", connectName, itemMysqlConfig.ToString())
		initMysql(connectName, &itemMysqlConfig, mysqlConfig.Default)
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

func initMysql(nameSpace string, config *Mysql, defaultNameString string) {
	db := getMysql(nameSpace, defaultNameString, config)
	var (
		err   error
		sqlDB *sql.DB
	)
	if sqlDB, err = db.DB(); err != nil {
		panic(err)
	}

	// 开启 Logger, 以展示详细的日志
	// db.LogMode(true)
	//
	// // mysql 日志处理
	// db.SetLogger(common.NewGOrmLog(db))
	//
	// // 全局禁用表名复数
	// db.SingularTable(true)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func getMysql(nameSpace, defaultNameString string, addr *Mysql) *gorm.DB {
	var db *gorm.DB
	var err error
	// 数据库连接不可用会自动报错
	if db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       addr.Addr, // data source name
		DefaultStringSize:         256,       // default size for string fields
		DisableDatetimePrecision:  true,      // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,      // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,      // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,     // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: common.NewWithLogger(app_obj.GetLog()), // common.NewGOrmLog(),
	}); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutPrintf(fmt.Sprintf("Fatal error database file: %v \n", err))
		panic(err)
	}
	app_obj.DbMysql[nameSpace] = db
	if nameSpace == defaultNameString && defaultNameString != "" {
		app_obj.DbMysql[app_obj.DefaultDbNameSpace] = db
	}
	return app_obj.DbMysql[nameSpace]
}
func createTable(dbc *gorm.DB) error {
	var models = make([]interface{}, 0)
	for _, m := range models {
		dbt := dbc.Migrator()
		if dbt.HasTable(m) {
			continue
		}

		if err := dbt.CreateTable(m); err != nil {
			return err
		}
	}
	return nil
}
