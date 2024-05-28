package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginRedis(arg *app_start.PluginsOperate) (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	loadRedisConfig()
	return
}

func loadRedisConfig() (err error) {

	io.SystemOutPrintln("Load redis start")
	var (
		yamlFile              []byte
		redisAppConfig        RedisAppConfig
		itemConfig            *Redis
		mapRedisConfig        map[string]*Redis
		filePath, connectName string
		ok                    bool
	)
	// 数据库配置信息存储对象
	if yamlFile, err = os.ReadFile(common.GetConfigFilePath("redis.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &redisAppConfig); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load redis config config err(%#v) \n", err)
	}
	//读取common_config配置文件中的信息
	filePath = common.GetCommonConfigFilePath("database.yaml", true)
	if yamlFile, err = os.ReadFile(filePath); err != nil {
		io.SystemOutFatalf("yamlFile.Get err(%s)  #%v \n", filePath, err)
	}

	if err = yaml.Unmarshal(yamlFile, &mapRedisConfig); err != nil {
		io.SystemOutFatalf("load database config err(%+v) \n", err)
	}

	for _, connectName = range redisAppConfig.Connects {
		if connectName == "" {
			continue
		}
		if itemConfig, ok = mapRedisConfig[connectName]; !ok {
			err = fmt.Errorf("当前common_config中不支持您要使用的数据库连接(%v)", connectName)
			io.SystemOutFatalf("load database config err(%+v) \n", err)
			return
		}
		initRedis(connectName, redisAppConfig.Default, itemConfig)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("load redis config finished \n"))
	return
}

func initRedis(nameSpace, defaultNameSpace string, configs *Redis) {
	var err error
	var conf = redis.Options{
		Addr:         configs.Addr,
		DB:           configs.DB,
		PoolSize:     configs.PoolSize,
		MinIdleConns: configs.MinIdleConns,
		Password:     configs.Password,
	}

	// 初始化Redis连接信息
	app_obj.DbRedis[nameSpace] = redis.NewClient(&conf)
	if defaultNameSpace == nameSpace && defaultNameSpace != "" {
		app_obj.DbRedis[app_obj.DefaultDbNameSpace] = app_obj.DbRedis[nameSpace]
	}
	_, err = app_obj.DbRedis[nameSpace].Ping(context.Background()).Result()

	if err != nil {
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Load  redis config is error \n"))
		io.SetInfoType(base.LogLevelFatal).SystemOutPrintf(fmt.Sprintf("err:%s,conf:%#v", err.Error(), conf))

	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("Load  redis config is finished \n"))

}

type (
	RedisAppConfig struct {
		Connects            []string `json:"connects" yaml:"connects"`                         //当前应用使用了的数据库连接
		Default             string   `json:"default"  yaml:"default"`                          //默认数据库
		DistributedConnects []string `json:"distributed_connects" yaml:"distributed_connects"` //需要使用的分布式数据库连接
	}
	Redis struct {
		NameSpace    string `json:"name_space"`
		Addr         string `json:"addr" yaml:"addr"`
		DB           int    `json:"db" yaml:"db"`
		Password     string `json:"password" yaml:"password"`
		PoolSize     int    `json:"pool_size" yaml:"poolsize"`
		MinIdleConns int    `json:"min_idle_conns" yaml:"minidleconns"`
	}
)

func RedisOptionToString(opt *redis.Options) string {
	type redisOption struct {
		// host:port address.
		Addr string `json:"addr"`

		// Optional password. Must match the password specified in the
		// requirepass server configuration option.
		Password string `json:"password"`
		// Database to be selected after connecting to the server.
		DB int `json:"db"`

		// Maximum number of retries before giving up.
		// Default is to not retry failed commands.
		MaxRetries int `json:"max_retries"`

		// Maximum number of socket connections.
		// Default is 10 connections per every CPU as reported by runtime.NumCPU.
		PoolSize int `json:"pool_size"`
		// Minimum number of idle connections which is useful when establishing
		// new connection is slow.
		MinIdleConns int `json:"min_idle_conns"`

		// Enables read only queries on slave nodes.
		readOnly bool `json:"read_only"`
	}
	var dta = redisOption{
		Addr:         opt.Addr,
		DB:           opt.DB,
		Password:     opt.Password,
		PoolSize:     opt.PoolSize,
		MinIdleConns: opt.MinIdleConns,
	}
	s, _ := json.Marshal(dta)
	return string(s)
}
