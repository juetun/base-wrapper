// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type HttpTraefik struct {
	Routers           map[string]HttpTraefikRouters           `json:"routers" yaml:"routers,omitempty" key_value:"routers,omitempty"`
	Services          map[string]HttpTraefikServiceConfig     `json:"services" yaml:"services,omitempty" key_value:"services,omitempty"`
	Middlewares       map[string]HttpTraefikMiddleware        `json:"middlewares" yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`
	ServersTransports map[string]HttpTraefikServersTransports `json:"serversTransports" yaml:"serversTransports,omitempty" key_value:"serversTransports,omitempty"`
}

func (r *HttpTraefik) MergeRouters(routerMap map[string]HttpTraefikRouters) {
	if len(routerMap) == 0 {
		return
	}
	if len(r.Routers) == 0 {
		r.Routers = make(map[string]HttpTraefikRouters, 1)
	}
	for s, it := range routerMap {
		if _, ok := r.Routers[s]; !ok {
			r.Routers[s] = it
			continue
		}
		r.Routers[s] = r.mergeRouter(r.Routers[s], it)
	}
}

func (r *HttpTraefik) MergeServices(serviceMap map[string]HttpTraefikServiceConfig) {
	if len(serviceMap) == 0 {
		return
	}
	if len(r.Services) == 0 {
		r.Services = make(map[string]HttpTraefikServiceConfig, 1)
	}
	for s, it := range serviceMap {
		if _, ok := r.Services[s]; !ok {
			r.Services[s] = it
			continue
		}
		r.Services[s] = r.mergeServices(r.Services[s], it)
	}
}

func (r *HttpTraefik) MergeMiddlewares(middlewareMap map[string]HttpTraefikMiddleware) {
	if len(middlewareMap) == 0 {
		return
	}
	if len(r.Middlewares) == 0 {
		r.Middlewares = make(map[string]HttpTraefikMiddleware, 1)
	}
	for s, it := range middlewareMap {
		r.Middlewares[s] = it
	}
}

func (r *HttpTraefik) MergeServersTransports(serversTransportsMap map[string]HttpTraefikServersTransports) {
	if len(serversTransportsMap) == 0 {
		return
	}
	if len(r.ServersTransports) == 0 {
		r.ServersTransports = make(map[string]HttpTraefikServersTransports, 1)
	}
	for s, it := range serversTransportsMap {
		if _, ok := r.ServersTransports[s]; !ok {
			r.ServersTransports[s] = it
			continue
		}
		r.ServersTransports[s] = r.mergeServersTransport(r.ServersTransports[s], it)
	}
}

//合并路由
func (r *HttpTraefik) mergeRouter(src, new HttpTraefikRouters) (res HttpTraefikRouters) {
	res = HttpTraefikRouters{
		Service:  new.Service,
		Rule:     new.Rule,
		Priority: new.Priority,
	}
	res.mergeTls(src.Tls, new.Tls)
	res.mergeEntryPoints(src.EntryPoints, new.EntryPoints)
	res.mergeMiddlewares(src.Middlewares, new.Middlewares)

	return
}

//合并服务
func (r *HttpTraefik) mergeServices(src, new HttpTraefikServiceConfig) (res HttpTraefikServiceConfig) {
	res = HttpTraefikServiceConfig{}
	res.mergeHttpWeighted(src.Weighted, new.Weighted)
	res.mergeLoadBalancer(src.LoadBalancer, new.LoadBalancer)
	res.mergeMirroring(src.Mirroring, new.Mirroring)
	return
}

//TODO 合并serverTransports
func (r *HttpTraefik) mergeServersTransport(src, new HttpTraefikServersTransports) (res HttpTraefikServersTransports) {
	res = HttpTraefikServersTransports{
		ServerName:          new.ServerName,
		InsecureSkipVerify:  new.InsecureSkipVerify,
		MaxIdleConnsPerHost: new.MaxIdleConnsPerHost,
	}
	res.mergeRootCAs(src.RootCAs, new.RootCAs)
	res.mergeCertificates(src.Certificates, new.Certificates)
	res.mergeForwardingTimeouts(src.ForwardingTimeouts, new.ForwardingTimeouts)
	return
}

