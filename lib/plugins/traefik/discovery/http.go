// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import "time"

type TraefikHttpMiddleware struct {
	Plugin interface{} `yaml:"plugin,omitempty" key_value:"plugin,omitempty"`
}
type TraefikHttpRouters struct {
	EntryPoints []string `yaml:"entry_points,omitempty" key_value:"entrypoints,omitempty"`
	Rule        string   `yaml:"rule,omitempty" key_value:"rule,omitempty"`
	Service     string   `yaml:"service,omitempty" key_value:"service,omitempty"`
	Middlewares []string `yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`

	//ServiceList []ServiceConfig `yaml:"service_list,omitempty"`
	//Priority    int             `yaml:"priority,omitempty"`
	//Tls         *Tls            `yaml:"tls,omitempty"`
}
type TraefikHttpServiceConfig struct {
	LoadBalancer *LoadBalancer `yaml:"loadBalancer,omitempty" key_value:"loadbalancer,omitempty"`
	Mirroring    *Mirroring    `yaml:"mirroring,omitempty" key_value:"mirroring,omitempty"`
	Weighted     *Weighted     `yaml:"weighted,omitempty" key_value:"weighted,omitempty"`
}
type Mirroring struct {
	Service string    `yaml:"service,omitempty" key_value:"service,omitempty"`
	Mirrors []Mirrors `yaml:"mirrors,omitempty" key_value:"mirrors,omitempty"`
}
type Mirrors struct {
	Name    string `yaml:"name,omitempty" key_value:"name,omitempty"`
	Percent int    `yaml:"percent,omitempty" key_value:"percent,omitempty"`
}
type Weighted struct {
	Services []WeightedService `yaml:"services,omitempty" key_value:"services,omitempty"`
	Sticky   *Sticky           `yaml:"sticky,omitempty" key_value:"sticky,omitempty"`
}
type WeightedService struct {
	Name   string `yaml:"name,omitempty" key_value:"name,omitempty"`
	Weight int    `yaml:"weight,omitempty" key_value:"weight,omitempty"`
}

type LoadBalancer struct {
	Servers            []LoadBalancerServer `yaml:"servers,omitempty" key_value:"servers,omitempty"`
	PassHostHeader     bool                 `yaml:"passHostHeader,omitempty" key_value:"passhostheader,omitempty"`
	ServersTransport   string               `yaml:"serverstransport,omitempty" key_value:"serverstransport,omitempty"`
	HealthCheck        *HealthCheck         `yaml:"healthCheck,omitempty" key_value:"healthcheck,omitempty"`
	Sticky             *Sticky              `yaml:"sticky,omitempty" key_value:"sticky,omitempty"`
	ResponseForwarding ResponseForwarding   `yaml:"responseforwarding,omitempty" key_value:"responseforwarding,omitempty"`
}
type ResponseForwarding struct {
	FlushInterval time.Duration `json:"flushinterval" yaml:"flushinterval,omitempty" key_value:"flushinterval,omitempty"`
}

type Sticky struct {
	Value  bool    `yaml:"value,omitempty" key_value:"value,omitempty"`
	Cookie *Cookie `yaml:"cookie,omitempty" key_value:"cookie,omitempty"`
}
type Cookie struct {
	HttpOnly bool   `yaml:"httponly,omitempty" key_value:"httponly,omitempty"`
	Name     string `yaml:"name,omitempty" key_value:"name,omitempty"`
	Secure   bool   `yaml:"secure,omitempty" key_value:"secure,omitempty"`
	SameSite string `yaml:"samesite,omitempty" key_value:"samesite,omitempty"`
}
type HealthCheck struct {
	Headers  map[string]string `json:"headers" yaml:"headers,omitempty" key_value:"headers,omitempty"`
	Hostname string            `json:"hostname" yaml:"hostname,omitempty" key_value:"hostname,omitempty"`
	Interval time.Duration            `json:"interval" yaml:"interval,omitempty" key_value:"interval,omitempty"`
	Path     string            `json:"path" yaml:"path,omitempty" key_value:"path,omitempty"`
	Port     int               `json:"port" yaml:"port,omitempty" key_value:"port,omitempty"`
	Scheme   string            `json:"scheme" yaml:"scheme,omitempty" key_value:"scheme,omitempty"`
	Timeout  time.Duration            `json:"timeout" yaml:"timeout,omitempty" key_value:"timeout,omitempty"`
}
type LoadBalancerServer struct {
	Url string `json:"url" yaml:"url,omitempty" key_value:"url,omitempty"`
}
type Tls struct {
	Value        bool            `json:"value"`
	CertResolver string          `json:"certresolver"`
	Domains      []TlsDomainItem `json:"domains"`
	Options      string          `json:"options"`
}

//type TlsDomain struct {
//	List []TlsDomainItem
//}
type TlsDomainItem struct {
	Domain      string   `json:"domain"`
	DomainValue string   `json:"value"`
	Sans        []string `json:"sans"`
}
