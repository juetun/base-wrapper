// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

type UdpTraefik struct {
	Routers  map[string]UdpTraefikRouters       `json:"routers" yaml:"routers,omitempty" key_value:"routers,omitempty"`
	Services map[string]UdpTraefikServiceConfig `yaml:"services,omitempty" key_value:"services,omitempty"`
}
type UdpTraefikRouters struct {
	EntryPoints []string `json:"entryPoints" yaml:"entryPoints,omitempty" key_value:"entryPoints,omitempty"`
	Service     string   `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
}
type UdpTraefikServiceConfig struct {
	LoadBalancer *UdpLoadBalancer `yaml:"loadBalancer,omitempty" key_value:"loadBalancer,omitempty"`
	Weighted     *UdpWeighted     `yaml:"weighted,omitempty" key_value:"weighted,omitempty"`
}
type UdpLoadBalancer struct {
	Servers []UdpLoadBalancerServer `json:"servers" yaml:"servers,omitempty" key_value:"servers,omitempty"`
}

type UdpLoadBalancerServer struct {
	Address string `json:"address" yaml:"address,omitempty" key_value:"address,omitempty"`
}

type UdpWeighted struct {
	Services []UdpWeightedService `yaml:"services,omitempty" key_value:"services,omitempty"`
}

type UdpWeightedService struct {
	Name   string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Weight int    `json:"weight" yaml:"weight,omitempty" key_value:"weight,omitempty"`
}
