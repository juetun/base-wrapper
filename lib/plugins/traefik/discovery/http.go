// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import "time"

type HttpTraefik struct {
	Routers     map[string]HttpTraefikRouters       `json:"routers" yaml:"routers,omitempty" key_value:"routers,omitempty"`
	Services    map[string]HttpTraefikServiceConfig `json:"services" yaml:"services,omitempty" key_value:"services,omitempty"`
	Middlewares map[string]HttpTraefikMiddleware    `json:"middlewares" yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`
}

type HttpTraefikMiddleware struct {
	Plugin interface{} `json:"plugin" yaml:"plugin,omitempty" key_value:"plugin,omitempty"`
}
type HttpTraefikRouters struct {
	EntryPoints []string `json:"entry_points" yaml:"entry_points,omitempty" key_value:"entrypoints,omitempty"`
	Rule        string   `json:"rule" yaml:"rule,omitempty" key_value:"rule,omitempty"`
	Service     string   `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	Middlewares []string `json:"middlewares" yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`

	//ServiceList []ServiceConfig `yaml:"service_list,omitempty"`
	Priority int      `json:"priority" yaml:"priority,omitempty" key_value:"priority,omitempty"`
	Tls      *HttpTls `json:"tls" yaml:"tls,omitempty" key_value:"tls,omitempty"`
}
type HttpTraefikServiceConfig struct {
	LoadBalancer *HttpLoadBalancer `json:"load_balancer" yaml:"loadBalancer,omitempty" key_value:"loadbalancer,omitempty"`
	Mirroring    *HttpMirroring    `json:"mirroring" yaml:"mirroring,omitempty" key_value:"mirroring,omitempty"`
	Weighted     *HttpWeighted     `json:"weighted" yaml:"weighted,omitempty" key_value:"weighted,omitempty"`
}

type HttpMirroring struct {
	Service string        `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	Mirrors []HttpMirrors `json:"mirrors" yaml:"mirrors,omitempty" key_value:"mirrors,omitempty"`
}
type HttpMirrors struct {
	Name    string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Percent int    `json:"percent" yaml:"percent,omitempty" key_value:"percent,omitempty"`
}
type HttpWeighted struct {
	Services []HttpWeightedService `json:"services" yaml:"services,omitempty" key_value:"services,omitempty"`
	Sticky   *HttpSticky           `json:"sticky" yaml:"sticky,omitempty" key_value:"sticky,omitempty"`
}

type HttpWeightedService struct {
	Name   string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Weight int    `json:"weight" yaml:"weight,omitempty" key_value:"weight,omitempty"`
}

type HttpLoadBalancer struct {
	Servers            []HttpLoadBalancerServer `json:"servers" yaml:"servers,omitempty" key_value:"servers,omitempty"`
	PassHostHeader     bool                     `json:"pass_host_header" yaml:"passHostHeader,omitempty" key_value:"passhostheader,omitempty"`
	ServersTransport   string                   `json:"servers_transport" yaml:"serverstransport,omitempty" key_value:"serverstransport,omitempty"`
	HealthCheck        *HttpHealthCheck         `json:"health_check" yaml:"healthCheck,omitempty" key_value:"healthcheck,omitempty"`
	Sticky             *HttpSticky              `json:"sticky" yaml:"sticky,omitempty" key_value:"sticky,omitempty"`
	ResponseForwarding HttpResponseForwarding   `json:"response_forwarding" yaml:"responseforwarding,omitempty" key_value:"responseforwarding,omitempty"`
}

type HttpResponseForwarding struct {
	FlushInterval time.Duration `json:"flush_interval" yaml:"flushinterval,omitempty" key_value:"flushinterval,omitempty"`
}

type HttpSticky struct {
	Value  bool        `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	Cookie *HttpCookie `json:"cookie" yaml:"cookie,omitempty" key_value:"cookie,omitempty"`
}
type HttpCookie struct {
	HttpOnly bool   `json:"http_only" yaml:"httponly,omitempty" key_value:"httponly,omitempty"`
	Name     string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Secure   bool   `json:"secure" yaml:"secure,omitempty" key_value:"secure,omitempty"`
	SameSite string `json:"same_site" yaml:"samesite,omitempty" key_value:"samesite,omitempty"`
}
type HttpHealthCheck struct {
	Headers  map[string]string `json:"headers" yaml:"headers,omitempty" key_value:"headers,omitempty"`
	Hostname string            `json:"hostname" yaml:"hostname,omitempty" key_value:"hostname,omitempty"`
	Interval time.Duration     `json:"interval" yaml:"interval,omitempty" key_value:"interval,omitempty"`
	Path     string            `json:"path" yaml:"path,omitempty" key_value:"path,omitempty"`
	Port     int               `json:"port" yaml:"port,omitempty" key_value:"port,omitempty"`
	Scheme   string            `json:"scheme" yaml:"scheme,omitempty" key_value:"scheme,omitempty"`
	Timeout  time.Duration     `json:"timeout" yaml:"timeout,omitempty" key_value:"timeout,omitempty"`
}
type HttpLoadBalancerServer struct {
	Url string `json:"url" yaml:"url,omitempty" key_value:"url,omitempty"`
}

type HttpTls struct {
	Value        bool            `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	CertResolver string          `json:"certresolver" yaml:"certresolver,omitempty" key_value:"certresolver,omitempty"`
	Domains      []HttpDomainTlsItem `json:"domains" yaml:"domains,omitempty" key_value:"domains,omitempty"`
	Options      string          `json:"options" yaml:"options,omitempty" key_value:"options,omitempty"`
}

type HttpDomainTlsItem struct {
	Main        string   `json:"main" yaml:"main,omitempty" key_value:"main,omitempty"`
	DomainValue string   `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
	Sans        []string `json:"sans" yaml:"sans,omitempty" key_value:"sans,omitempty"`
}
