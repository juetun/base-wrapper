package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/go-redis/redis"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

func PluginRedis() (err error) {
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()
	loadRedisConfig()
	return
}

func loadRedisConfig() (err error) {

	io.SystemOutPrintln("Load redis start")

	// 数据库配置信息存储对象
	var config = make(map[string]Redis)
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(common.GetConfigFilePath("redis.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Load redis config config err(%#v) \n", err)
	}
	for key, value := range config {
		initRedis(key, &value)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("load redis config finished \n"))
	return
}

func initRedis(nameSpace string, configs *Redis) {
	var err error
	var conf = redis.Options{
		Addr:         configs.Addr,
		DB:           configs.DB,
		PoolSize:     configs.PoolSize,
		MinIdleConns: configs.MinIdleConns,
		Password:     configs.Password,
	}

	io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Init redis is  '%s'", RedisOptionToString(&conf))
	// 初始化Redis连接信息
	app_obj.DbRedis[nameSpace] = redis.NewClient(&conf)

	_, err = app_obj.DbRedis[nameSpace].Ping().Result()

	if err != nil {
		io.SetInfoType(base.LogLevelError).SystemOutPrintf(fmt.Sprintf("Load  redis config is error \n"))
		panic(err)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("Load  redis config is finished \n"))

}

type Redis struct {
	NameSpace    string `json:"name_space"`
	Addr         string `json:"addr" yaml:"addr"`
	DB           int    `json:"db" yaml:"db"`
	Password     string `json:"password" yaml:"password"`
	PoolSize     int    `json:"pool_size" yaml:"poolsize"`
	MinIdleConns int    `json:"min_idle_conns" yaml:"minidleconns"`
}

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
