// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package service_discory

type ServerRegistry struct {
	Endpoints     []string `json:"endpoints" yaml:"endpoints"`
	Dir           string   `json:"dir" yaml:"dir"`
	LockKey       string   `json:"lock_key" yaml:"lockkey"`
	EtcdEndPoints []string `json:"etcd_end_points" yaml:"etcdendpoints"`
	Host          string   `json:"host" yaml:"host"`
}
