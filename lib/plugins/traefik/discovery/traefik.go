// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type TraefikConfig struct {
	RouterHttpConfig []TraefikService    //http请求参数
	RouterTcpConfig  []TraefikTcpService //TCP请求参数

	MapValue []KeyValue //数据结果
}

func NewTraefikConfig() (res *TraefikConfig) {
	res = &TraefikConfig{
		RouterHttpConfig: []TraefikService{
			{
				RouterName: "myrouter",
				Rule:       "Host(`example.com`)",
				EntryPoints: []string{
					"web",
					"websecure",
				},

				Middlewares: []string{

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

func (r *TraefikConfig) Action() {
	{ //init router
		r.router()
		r.entryPoints()
		r.middlewares()
		r.service()
		r.priority()
		r.tls()
	}
	{ //init service list
		r.serviceList()
	}

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
