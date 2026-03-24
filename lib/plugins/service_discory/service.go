// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package service_discory

type ServerRegistry struct {
	ETCD *struct {
		Endpoints     []string `json:"endpoints" yaml:"endpoints"`
		Dir           string   `json:"dir" yaml:"dir"`
		LockKey       string   `json:"lock_key" yaml:"lockkey"`
		EtcdEndPoints []string `json:"etcd_end_points" yaml:"etcdendpoints"`
		Host          string   `json:"host" yaml:"host"`
	} `json:"etcd" yaml:"etcd"` //etcd 配置信息
	
	Consul *struct {
		Endpoints []string `json:"endpoints" yaml:"endpoints"` //consul 注册地址
	} `json:"consul" yaml:"consul"` //consul 配置信息
}
