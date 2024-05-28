package plugins

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/juetun/base-wrapper/lib/app/app_start"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"gopkg.in/yaml.v2"
)

// PluginElasticSearchV7 ElasticSearch检索初始化入口
func PluginElasticSearchV7(arg *app_start.PluginsOperate) (err error) {
	_ = arg
	var syncLock sync.Mutex
	syncLock.Lock()
	defer syncLock.Unlock()

	io.SystemOutPrintln("Load ElasticSearch start")
	defer io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf(fmt.Sprintf("ElasticSearch load config finished \n"))

	var yamlFile []byte
	if yamlFile, err = os.ReadFile(common.GetConfigFilePath("elasticsearch.yaml")); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("yamlFile.Get err   #%v \n", err)
	}
	// 数据库配置信息存储对象
	var configs = make(map[string]Config)
	if err = yaml.Unmarshal(yamlFile, &configs); err != nil {
		io.SetInfoType(base.LogLevelFatal).SystemOutFatalf("Fatal error elastic_search file(%#v) \n", err)
	}

	for key, value := range configs {
		esConfig := orgConfig(&value)
		initEs(key, esConfig)
	}
	io.SetInfoType(base.LogLevelInfo).SystemOutPrintf("Load elastic_search config finished ")
	return
}

func orgConfig(config *Config) (configOption []SetElasticSearchConfigOption) {
	if len(config.Addresses) > 0 {
		configOption = append(configOption, SetAddresses(config.Addresses))
	}
	if config.Username != "" {
		configOption = append(configOption, SetUsername(config.Username))
	}
	if config.Password != "" {
		configOption = append(configOption, SetPassword(config.Password))
	}
	configOption = append(configOption, SetDiscoverNodesOnStart(config.DiscoverNodesOnStart))
	configOption = append(configOption, SetTransport(&http.Transport{
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second,
		DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS11,
		},
	}))
	return
}

func initEs(nameSpace string, configOption []SetElasticSearchConfigOption) {
	esConfig := NewElasticSearchConfig(configOption...)
	var err error
	var handler *elasticsearch.Client

	var configInterface map[string]interface{}
	var (
		bt []byte
	)
	io.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("Load elastic_search config(%s) ", nameSpace)

	var showEsConfig ShowEsConfig
	showEsConfig.ParseFromEsConfig(&esConfig.Config)
	if bt, err = json.Marshal(showEsConfig); err != nil {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("init es failure \n"))
		panic(err)
	}
	if err = json.Unmarshal(bt, &configInterface); err != nil {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("init es failure \n"))
		panic(err)
	}
	for key, data := range configInterface {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("【%s】%+v \n", key, data))
	}
	handler, err = elasticsearch.NewClient(esConfig.Config)
	if err != nil {
		io.SetInfoType(base.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("init es failure \n"))
		panic(err)
	}
	app_obj.ElasticSearchV7Maps[nameSpace] = handler

}

type ShowEsConfig struct {
	Addresses             []string      `json:"addresses"`
	Username              string        `json:"username"`
	Password              string        `json:"password"`
	CloudID               string        `json:"cloud_id"`
	APIKey                string        `json:"api_key"`
	ServiceToken          string        `json:"service_token"`
	Header                http.Header   `json:"header"`
	CACert                string        `json:"ca_cert"`
	RetryOnStatus         []int         `json:"retry_on_status"`
	DisableRetry          bool          `json:"disable_retry"`
	EnableRetryOnTimeout  bool          `json:"enable_retry_on_timeout"`
	MaxRetries            int           `json:"max_retries"`
	DiscoverNodesOnStart  bool          `json:"discover_nodes_on_start"`
	DiscoverNodesInterval time.Duration `json:"discover_nodes_interval"`
	EnableMetrics         bool          `json:"enable_metrics"`
	EnableDebugLogger     bool          `json:"enable_debug_logger"`
	DisableMetaHeader     bool          `json:"disable_meta_header"`
}

func (r *ShowEsConfig) ParseFromEsConfig(data *elasticsearch.Config) {
	r.Addresses = data.Addresses
	r.Username = data.Username
	r.Password = data.Password
	r.CloudID = data.CloudID
	r.APIKey = data.APIKey
	r.ServiceToken = data.ServiceToken
	r.Header = data.Header
	r.CACert = string(data.CACert)
	r.RetryOnStatus = data.RetryOnStatus
	r.DisableRetry = data.DisableRetry
	r.EnableRetryOnTimeout = data.EnableRetryOnTimeout
	r.MaxRetries = data.MaxRetries
	r.DiscoverNodesOnStart = data.DiscoverNodesOnStart
	r.DiscoverNodesInterval = data.DiscoverNodesInterval
	r.EnableMetrics = data.EnableMetrics
	r.EnableDebugLogger = data.EnableDebugLogger
	r.DisableMetaHeader = data.DisableMetaHeader
}

type ElasticSearchConfig struct {
	elasticsearch.Config
}

