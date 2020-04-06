package plugins

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	"github.com/juetun/app-web/lib/app_obj"
	"github.com/juetun/app-web/lib/common"
	"github.com/spf13/viper"
)

type Mysql struct {
	NameSpace    string `json:"name_space"`
	Addr         string `json:"addr" yaml:"addr"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
}

func PluginMysql() (err error) {
	loadMysqlConfig()
	return
}

var io = common.NewSystemOut().SetInfoType(common.LogLevelInfo)

func loadMysqlConfig() (err error) {

	io.SystemOutPrintln("Load database start")
	configSource := viper.New()
	configSource.SetConfigName("database") // name of config file (without extension)
	configSource.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	dir := common.GetConfigFileDirectory()

	configSource.AddConfigPath(dir)   // path to look for the config file in
	err = configSource.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		io.SetInfoType(common.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error database file: %v \n", err))
		return
	}
	// 数据库配置信息存储对象
	var mysqlConfig = make(map[string]Mysql)

	if err = configSource.Unmarshal(&mysqlConfig); err != nil {
		io.SetInfoType(common.LogLevelInfo).
			SystemOutPrintf("Load database config failure  '%v' ", mysqlConfig)
		panic(err)
	}
	for key, value := range mysqlConfig {
		initMysql(key, &value)
	}

	viper.WatchConfig()
	viper.OnConfigChange(databaseFileChange)
	io.SetInfoType(common.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("Database load config finished \n"))
	return
}

// 数据库配置文件改变了加载动作
func databaseFileChange(e fsnotify.Event) { // 热加载

	fmt.Println("Database config file changed:", e.Name)
}
func initMysql(nameSpace string, config *Mysql) {
	db := getMysql(nameSpace, config.Addr)

	// 开启 Logger, 以展示详细的日志
	db.LogMode(true)

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(config.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Hour)
	if err := createTable(db); err != nil {
		panic(err)
	}

}
func getMysql(nameSpace, addr string) *gorm.DB {
	io.SetInfoType(common.LogLevelInfo).
		SystemOutPrintf("init mysql %s", addr)
	db, err := gorm.Open("mysql", addr)
	if err != nil {
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
