// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

type TraefikTcpService struct {
	RouterName  string             `json:"router_name"`
	Rule        string             `json:"rule"`
	Service     string             `json:"service"`
	Tls         *Tls               `json:"tls"`
	ServiceList []ServiceTcpConfig `json:"service_list"`
}
type ServiceTcpConfig struct {
	ServiceName  string          `json:"service_name"`
	LoadBalancer LoadTcpBalancer `json:"loadbalancer"`
	Weighted     TcpWeighted     `json:"weighted"`
}
type TcpWeighted struct {
	Services []WeightedService `json:"services"`
}
type LoadTcpBalancer struct {
	Servers          []LoadBalancerServer `json:"servers"`
	TerminationDelay string               `json:"terminationdelay"`
	ProxyProtocol    struct {
		Version string `json:"version"`
	} `json:"proxyprotocol"`
}
