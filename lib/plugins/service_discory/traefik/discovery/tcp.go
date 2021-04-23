// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

type TcpTraefik struct {
	Routers  map[string]TcpTraefikRouters       `yaml:"routers,omitempty" key_value:"routers,omitempty"`
	Services map[string]TcpTraefikServiceConfig `yaml:"services,omitempty" key_value:"services,omitempty"`
}
type TcpTraefikRouters struct {
	EntryPoints []string `json:"entryPoints" yaml:"entryPoints,omitempty" key_value:"entryPoints,omitempty"`
	Service     string   `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	Rule        string   `json:"rule" yaml:"rule,omitempty" key_value:"rule,omitempty"`
	Tls         *TCPTls  `json:"tls" yaml:"tls,omitempty" key_value:"tls,omitempty"`
	//ServiceList []TcpServiceConfig `json:"service_list" yaml:"service_list,omitempty" key_value:"service_list,omitempty"`
}
type TCPTls struct {
	PassThrough bool   `json:"passthrough" yaml:"passthrough,omitempty" key_value:"passthrough,omitempty"`
	Options     string `json:"options" yaml:"options,omitempty" key_value:"options,omitempty"`
	//Value        bool               `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	CertResolver string             `json:"certResolver" yaml:"certResolver,omitempty" key_value:"certResolver,omitempty"`
	Domains      []TcpDomainTlsItem `json:"domains" yaml:"domains,omitempty" key_value:"domains,omitempty"`
}
type TcpTraefikServiceConfig struct {
	LoadBalancer *TcpLoadBalancer `yaml:"loadBalancer,omitempty" key_value:"loadBalancer,omitempty"`
	Weighted     *TcpWeighted     `yaml:"weighted,omitempty" key_value:"weighted,omitempty"`
}

type TcpDomainTlsItem struct {
	Main        string   `json:"main" yaml:"main,omitempty" key_value:"main,omitempty"`
	DomainValue string   `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	Sans        []string `json:"sans" yaml:"sans,omitempty" key_value:"sans,omitempty"`
}
type TcpLoadBalancer struct {
	TerminationDelay int                     `json:"terminationDelay" yaml:"terminationDelay,omitempty" key_value:"terminationDelay,omitempty"`
	ProxyProtocol    TcpProxyProtocol        `json:"proxyProtocol" yaml:"proxyProtocol,omitempty" key_value:"proxyProtocol,omitempty"`
	Servers          []TcpLoadBalancerServer `json:"servers" yaml:"servers,omitempty" key_value:"servers,omitempty"`
}
type TcpWeighted struct {
	Services []TcpWeightedService `yaml:"services,omitempty" key_value:"services,omitempty"`
}

type TcpProxyProtocol struct {
	Version int `json:"version" yaml:"version,omitempty" key_value:"version,omitempty"`
}
type TcpWeightedService struct {
	Name   string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Weight int    `json:"weight" yaml:"weight,omitempty" key_value:"weight,omitempty"`
}
type TcpLoadBalancerServer struct {
	Address string `json:"address" yaml:"address,omitempty" key_value:"address,omitempty"`
}