func (r *HttpTraefikServersTransports) mergeRootCAs(src, new map[string]string) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	var l = ls + len(new)
	if len(r.RootCAs) == 0 {
		r.RootCAs = make(map[string]string, l)
	}
	var m = make(map[string]string, ls)
	for k, i2 := range src {
		m[i2] = k
	}
	var ind int
	for _, s := range new {
		if k, ok := m[s]; !ok {
			r.RootCAs[strconv.Itoa(ls+ind)] = s
			ind++
			continue
		} else {
			r.RootCAs[k] = s
		}
	}
}

func (r *HttpTraefikServersTransports) mergeCertificates(src, new map[string]HttpTraefikServersTransportsCertificates) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	var l = ls + len(new)
	if len(r.Certificates) == 0 {
		r.Certificates = make(map[string]HttpTraefikServersTransportsCertificates, l)
	}
	var m = make(map[string]string, ls)
	for k, i2 := range src {
		m[i2.ToString()] = k
	}
	var ind int
	for _, s := range new {
		if k, ok := m[s.ToString()]; !ok {
			r.Certificates[strconv.Itoa(ls+ind)] = s
			ind++
			continue
		} else {
			r.Certificates[k] = s
		}
	}
}

func (r *HttpTraefikServersTransports) mergeForwardingTimeouts(src, new *HttpTraefikServersTransportsForwardingTimeouts) {
	if new == nil {
		return
	}
	r.ForwardingTimeouts = &HttpTraefikServersTransportsForwardingTimeouts{
		DialTimeout:           new.DialTimeout,
		ResponseHeaderTimeout: new.ResponseHeaderTimeout,
		IdleConnTimeout:       new.IdleConnTimeout,
	}
}

type HttpTraefikServersTransports struct {
	ServerName          string                                              `json:"serverName" yaml:"serverName,omitempty" key_value:"serverName,omitempty"`
	InsecureSkipVerify  bool                                                `json:"insecureSkipVerify" yaml:"insecureSkipVerify,omitempty" key_value:"insecureSkipVerify,omitempty"`
	RootCAs             map[string]string                                   `json:"rootCAs" yaml:"rootCAs,omitempty" key_value:"rootCAs,omitempty"`
	Certificates        map[string]HttpTraefikServersTransportsCertificates `json:"certificates" yaml:"certificates,omitempty" key_value:"certificates,omitempty"`
	MaxIdleConnsPerHost int                                                 `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost,omitempty" key_value:"maxIdleConnsPerHost,omitempty"`
	ForwardingTimeouts  *HttpTraefikServersTransportsForwardingTimeouts     `json:"forwardingTimeouts" yaml:"forwardingTimeouts,omitempty" key_value:"forwardingTimeouts,omitempty"`
}
type HttpTraefikServersTransportsForwardingTimeouts struct {
	DialTimeout           time.Duration `json:"dialTimeout" yaml:"dialTimeout,omitempty" key_value:"dialTimeout,omitempty"`
	ResponseHeaderTimeout time.Duration `json:"responseHeaderTimeout" yaml:"responseHeaderTimeout,omitempty" key_value:"responseHeaderTimeout,omitempty"`
	IdleConnTimeout       time.Duration `json:"idleConnTimeout" yaml:"idleConnTimeout,omitempty" key_value:"idleConnTimeout,omitempty"`
}
type HttpTraefikServersTransportsCertificates struct {
	CertFile string `json:"certFile" yaml:"certFile,omitempty" key_value:"certFile,omitempty"`
	KeyFile  string `json:"keyFile" yaml:"keyFile,omitempty" key_value:"keyFile,omitempty"`
}

func (r *HttpTraefikServersTransportsCertificates) ToString() (res string) {
	bt, _ := json.Marshal(r)
	res = string(bt)
	return
}

