// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/plugins/traefik/middleware_traefik"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TraefikDynamic struct {
	Http HttpTraefik `yaml:"http,omitempty" key_value:"http,omitempty"`
	Tcp  TcpTraefik  `yaml:"tcp,omitempty" key_value:"tcp,omitempty"`
}

type TraefikConfig struct {
	TraefikDynamic
	MapValue []KeyValue //数据结果
}

func NewTraefikConfig() (res *TraefikConfig) {
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
							Value:        true,
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
							Value:        true,
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
								Value: true,
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
				},
			},
			Tcp: TcpTraefik{
				Routers: map[string]TcpTraefikRouters{
					"<router_name>": TcpTraefikRouters{
						EntryPoints: []string{"ep1", "ep2"},
						Rule:        "HostSNI(`example.com`)",
						Service:     "myservice",
						Tls: &TCPTls{
							Value:        true,
							CertResolver: "myresolver",
							Domains: []TcpDomainTlsItem{
								{
									Main: "example.org",
									Sans: []string{"test.example.org", "dev.example.org"},
								},
							},
							Options:     "foobar",
							PassThrough: true,
						},
					},
				},
				Services: map[string]TcpTraefikServiceConfig{
					"<service_name>": {
						Weighted: &TcpWeighted{
							Services: []TcpWeightedService{
								{
									Name:   "foobar",
									Weight: 42,
								},
							},
						},
						LoadBalancer: &TcpLoadBalancer{
							ProxyProtocol: TcpProxyProtocol{
								Version: "1",
							},
							TerminationDelay: 100,
							Servers: []TcpLoadBalancerServer{
								{
									Url: "xx.xx.xx.xx:xx",
								},
							},
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
	for i, i2 := range stringMap {
		fmt.Printf("%s \t %s \n", i, i2)
	}
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
