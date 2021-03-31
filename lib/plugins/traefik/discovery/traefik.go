// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/plugins/traefik/discovery/middleware_traefik"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TraefikDynamic struct {
	Http HttpTraefik `yaml:"http,omitempty" key_value:"http,omitempty"`
	Tcp  TcpTraefik  `yaml:"tcp,omitempty" key_value:"tcp,omitempty"`
	Udp  UdpTraefik  `yaml:"udp,omitempty" key_value:"udp,omitempty"`
	Tls  TlsTraefik  `yaml:"tls,omitempty" key_value:"tls,omitempty"`
}

type TraefikConfig struct {
	TraefikDynamic
	MapValue []KeyValue //数据结果
}

func NewTraefikConfig() (res *TraefikConfig) {
	res = &TraefikConfig{
	}
	return
}
func NewTraefikConfigTest() (res *TraefikConfig) {
	res = &TraefikConfig{
		TraefikDynamic: TraefikDynamic{
			Http: HttpTraefik{
				Routers: map[string]HttpTraefikRouters{
					"Router0": {
						Rule: "foobar",
						EntryPoints: []string{
							"foobar",
							"foobar",
						},
						Service:     "foobar",
						Middlewares: []string{"foobar", "foobar"},
						Priority:    42,
						Tls: &HttpTls{
							Options:      "foobar",
							CertResolver: "foobar",
							Domains: []HttpDomainTlsItem{
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
							},
						},
					},
					"Router1": {
						Rule: "foobar",
						EntryPoints: []string{
							"foobar",
							"foobar",
						},
						Service:     "foobar",
						Middlewares: []string{"foobar", "foobar"},
						Priority:    42,
						Tls: &HttpTls{
							Options:      "foobar",
							CertResolver: "foobar",
							Domains: []HttpDomainTlsItem{
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
							},
						},
					},
				},
				Services: map[string]HttpTraefikServiceConfig{
					"Service01": {
						LoadBalancer: &HttpLoadBalancer{
							Sticky: &HttpSticky{
								Cookie: &HttpCookie{
									Name:     "foobar",
									Secure:   true,
									SameSite: "foobar",
									HttpOnly: true,
								},
							},
							Servers: []HttpLoadBalancerServer{
								{
									Url: "foobar",
								},
								{
									Url: "foobar",
								},
							},
							HealthCheck: &HttpHealthCheck{
								Scheme:          "foobar",
								Path:            "foobar",
								Port:            42,
								Headers:         map[string]string{"name0": "foobar", "name1": "foobar"},
								Hostname:        "foobar",
								Interval:        10 * time.Second,
								Timeout:         15 * time.Second,
								FollowRedirects: true,
							},
							ResponseForwarding: HttpResponseForwarding{
								FlushInterval: 10 * time.Hour,
							},
							PassHostHeader:   true,
							ServersTransport: "foobar",
						},
					},
					"Service02": {
						Mirroring: &HttpMirroring{
							Service:     "foobar",
							MaxBodySize: 42,
							Mirrors: []HttpMirrors{
								{
									Name:    "foobar",
									Percent: 42,
								},
								{
									Name:    "foobar",
									Percent: 42,
								},
							},
						},
					},
					"Service03": {
						Weighted: &HttpWeighted{
							Services: []HttpWeightedService{
								{
									Name:   "foobar",
									Weight: 42,
								},
								{
									Name:   "foobar",
									Weight: 42,
								},
							},
							Sticky: &HttpSticky{
								Cookie: &HttpCookie{
									Name:     "foobar",
									Secure:   true,
									HttpOnly: true,
									SameSite: "foobar",
								},
							},
						},
					},
				},
				Middlewares: map[string]HttpTraefikMiddleware{
					"Middleware00": middleware_traefik.HttpMiddlewareAddPrefix{
						AddPrefix: middleware_traefik.HttpAddPrefixArg{
							Prefix: "foobar",
						},
					},
					"Middleware01": middleware_traefik.HttpMiddlewareBasicAuth{
						BasicAuth: middleware_traefik.HttpMiddlewareBasicAuthArg{
							Users:        []string{"foobar", "foobar"},
							UsersFile:    "foobar",
							Realm:        "foobar",
							RemoveHeader: true,
							HeaderField:  "foobar",
						},
					},
					"Middleware02": middleware_traefik.HttpMiddlewareBuffering{
						Buffering: middleware_traefik.HttpBufferingArg{
							MaxRequestBodyBytes:  42,
							MemRequestBodyBytes:  42,
							MaxResponseBodyBytes: 42,
							MemResponseBodyBytes: 42,
							RetryExpression:      "foobar",
						},
					},
					"Middleware03": middleware_traefik.HttpMiddlewareChain{
						Chain: middleware_traefik.HttpChainArg{
							Middlewares: []string{"foobar", "foobar"},
						},
					},
					"Middleware04": middleware_traefik.HttpMiddlewareCircuitBreaker{
						CircuitBreaker: middleware_traefik.HttpCircuitBreakerArg{
							Expression: "foobar",
						},
					},
					"Middleware05": middleware_traefik.HttpMiddlewareCompress{
						Compress: middleware_traefik.HttpCompressArg{
							ExcludedContentTypes: []string{"foobar", "foobar"},
						},
					},
					"Middleware06": middleware_traefik.HttpMiddlewareContentType{
						ContentType: middleware_traefik.HttpContentTypeArg{
							AutoDetect: true,
						},
					},
					"Middleware07": middleware_traefik.HttpMiddlewareDigestAuth{
						DigestAuth: middleware_traefik.HttpDigestAuthArg{
							Users:        []string{"foobar", "foobar"},
							UsersFile:    "foobar",
							RemoveHeader: true,
							Realm:        "foobar",
							HeaderField:  "foobar",
						},
					},
					"Middleware08": middleware_traefik.HttpMiddlewareErrors{
						Errors: middleware_traefik.HttpMiddlewareErrorsArg{
							Status:  []string{"foobar", "foobar"},
							Service: "foobar",
							Query:   "foobar",
						},
					},
					"Middleware09": middleware_traefik.HttpMiddlewareForwardAuth{
						ForwardAuth: middleware_traefik.HttpMiddlewareForwardAuthArg{

							Address: "foobar",
							Tls: middleware_traefik.HttpMiddlewareForwardAuthArgTls{
								CA:                 "foobar",
								CaOptional:         true,
								Cert:               "foobar",
								Key:                "foobar",
								InsecureSkipVerify: true,
							},
							TrustForwardHeader:       true,
							AuthResponseHeaders:      []string{"foobar", "foobar"},
							AuthResponseHeadersRegex: "foobar",
							AuthRequestHeaders:       []string{"foobar", "foobar"},
						},
					},
					"Middleware10": middleware_traefik.HttpMiddlewareHeaders{
						Headers: middleware_traefik.HttpMiddlewareHeadersArg{
							CustomRequestHeaders:              map[string]string{"name0": "foobar", "name1": "foobar"},
							CustomResponseHeaders:             map[string]string{"name0": "foobar", "name1": "foobar"},
							AccessControlAllowCredentials:     true,
							AccessControlAllowHeaders:         []string{"foobar", "foobar"},
							AccessControlAllowMethods:         []string{"foobar", "foobar"},
							AccessControlAllowOrigin:          "foobar",
							AccessControlAllowOriginList:      []string{"foobar", "foobar"},
							AccessControlAllowOriginListRegex: []string{"foobar", "foobar"},
							AccessControlExposeHeaders:        []string{"foobar", "foobar"},
							AccessControlMaxAge:               42,
							AddVaryHeader:                     true,
							AllowedHosts:                      []string{"foobar", "foobar"},
							HostsProxyHeaders:                 []string{"foobar", "foobar"},
							SslRedirect:                       true,
							SslTemporaryRedirect:              true,
							SslHost:                           "foobar",
							SslProxyHeaders:                   map[string]string{"name0": "foobar", "name1": "foobar"},
							SslForceHost:                      true,
							StsSeconds:                        42,
							StsIncludeSubdomains:              true,
							StsPreload:                        true,
							ForceSTSHeader:                    true,
							FrameDeny:                         true,
							CustomFrameOptionsValue:           "foobar",
							ContentTypeNosniff:                true,
							BrowserXssFilter:                  true,
							CustomBrowserXSSValue:             "foobar",
							ContentSecurityPolicy:             "foobar",
							PublicKey:                         "foobar",
							ReferrerPolicy:                    "foobar",
							FeaturePolicy:                     "foobar",
							IsDevelopment:                     true,
						},
					},
					"Middleware11": middleware_traefik.HttpMiddlewareIpWhiteList{
						IpWhiteList: middleware_traefik.HttpMiddlewareIpWhiteListArg{
							SourceRange: []string{"foobar", "foobar"},
							IpStrategy: middleware_traefik.HttpMiddlewareIpWhiteListArgIpStrategy{
								Depth:       42,
								ExcludedIPs: []string{"foobar", "foobar"},
							},
						},
					},
					"Middleware12": middleware_traefik.HttpMiddlewareInFlightReq{
						InFlightReq: middleware_traefik.HttpMiddlewareInFlightReqArg{
							Amount: 42,
							SourceCriterion: middleware_traefik.HttpMiddlewareInFlightReqSourceCriterionArg{
								IpStrategy: middleware_traefik.HttpMiddlewareInFlightReqSourceCriterionIpStrategyArg{
									Depth:       42,
									ExcludedIPs: []string{"foobar", "foobar"},
								},
								RequestHeaderName: "foobar",
								RequestHost:       true,
							},
						},
					},
					"Middleware13": middleware_traefik.HttpMiddlewarePassTLSClientCert{
						PassTLSClientCert: middleware_traefik.HttpMiddlewarePassTLSClientCertArg{
							Pem: true,
							Info: middleware_traefik.HttpMiddlewarePassTLSClientCertArgInfo{
								NotAfter:  true,
								NotBefore: true,
								Sans:      true,
								Subject: middleware_traefik.HttpMiddlewarePassTLSClientCertArgSubject{
									Country:         true,
									Province:        true,
									Locality:        true,
									Organization:    true,
									CommonName:      true,
									SerialNumber:    true,
									DomainComponent: true,
								},
								Issuer: middleware_traefik.HttpMiddlewarePassTLSClientCertArgSubject{
									Country:         true,
									Province:        true,
									Locality:        true,
									Organization:    true,
									CommonName:      true,
									SerialNumber:    true,
									DomainComponent: true,
								},
								SerialNumber: true,
							},
						},
					},
					"Middleware14": middleware_traefik.HttpMiddlewarePlugin{
						Plugin: middleware_traefik.HttpMiddlewarePluginArg{
							PluginConf: middleware_traefik.HttpMiddlewarePluginConfArg{
								Foo: "bar",
							},
						},
					},
					"Middleware15": middleware_traefik.HttpMiddlewareRateLimit{
						RateLimit: middleware_traefik.HttpMiddlewareRateLimitArg{
							Average: 42,
							Period:  42,
							Burst:   42,
							SourceCriterion: middleware_traefik.HttpMiddlewareRateLimitSourceCriterionArg{
								IpStrategy: middleware_traefik.HttpMiddlewareRateLimitSourceCriterionIpStrategy{
									Depth:       42,
									ExcludedIPs: []string{"foobar", "foobar"},
								},
								RequestHeaderName: "foobar",
								RequestHost:       true,
							},
						},
					},
					"Middleware16": middleware_traefik.HttpMiddlewareRedirectRegex{
						RedirectRegex: middleware_traefik.HttpMiddlewareRedirectRegexArg{
							Regex:       "foobar",
							Replacement: "foobar",
							Permanent:   true,
						},
					},
					"Middleware17": middleware_traefik.HttpMiddlewareRedirectScheme{
						RedirectScheme: middleware_traefik.HttpMiddlewareRedirectSchemeArg{
							Scheme:    "foobar",
							Port:      "foobar",
							Permanent: true,
						},
					},
					"Middleware18": middleware_traefik.HttpMiddlewareReplacePath{
						ReplacePath: middleware_traefik.HttpMiddlewareReplacePathArg{
							Path: "foobar",
						},
					},
					"Middleware19": middleware_traefik.HttpMiddlewareReplacePathRegex{
						ReplacePathRegex: middleware_traefik.HttpMiddlewareReplacePathRegexArg{
							Regex:       "foobar",
							Replacement: "foobar",
						},
					},
					"Middleware20": middleware_traefik.HttpMiddlewareRetry{
						Retry: middleware_traefik.HttpMiddlewareRetryArg{
							Attempts:        42,
							InitialInterval: 42,
						},
					},
					"Middleware21": middleware_traefik.HttpMiddlewareStripPrefix{
						StripPrefix: middleware_traefik.HttpMiddlewareStripPrefixArg{
							Prefixes:   []string{"foobar", "foobar"},
							ForceSlash: true,
						},
					},
					"Middleware22": middleware_traefik.HttpMiddlewareStripPrefixRegex{
						StripPrefixRegex: middleware_traefik.HttpMiddlewareStripPrefixRegexArg{
							Regex: []string{"foobar", "foobar"},
						},
					},
				},
				ServersTransports: map[string]HttpTraefikServersTransports{
					"ServersTransport0": {
						ServerName:         "foobar",
						InsecureSkipVerify: true,
						RootCAs:            []string{"foobar", "foobar"},
						Certificates: []HttpTraefikServersTransportsCertificates{
							{
								CertFile: "foobar",
								KeyFile:  "foobar",
							},
							{
								CertFile: "foobar",
								KeyFile:  "foobar",
							},
						},
						MaxIdleConnsPerHost: 42,
						ForwardingTimeouts: HttpTraefikServersTransportsForwardingTimeouts{
							DialTimeout:           42 * time.Second,
							ResponseHeaderTimeout: 42 * time.Second,
							IdleConnTimeout:       42 * time.Second,
						},
					},
					"ServersTransport1": {
						ServerName:         "foobar",
						InsecureSkipVerify: true,
						RootCAs:            []string{"foobar", "foobar"},
						Certificates: []HttpTraefikServersTransportsCertificates{
							{
								CertFile: "foobar",
								KeyFile:  "foobar",
							},
							{
								CertFile: "foobar",
								KeyFile:  "foobar",
							},
						},
						MaxIdleConnsPerHost: 42,
						ForwardingTimeouts: HttpTraefikServersTransportsForwardingTimeouts{
							DialTimeout:           42 * time.Second,
							ResponseHeaderTimeout: 42 * time.Second,
							IdleConnTimeout:       42 * time.Second,
						},
					},
				},
			},
			Tcp: TcpTraefik{
				Routers: map[string]TcpTraefikRouters{
					"TCPRouter0": {
						EntryPoints: []string{"foobar", "foobar"},
						Rule:        "foobar",
						Service:     "foobar",
						Tls: &TCPTls{
							//Value:        true,
							CertResolver: "foobar",
							Domains: []TcpDomainTlsItem{
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
							},
							Options:     "foobar",
							PassThrough: true,
						},
					},
					"TCPRouter1": {
						EntryPoints: []string{"foobar", "foobar"},
						Rule:        "foobar",
						Service:     "foobar",
						Tls: &TCPTls{
							//Value:        true,
							CertResolver: "foobar",
							Domains: []TcpDomainTlsItem{
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
								{
									Main: "foobar",
									Sans: []string{"foobar", "foobar"},
								},
							},
							Options:     "foobar",
							PassThrough: true,
						},
					},
				},
				Services: map[string]TcpTraefikServiceConfig{
					"TCPService01": {
						LoadBalancer: &TcpLoadBalancer{
							ProxyProtocol: TcpProxyProtocol{
								Version: 42,
							},
							TerminationDelay: 42,
							Servers: []TcpLoadBalancerServer{
								{
									Address: "foobar",
								},
								{
									Address: "foobar",
								},
							},
						},
					},
					"TCPService02": {
						Weighted: &TcpWeighted{
							Services: []TcpWeightedService{
								{
									Name:   "foobar",
									Weight: 42,
								},
								{
									Name:   "foobar",
									Weight: 42,
								},
							},
						},
					},
				},
			},
			Udp: UdpTraefik{
				Routers: map[string]UdpTraefikRouters{
					"UDPRouter0": {
						EntryPoints: []string{"foobar", "foobar"},
						Service:     "foobar",
					},
					"UDPRouter1": {
						EntryPoints: []string{"foobar", "foobar"},
						Service:     "foobar",
					},
				},
				Services: map[string]UdpTraefikServiceConfig{
					"UDPService01": {
						LoadBalancer: &UdpLoadBalancer{
							Servers: []UdpLoadBalancerServer{
								{Address: "foobar"},
								{Address: "foobar"},
							},
						},
					},
					"UDPService02": {
						Weighted: &UdpWeighted{
							Services: []UdpWeightedService{
								{Name: "foobar", Weight: 42},
								{Name: "foobar", Weight: 42},
							},
						},
					},
				},
			},
			Tls: TlsTraefik{
				Certificates: []TlsTraefikCertificates{
					{
						CertFile: "foobar",
						KeyFile:  "foobar",
						Stores: []string{
							"foobar", "foobar",
						},
					},
					{
						CertFile: "foobar",
						KeyFile:  "foobar",
						Stores: []string{
							"foobar", "foobar",
						},
					},
				},
				Options: map[string]TlsTraefikOptions{
					"Options0": {
						MinVersion: "foobar",
						MaxVersion: "foobar",
						CipherSuites: []string{
							"foobar",
							"foobar",
						},
						CurvePreferences: []string{
							"foobar",
							"foobar",
						},
						ClientAuth: TlsTraefikOptionsClientAuth{
							CaFiles: []string{
								"foobar",
								"foobar",
							},
							ClientAuthType: "foobar",
						},
						SniStrict:                true,
						PreferServerCipherSuites: true,
					},
					"Options1": {
						MinVersion: "foobar",
						MaxVersion: "foobar",
						CipherSuites: []string{
							"foobar",
							"foobar",
						},
						CurvePreferences: []string{
							"foobar",
							"foobar",
						},
						ClientAuth: TlsTraefikOptionsClientAuth{
							CaFiles: []string{
								"foobar",
								"foobar",
							},
							ClientAuthType: "foobar",
						},
						SniStrict:                true,
						PreferServerCipherSuites: true,
					},
				},
				Stores: map[string]TlsTraefikStores{
					"Store0": {
						DefaultCertificate: TlsTraefikStoresDefaultCertificate{
							CertFile: "foobar",
							KeyFile:  "foobar",
						},
					},
					"Store1": {
						DefaultCertificate: TlsTraefikStoresDefaultCertificate{
							CertFile: "foobar",
							KeyFile:  "foobar",
						},
					},
				},
			},
		},
	}
	return
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *TraefikConfig) parseTag(tagValue string) (name string, omitempty, ignore bool) {
	tagList := strings.Split(tagValue, ",")
	name = tagList[0]
	if name == "" {
		return
	}
	for _, s := range tagList {
		if s == "omitempty" {
			omitempty = true
		}
		if s == "-" {
			ignore = true
		}
	}
	return
}
func (r *TraefikConfig) reflectName(tag, prefixName string, fieldStruct *reflect.StructField) (name, newPrefixName string, omitempty, ignore bool) {
	name, omitempty, ignore = r.parseTag(tag)
	if name == "" {
		name = fieldStruct.Name
	}
	newPrefixName = r.getPrefixName(prefixName, name)
	return
}
func (r *TraefikConfig) getPrefixName(prefixName, name string) (res string) {
	res = fmt.Sprintf("%s/%s", prefixName, name)
	return
}
func (r *TraefikConfig) ergodic(data interface{}, prefixName string, res map[string]string) {

	values := reflect.ValueOf(data)
	if r.generalType(values) {
		res[prefixName] = values.String()
		return
	}

	tagValue := "key_value"
	types := reflect.TypeOf(data)

	var fieldStruct reflect.StructField
	var valueStruct reflect.Value
	fieldNum := types.NumField()
	for i := 0; i < fieldNum; i++ {

		fieldStruct = types.Field(i)
		tag := fieldStruct.Tag.Get(tagValue)
		name, newPrefixName, omitempty, ignore := r.reflectName(tag, prefixName, &fieldStruct)
		if ignore {
			continue
		}
		valueStruct = values.Field(i)
		kind := valueStruct.Kind()
		if omitempty && r.isBlank(valueStruct) {
			continue
		}
		var iValue reflect.Value
		switch kind {
		case reflect.Ptr: //如果是指针
			r.ergodic(valueStruct.Elem().Interface(), newPrefixName, res)
		case reflect.Struct: //如果是结构体
			r.ergodic(valueStruct.Interface(), newPrefixName, res)
		case reflect.Map:
			for _, key := range valueStruct.MapKeys() {
				iValue = valueStruct.MapIndex(key)
				keyName := r.getPrefixName(newPrefixName, key.String())
				if r.generalType(iValue) {
					res[keyName] = iValue.String()
					continue
				}
				r.ergodic(valueStruct.MapIndex(key).Interface(), keyName, res)
			}
		case reflect.Slice:
			for j := 0; j < valueStruct.Len(); j++ {
				iValue = valueStruct.Index(j)
				keyName := r.getPrefixName(newPrefixName, strconv.Itoa(j))
				if r.generalType(iValue) {
					res[keyName] = iValue.String()
					continue
				}
				r.ergodic(iValue.Interface(), keyName, res)
			}

		case reflect.Bool:
			if name == "value" {
				//如果是空,则跳过了
				res[prefixName] = strconv.FormatBool(valueStruct.Bool())
				continue
			}

			//如果是空,则跳过了
			res[newPrefixName] = strconv.FormatBool(valueStruct.Bool())

		case reflect.Interface:
			//fmt.Println(name,valueStruct.Interface().(type))
			r.ergodic(valueStruct.Interface(), newPrefixName, res)
		default:
			if name == "value" {
				//如果是空,则跳过了
				res[prefixName] = r.formatValueToString(name, valueStruct, fieldStruct)
				continue
			}
			//如果是空,则跳过了
			res[newPrefixName] = r.formatValueToString(name, valueStruct, fieldStruct)
		}
	}
}