type HttpTraefikMiddleware interface{} //`json:"plugin" yaml:"plugin,omitempty" key_value:"plugin,omitempty"`

type HttpTraefikRouters struct {
	Service     string            `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	Rule        string            `json:"rule" yaml:"rule,omitempty" key_value:"rule,omitempty"`
	Priority    int               `json:"priority" yaml:"priority,omitempty" key_value:"priority,omitempty"`
	Tls         *HttpTls          `json:"tls" yaml:"tls,omitempty" key_value:"tls,omitempty"`
	EntryPoints map[string]string `json:"entryPoints" yaml:"entryPoints,omitempty" key_value:"entryPoints,omitempty"`
	Middlewares map[string]string `json:"middlewares" yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`
}

func (r *HttpTraefikRouters) mergeTls(src, new *HttpTls) {
	if new == nil {
		return
	}
	if r.Tls == nil {
		r.Tls = &HttpTls{}
	}
	r.Tls.Options = new.Options
	r.Tls.CertResolver = new.CertResolver
	if src == nil {
		src = &HttpTls{}
	}
	r.Tls.mergeDomains(src.Domains, new.Domains)
	return
}
func (r *HttpTraefikRouters) mergeMiddlewares(src, new map[string]string) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	var l = ls + len(new)
	if len(r.Middlewares) == 0 {
		r.Middlewares = make(map[string]string, l)
	}
	var m = make(map[string]string, ls)
	for k, i2 := range src {
		m[i2] = k
	}
	var ind int
	for _, s := range new {
		if k, ok := m[s]; !ok {
			r.Middlewares[strconv.Itoa(ls+ind)] = s
			ind++
			continue
		} else {
			r.Middlewares[k] = s
		}
	}
}
func (r *HttpTraefikRouters) mergeEntryPoints(src, new map[string]string) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	var l = ls + len(new)
	if len(r.EntryPoints) == 0 {
		r.EntryPoints = make(map[string]string, l)
	}
	var m = make(map[string]string, ls)
	for k, i2 := range src {
		m[i2] = k
	}
	var ind int
	for _, s := range new {
		if loc, ok := m[s]; !ok {
			r.EntryPoints[strconv.Itoa(ls+ind)] = s
			ind++
			continue
		} else {
			r.EntryPoints[loc] = s
		}
	}
}

type HttpTraefikServiceConfig struct {
	LoadBalancer *HttpLoadBalancer `json:"loadBalancer" yaml:"loadBalancer,omitempty" key_value:"loadBalancer,omitempty"`
	Mirroring    *HttpMirroring    `json:"mirroring" yaml:"mirroring,omitempty" key_value:"mirroring,omitempty"`
	Weighted     *HttpWeighted     `json:"weighted" yaml:"weighted,omitempty" key_value:"weighted,omitempty"`
}

func (r *HttpTraefikServiceConfig) mergeLoadBalancer(src, new *HttpLoadBalancer) {
	if new == nil {
		return
	}
	if src == nil {
		src = &HttpLoadBalancer{}
	}
	r.LoadBalancer = &HttpLoadBalancer{
		PassHostHeader:   new.PassHostHeader,
		ServersTransport: new.ServersTransport,
	}
	r.LoadBalancer.mergeSticky(src.Sticky, new.Sticky)
	r.LoadBalancer.mergeServer(src.Servers, new.Servers)
	r.LoadBalancer.mergeHealthCheck(src.HealthCheck, new.HealthCheck)
	r.LoadBalancer.mergeResponseForwarding(src.ResponseForwarding, new.ResponseForwarding)

}
func (r *HttpLoadBalancer) mergeResponseForwarding(src, new *HttpResponseForwarding) {
	if new == nil {
		return
	}
	if src == nil {
		src = &HttpResponseForwarding{}
	}
	r.ResponseForwarding = &HttpResponseForwarding{
		FlushInterval: new.FlushInterval,
	}
}
func (r *HttpTraefikServiceConfig) mergeMirroring(src, new *HttpMirroring) {
	if new == nil {
		return
	}
	r.Mirroring = &HttpMirroring{
		Service:     new.Service,
		MaxBodySize: new.MaxBodySize,
	}
	if src == nil {
		src = &HttpMirroring{}
	}
	r.Mirroring.mergeMirrors(src.Mirrors, new.Mirrors)

}
func (r *HttpTraefikServiceConfig) mergeHttpWeighted(src, new *HttpWeighted) {
	if new == nil {
		return
	}
	if src == nil {
		src = &HttpWeighted{}
	}
	r.Weighted = &HttpWeighted{}
	r.Weighted.mergeServices(src.Services, new.Services)
	r.Weighted.mergeSticky(src.Sticky, new.Sticky)

}

type HttpMirroring struct {
	Service     string                 `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	MaxBodySize int                    `json:"maxBodySize" yaml:"maxBodySize,omitempty" key_value:"maxBodySize,omitempty"`
	Mirrors     map[string]HttpMirrors `json:"mirrors" yaml:"mirrors,omitempty" key_value:"mirrors,omitempty"`
}

