// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TraefikDynamic struct {
	Http TraefikHttp `yaml:"http,omitempty" key_value:"http,omitempty"`
	Tcp  TraefikTcp  `yaml:"tcp,omitempty" key_value:"tcp,omitempty"`
}
type TraefikTcp struct {
	RouterTcpConfig map[string]TraefikTcpService //TCP请求参数
}
type TraefikHttp struct {
	Routers     map[string]TraefikHttpRouters       `yaml:"routers,omitempty" key_value:"routers,omitempty"`
	Services    map[string]TraefikHttpServiceConfig `yaml:"services,omitempty" key_value:"services,omitempty"`
	Middlewares map[string]TraefikHttpMiddleware    `yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`
}

type TraefikConfig struct {
	TraefikDynamic
	MapValue []KeyValue //数据结果
}

func NewTraefikConfig() (res *TraefikConfig) {
	res = &TraefikConfig{
		TraefikDynamic: TraefikDynamic{
			Http: TraefikHttp{
				Routers: map[string]TraefikHttpRouters{
					"<router_name>": {
						Rule: "Host(`api.test.com`) && PathPrefix(`/api-user`)",
						EntryPoints: []string{
							"web",
							"websecure",
						},
						Service:     "api-user",
						Middlewares: []string{"my-plugin"},
					},
				},
				Services: map[string]TraefikHttpServiceConfig{
					"<service_name>": {
						LoadBalancer: &LoadBalancer{
							Servers: []LoadBalancerServer{
								{
									Url: "http://localhost:8093",
								},
							},
							ResponseForwarding: ResponseForwarding{
								FlushInterval: 10 * time.Hour,
							},
							PassHostHeader: true,
							HealthCheck: &HealthCheck{
								Headers:  map[string]string{"Content-type": "application/json"},
								Hostname: "example.org",
								Interval: 10 * time.Second,
								Path:     "/foo",
								Port:     8080,
								Scheme:   "https",
								Timeout:  15 * time.Second,
							},
							Sticky: &Sticky{
								Value: true,
								Cookie: &Cookie{
									HttpOnly: true,
								},
							},
						},
					},
				},
				Middlewares: map[string]TraefikHttpMiddleware{
					"my-plugin": {
						Plugin: "123123",
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