func (r *TraefikConfig) formatValueToString(name string, valueStruct reflect.Value, fieldStruct reflect.StructField) (res string) {
	switch fieldStruct.Type.String() {
	case "time.Duration":
		res = valueStruct.Interface().(time.Duration).String()
	case "int":
		res = fmt.Sprintf("%v", valueStruct.Int())
	default:
		res = valueStruct.String()
	}
	return
}

func (r *TraefikConfig) generalType(value reflect.Value) (res bool) {
	switch value.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		res = true
	case reflect.Ptr:
		v := value.Elem()
		res = r.generalType(v)
	}
	return
}

func (r *TraefikConfig) isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func (r *TraefikConfig) ToKV() (res map[string]string) {
	res = make(map[string]string, 200)
	prefixName := "traefik"
	r.ergodic(r.TraefikDynamic, prefixName, res)
	return
}
func (r *TraefikConfig) KVShow() {
	stringMap := r.ToKV()
	var keySortSlice KeySortSlice = make([]string, 0, len(stringMap))
	for i, _ := range stringMap {
		keySortSlice = append(keySortSlice, i)
	}
	sort.Sort(&keySortSlice)
	for _, s := range keySortSlice {
		fmt.Printf("%s\t%s\n", s, stringMap[s])
	}

}

type KeySortSlice []string

func (k KeySortSlice) Len() int {
	return len(k)
}

func (k KeySortSlice) Less(i, j int) bool {
	return k[i] < k[j]
}

func (k KeySortSlice) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (r *TraefikConfig) writeToFile(filename string, data []byte) (n int, err1 error) {

	var f *os.File
	f, err1 = os.Create(filename) //创建文件
	if err1 != nil {
		return
	}
	if n, err1 = f.Write(data); err1 != nil { //写入文件(字符串)
		return
	}
	return
}

func (r *TraefikConfig) AppendToFile(filename string) (err error) {
	var data []byte
	data, err = yaml.Marshal(r.TraefikDynamic)
	r.writeToFile(filename, data)
	return
}

func (r *TraefikConfig) checkFileIsExist(fileName string) (res bool) {
	res = true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		res = false
	}
	return
}
func (r *TraefikConfig) ResultToJson() (res string, err error) {
	d, err := json.Marshal(r.MapValue)
	if err != nil {
		return
	}
	res = string(d)
	return
}

func (r *TraefikConfig) ResultToYaml() (res string, err error) {
	d, err := yaml.Marshal(&r.MapValue)
	if err != nil {
		return
	}
	res = string(d)
	return
}