func (r *HttpMirroring) mergeMirrors(src, new map[string]HttpMirrors) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	l := ls + len(new)
	r.Mirrors = make(map[string]HttpMirrors, l)
	var m = make(map[string]string, l)
	for k, item := range src {
		m[item.ToString()] = k
	}
	var ind int
	for _, item := range new {
		if loc, ok := m[item.ToString()]; !ok {
			r.Mirrors[strconv.Itoa(ls+ind)] = item
			ind++
			continue
		} else {
			r.Mirrors[loc] = item
		}
	}
}

type HttpMirrors struct {
	Name    string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Percent int    `json:"percent" yaml:"percent,omitempty" key_value:"percent,omitempty"`
}

func (r *HttpMirrors) ToString() (res string) {
	bt, _ := json.Marshal(r)
	res = string(bt)
	return
}

type HttpWeighted struct {
	Services map[string]HttpWeightedService `json:"services" yaml:"services,omitempty" key_value:"services,omitempty"`
	Sticky   *HttpSticky                    `json:"sticky" yaml:"sticky,omitempty" key_value:"sticky,omitempty"`
}

func (r *HttpWeighted) mergeSticky(src, new *HttpSticky) {
	if new == nil {
		return
	}
	if src == nil {
		src = &HttpSticky{}
	}
	r.Sticky = &HttpSticky{}
	r.Sticky.mergeCookie(src.Cookie, new.Cookie)
}
func (r *HttpWeighted) mergeServices(src, new map[string]HttpWeightedService) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	l := ls + len(new)
	r.Services = make(map[string]HttpWeightedService, l)
	var m = make(map[string]string, l)
	for k, item := range src {
		m[item.ToString()] = k
	}
	var ind int
	for _, item := range new {
		if k, ok := m[item.ToString()]; !ok {
			r.Services[strconv.Itoa(ls+ind)] = item
			ind++
			continue
		} else {
			r.Services[k] = item
		}
	}
}

type HttpWeightedService struct {
	Name   string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Weight int    `json:"weight" yaml:"weight,omitempty" key_value:"weight,omitempty"`
}

func (r *HttpWeightedService) ToString() (res string) {
	bt, _ := json.Marshal(r)
	res = string(bt)
	return
}

type HttpLoadBalancer struct {
	Sticky             *HttpSticky                       `json:"sticky" yaml:"sticky,omitempty" key_value:"sticky,omitempty"`
	Servers            map[string]HttpLoadBalancerServer `json:"servers" yaml:"servers,omitempty" key_value:"servers,omitempty"`
	HealthCheck        *HttpHealthCheck                  `json:"healthCheck" yaml:"healthCheck,omitempty" key_value:"healthCheck,omitempty"`
	PassHostHeader     bool                              `json:"passHostHeader" yaml:"passHostHeader,omitempty" key_value:"passHostHeader,omitempty"`
	ResponseForwarding *HttpResponseForwarding           `json:"responseForwarding" yaml:"responseForwarding,omitempty" key_value:"responseForwarding,omitempty"`
	ServersTransport   string                            `json:"serversTransport" yaml:"serversTransport,omitempty" key_value:"serversTransport,omitempty"`
}

