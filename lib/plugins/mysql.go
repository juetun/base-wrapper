package plugins

import (
	"fmt"
	"io/ioutil"
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

type Mysql struct {
	NameSpace    string `json:"name_space" yaml:"name_space"`
	Addr         string `json:"addr" yaml:"addr"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"maxidleconns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"maxopenconns"`
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
	sqlDB, err := db.DB()
	if err != nil {
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

// func getLogger() logger.Interface {
//   	return logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
// 		logger.Config{
// 			SlowThreshold:             time.Second,   // Slow SQL threshold
// 			LogLevel:                  logger.Silent, // Log level
// 			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
// 			Colorful:                  true,         // Disable color
// 		},
// 	)
// }
func getMysql(nameSpace string, addr *Mysql) *gorm.DB {
	io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("init mysql :%#v", addr)
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