func NewElasticSearchConfig(arg ...SetElasticSearchConfigOption) (elasticSearchConfig *ElasticSearchConfig) {
	elasticSearchConfig = &ElasticSearchConfig{
		Config: elasticsearch.Config{
			Addresses: []string{}, // 默认访问位置9200
		},
	}
	for _, handler := range arg {
		handler(elasticSearchConfig)
	}
	return
}

type SetElasticSearchConfigOption func(config *ElasticSearchConfig)

func SetAddresses(Addresses []string) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Addresses = Addresses
	}
}
func SetUsername(Username string) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Username = Username
	}
}
func SetPassword(Password string) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Password = Password
	}
}

func SetCloudID(CloudID string) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.CloudID = CloudID
	}
}
func SetAPIKey(APIKey string) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.APIKey = APIKey
	}
}

func SetHeader(Header http.Header) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Header = Header
	}
}
func SetCACert(CACert []byte) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.CACert = CACert
	}
}
func SetRetryOnStatus(RetryOnStatus []int) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.RetryOnStatus = RetryOnStatus
	}
}
func SetDisableRetry(DisableRetry bool) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.DisableRetry = DisableRetry
	}
}
func SetEnableRetryOnTimeout(EnableRetryOnTimeout bool) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.EnableRetryOnTimeout = EnableRetryOnTimeout
	}
}
func SetMaxRetries(MaxRetries int) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.MaxRetries = MaxRetries
	}
}
func SetDiscoverNodesOnStart(DiscoverNodesOnStart bool) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.DiscoverNodesOnStart = DiscoverNodesOnStart
	}
}
func SetDiscoverNodesInterval(DiscoverNodesInterval time.Duration) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.DiscoverNodesInterval = DiscoverNodesInterval
	}
}
func SetEnableMetrics(EnableMetrics bool) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.EnableMetrics = EnableMetrics
	}
}
func SetEnableDebugLogger(EnableDebugLogger bool) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.EnableDebugLogger = EnableDebugLogger
	}
}
func SetRetryBackoff(RetryBackoff func(attempt int) time.Duration) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.RetryBackoff = RetryBackoff
	}
}
func SetTransport(Transport http.RoundTripper) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Transport = Transport
	}
}
func SetLogger(Logger estransport.Logger) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Logger = Logger
	}
}
func SetSelector(Selector estransport.Selector) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.Selector = Selector
	}
}
func SetConnectionPoolFunc(ConnectionPoolFunc func([]*estransport.Connection, estransport.Selector) estransport.ConnectionPool) SetElasticSearchConfigOption {
	return func(config *ElasticSearchConfig) {
		config.ConnectionPoolFunc = ConnectionPoolFunc
	}
}

type Config struct {
	Addresses []string `json:"addresses" yaml:"addresses"` // A list of Elasticsearch nodes to use.
	Username  string   `json:"username" yaml:"username"`   // Username for HTTP Basic Authentication.
	Password  string   `json:"password" yaml:"password"`   // Password for HTTP Basic Authentication.

	CloudID string `json:"cloud_id" yaml:"cloudid"` // Endpoint for the Elastic Service (https://elastic.co/cloud).
	APIKey  string `json:"api_key" yaml:"apikey"`   // Base64-encoded token for authorization; if set, overrides username and password.

	Header http.Header // Global HTTP request header.

	// PEM-encoded certificate authorities.
	// When set, an empty certificate pool will be created, and the certificates will be appended to it.
	// The option is only valid when the transport is not specified, or when it's http.Transport.
	CACert []byte `json:"ca_cert" yaml:"ca_cert"`

	RetryOnStatus        []int `json:"retry_on_status" yaml:"retryonstatus"`                // List of status codes for retry. Default: 502, 503, 504.
	DisableRetry         bool  `json:"disable_retry" yaml:"disableretry"`                   // Default: false.
	EnableRetryOnTimeout bool  `json:"enable_retry_on_timeout" yaml:"enableretryontimeout"` // Default: false.
	MaxRetries           int   `json:"max_retries" yaml:"maxretries"`                       // Default: 3.

	DiscoverNodesOnStart  bool          `json:"discover_nodes_on_start" yaml:"discovernodesonstart"`  // Discover nodes when initializing the client. Default: false.
	DiscoverNodesInterval time.Duration `json:"discover_nodes_interval" yaml:"discovernodesinterval"` // Discover nodes periodically. Default: disabled.

	EnableMetrics     bool `json:"enable_metrics" yaml:"enable_metrics"`         // Enable the metrics collection.
	EnableDebugLogger bool `json:"enable_debug_logger" yaml:"enabledebuglogger"` // Enable the debug logging.

	RetryBackoff func(attempt int) time.Duration // Optional backoff duration. Default: nil.

	Transport http.RoundTripper    // The HTTP transport object.
	Logger    estransport.Logger   // The logger object.
	Selector  estransport.Selector // The selector object.

	// Optional constructor function for a custom ConnectionPool. Default: nil.
	ConnectionPoolFunc func([]*estransport.Connection, estransport.Selector) estransport.ConnectionPool
}