func (r *HttpLoadBalancer) mergeSticky(src, new *HttpSticky) {
	if new == nil {
		return
	}
	if src == nil {
		src = &HttpSticky{}
	}
	r.Sticky = &HttpSticky{}
	r.Sticky.mergeCookie(src.Cookie, new.Cookie)
}

func (r *HttpLoadBalancer) mergeHealthCheck(src, new *HttpHealthCheck) {
	if new == nil {
		return
	}
	r.HealthCheck = &HttpHealthCheck{
		Scheme:          new.Scheme,
		Path:            new.Path,
		Port:            new.Port,
		Hostname:        new.Hostname,
		FollowRedirects: new.FollowRedirects,
		Interval: new.Interval,
		Timeout:  new.Timeout,
	}
	if src == nil {
		src = &HttpHealthCheck{}
	}
	r.HealthCheck.mergeHeaders(src.Headers, new.Headers)
}

func (r *HttpLoadBalancer) mergeServer(src, new map[string]HttpLoadBalancerServer) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	l := ls + len(new)
	r.Servers = make(map[string]HttpLoadBalancerServer, l)
	var m = make(map[string]string, l)
	for k, item := range src {
		m[item.ToString()] = k
	}
	var ind int
	for _, item := range new {
		if k, ok := m[item.ToString()]; !ok {
			r.Servers[strconv.Itoa(ls+ind)] = item
			ind++
			continue
		} else {
			r.Servers[k] = item
		}
	}
}

type HttpResponseForwarding struct {
	FlushInterval time.Duration `json:"flushInterval" yaml:"flushInterval,omitempty" key_value:"flushInterval,omitempty"`
}

type HttpSticky struct {
	Cookie *HttpCookie `json:"cookie" yaml:"cookie,omitempty" key_value:"cookie,omitempty"`
}

func (r *HttpSticky) mergeCookie(src, new *HttpCookie) {
	if new == nil {
		return
	}
	if src == nil {
		src = &HttpCookie{}
	}
	r.Cookie = &HttpCookie{
		Name:     new.Name,
		Secure:   new.Secure,
		HttpOnly: new.HttpOnly,
		SameSite: new.SameSite,
	}
}

type HttpCookie struct {
	Name     string `json:"name" yaml:"name,omitempty" key_value:"name,omitempty"`
	Secure   bool   `json:"secure" yaml:"secure,omitempty" key_value:"secure,omitempty"`
	HttpOnly bool   `json:"httpOnly" yaml:"httpOnly,omitempty" key_value:"httpOnly,omitempty"`
	SameSite string `json:"sameSite" yaml:"sameSite,omitempty" key_value:"sameSite,omitempty"`
}

type HttpHealthCheck struct {
	Scheme          string            `json:"scheme" yaml:"scheme,omitempty" key_value:"scheme,omitempty"`
	Path            string            `json:"path" yaml:"path,omitempty" key_value:"path,omitempty"`
	Port            string            `json:"port" yaml:"port,omitempty" key_value:"port,omitempty"`
	Hostname        string            `json:"hostname" yaml:"hostname,omitempty" key_value:"hostname,omitempty"`
	FollowRedirects bool              `json:"followRedirects" yaml:"followRedirects,omitempty" key_value:"followRedirects,omitempty"`
	Headers         map[string]string `json:"headers" yaml:"headers,omitempty" key_value:"headers,omitempty"`
	Interval        time.Duration     `json:"interval" yaml:"interval,omitempty" key_value:"interval,omitempty"`
	Timeout         time.Duration     `json:"timeout" yaml:"timeout,omitempty" key_value:"timeout,omitempty"`
}

