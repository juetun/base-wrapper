// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery
type TcpTraefik struct {
	Routers  map[string]TcpTraefikRouters       `yaml:"routers,omitempty" key_value:"routers,omitempty"`
	Services map[string]TcpTraefikServiceConfig `yaml:"services,omitempty" key_value:"services,omitempty"`
}
type TcpTraefikRouters struct {
	EntryPoints []string           `json:"entry_points" yaml:"entry_points,omitempty" key_value:"entrypoints,omitempty"`
	Rule        string             `json:"rule" yaml:"rule,omitempty" key_value:"rule,omitempty"`
	Service     string             `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	Tls         *TCPTls            `json:"tls" yaml:"tls,omitempty" key_value:"tls,omitempty"`
	ServiceList []TcpServiceConfig `json:"service_list" yaml:"service_list,omitempty" key_value:"service_list,omitempty"`
}
type TCPTls struct {
	Value        bool               `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	CertResolver string             `json:"certresolver" yaml:"certresolver,omitempty" key_value:"certresolver,omitempty"`
	Domains      []TcpDomainTlsItem `json:"domains" yaml:"domains,omitempty" key_value:"domains,omitempty"`
	Options      string             `json:"options" yaml:"options,omitempty" key_value:"options,omitempty"`
	PassThrough  bool               `json:"passthrough" yaml:"passthrough,omitempty" key_value:"passthrough,omitempty"`
}
type TcpTraefikServiceConfig struct {
	LoadBalancer *TcpLoadBalancer `yaml:"loadBalancer,omitempty" key_value:"loadbalancer,omitempty"`
	Weighted     *TcpWeighted     `yaml:"weighted,omitempty" key_value:"weighted,omitempty"`
}
type TcpServiceConfig struct {
	ServiceName  string          `json:"service_name"`
	LoadBalancer LoadTcpBalancer `json:"loadbalancer"`
	Weighted     TcpWeighted     `json:"weighted"`
}

type LoadTcpBalancer struct {
	Servers          []TcpLoadBalancerServer `json:"servers"`
	TerminationDelay string               `json:"terminationdelay"`
	ProxyProtocol    struct {
		Version string `json:"version"`
	} `json:"proxyprotocol"`
}

//type TlsDomain struct {
//	List []TlsDomainItem
//}
type TcpDomainTlsItem struct {
	Main        string   `json:"main" yaml:"main,omitempty" key_value:"main,omitempty"`
	DomainValue string   `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	Sans        []string `json:"sans" yaml:"sans,omitempty" key_value:"sans,omitempty"`
}
type TcpLoadBalancer struct {
	Servers          []TcpLoadBalancerServer `json:"servers" yaml:"servers,omitempty" key_value:"servers,omitempty"`
	TerminationDelay int                     `json:"terminationdelay" yaml:"terminationdelay,omitempty" key_value:"terminationdelay,omitempty"`
	ProxyProtocol    TcpProxyProtocol        `json:"proxyprotocol" yaml:"proxyprotocol,omitempty" key_value:"proxyprotocol,omitempty"`
}
type TcpWeighted struct {
	Services []TcpWeightedService `yaml:"services,omitempty" key_value:"services,omitempty"`
}

type TcpProxyProtocol struct {
	Version string `json:"version" yaml:"version,omitempty" key_value:"version,omitempty"`
}
type TcpWeightedService struct {
	Name   string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Weight int    `json:"weight" yaml:"weight,omitempty" key_value:"weight,omitempty"`
}
type TcpLoadBalancerServer struct {
	Url string `json:"url" yaml:"url,omitempty" key_value:"url,omitempty"`
}