//实现 HttpHealthCheck结构体 json反序列化方法
func (r *HttpHealthCheck) UnmarshalJSON(data []byte) (err error) {
	type httpHealthCheck struct {
		Scheme          string            `json:"scheme" yaml:"scheme,omitempty" key_value:"scheme,omitempty"`
		Path            string            `json:"path" yaml:"path,omitempty" key_value:"path,omitempty"`
		Port            string            `json:"port" yaml:"port,omitempty" key_value:"port,omitempty"`
		Hostname        string            `json:"hostname" yaml:"hostname,omitempty" key_value:"hostname,omitempty"`
		FollowRedirects bool              `json:"followRedirects" yaml:"followRedirects,omitempty" key_value:"followRedirects,omitempty"`
		Headers         map[string]string `json:"headers" yaml:"headers,omitempty" key_value:"headers,omitempty"`

		Interval string `json:"interval" yaml:"interval,omitempty" key_value:"interval,omitempty"`
		Timeout  string `json:"timeout" yaml:"timeout,omitempty" key_value:"timeout,omitempty"`
	}
	var tmp httpHealthCheck
 	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.Scheme = tmp.Scheme
	r.Path = tmp.Path
	r.Port = tmp.Port
	r.Hostname = tmp.Hostname
	r.FollowRedirects = tmp.FollowRedirects
	r.Headers = tmp.Headers
	if r.Interval, err = time.ParseDuration(tmp.Interval); err != nil {
		fmt.Println("*********Interval*********")
		return
	}
	r.Timeout, err = time.ParseDuration(tmp.Timeout)
	if err != nil {
		fmt.Println("*********Interval*********")

		return
	}
	return
}
func (r *HttpHealthCheck) mergeHeaders(src, new map[string]string) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	l := ls + len(new)
	r.Headers = make(map[string]string, l)
	var m = make(map[string]string, l)
	for k, item := range src {
		m[item] = k
	}
	var ind int
	for _, item := range new {
		if k, ok := m[item]; !ok {
			r.Headers[strconv.Itoa(ls+ind)] = item
			ind++
			continue
		} else {
			r.Headers[k] = item
		}
	}
}

type HttpLoadBalancerServer struct {
	Url string `json:"url" yaml:"url,omitempty" key_value:"url,omitempty"`
}

func (r *HttpLoadBalancerServer) ToString() (res string) {
	bt, _ := json.Marshal(r)
	res = string(bt)
	return
}

type HttpTls struct {
	Options      string                       `json:"options" yaml:"options,omitempty" key_value:"options,omitempty"`
	CertResolver string                       `json:"certResolver" yaml:"certResolver,omitempty" key_value:"certResolver,omitempty"`
	Domains      map[string]HttpDomainTlsItem `json:"domains" yaml:"domains,omitempty" key_value:"domains,omitempty"`
}

func (r *HttpTls) mergeDomains(src, new map[string]HttpDomainTlsItem) {
	if len(new) == 0 {
		return
	}
	ls := len(src)
	l := ls + len(new)
	r.Domains = make(map[string]HttpDomainTlsItem, l)
	var m = make(map[string]string, l)
	for k, item := range src {
		m[item.ToString()] = k
	}
	var ind int
	for _, item := range new {
		if k, ok := m[item.ToString()]; !ok {
			r.Domains[strconv.Itoa(ls+ind)] = item
			ind++
			continue
		} else {
			r.Domains[k] = item
		}
	}
}

type HttpDomainTlsItem struct {
	Main        string   `json:"main" yaml:"main,omitempty" key_value:"main,omitempty"`
	Sans        []string `json:"sans" yaml:"sans,omitempty" key_value:"sans,omitempty"`
	DomainValue string   `json:"value" yaml:"value,omitempty" key_value:"value,omitempty"`
}

func (r *HttpDomainTlsItem) ToString() (res string) {
	bt, _ := json.Marshal(r)
	res = string(bt)
	return
}
