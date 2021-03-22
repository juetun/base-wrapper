// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

import (
	"fmt"
	"strconv"
)

type TraefikService struct {
	RouterName  string          `json:"router_name"`
	Rule        string          `json:"rule"`
	Service     string          `json:"service"`
	EntryPoints []string        `json:"entry_points"`
	Middlewares []string        `json:"middlewares"`
	ServiceList []ServiceConfig `json:"service_list"`
	Priority    int             `json:"priority"`
	Tls         *Tls            `json:"tls"`
}
type ServiceConfig struct {
	ServiceName  string       `json:"service_name"`
	LoadBalancer LoadBalancer `json:"load_balancer"`
	Mirroring    Mirroring    `json:"mirroring"`
	Weighted     Weighted     `json:"weighted"`
}
type Mirroring struct {
	Service string    `json:"service"`
	Mirrors []Mirrors `json:"mirrors"`
}
type Mirrors struct {
	Name    string `json:"name"`
	Percent int    `json:"percent"`
}
type Weighted struct {
	Services []WeightedService `json:"services"`
	Sticky   *Sticky           `json:"sticky"`
}
type WeightedService struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}
type LoadBalancer struct {
	Servers            []LoadBalancerServer
	ServersTransport   string            `json:"serverstransport"`
	HealthCheck        HealthCheck       `json:"health_check"`
	Sticky             *Sticky           `json:"sticky"`
	ResponseForwarding map[string]string `json:"responseforwarding"`
}
type Sticky struct {
	Value  bool   `json:"value"`
	Cookie Cookie `json:"cookie"`
}
type Cookie struct {
	HttpOnly bool   `json:"httponly"`
	Name     string `json:"name"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"samesite"`
}
type HealthCheck struct {
	Headers  map[string]string `json:"headers"`
	Hostname string            `json:"hostname"`
	Interval int               `json:"interval"`
	Path     string            `json:"path"`
	Port     int               `json:"port"`
	Scheme   string            `json:"scheme"`
	Timeout  int               `json:"timeout"`
}
type LoadBalancerServer struct {
	Url string `json:"url"`
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


func (r *TraefikConfig) serviceList() {

}
func (r *TraefikConfig) tls() {
	for _, item := range r.RouterHttpConfig {
		if item.Tls == nil {
			continue
		}
		r.MapValue = append(r.MapValue, KeyValue{
			Key: fmt.Sprintf("traefik/http/routers/%s/tls",
				item.RouterName),
			Value: fmt.Sprintf("%v", item.Tls.Value),
		})

		if item.Tls.CertResolver != "" {

			r.MapValue = append(r.MapValue, KeyValue{
				Key:   fmt.Sprintf("traefik/http/routers/%s/tls/certresolver", item.RouterName),
				Value: item.Tls.CertResolver,
			})

		}
		if len(item.Tls.Domains) > 0 {
			for k, domain := range item.Tls.Domains {
				r.MapValue = append(r.MapValue, KeyValue{
					Key: fmt.Sprintf("traefik/http/routers/%s/tls/domains/%d/%s",
						item.RouterName, k, domain.Domain),
					Value: domain.DomainValue,
				})
				for i, v := range domain.Sans {
					r.MapValue = append(r.MapValue, KeyValue{
						Key: fmt.Sprintf("traefik/http/routers/%s/tls/domains/%d/sans/%d",
							item.RouterName, k, i),
						Value: v,
					})

				}
			}
		}
		if item.Tls.Options != "" {
			r.MapValue = append(r.MapValue, KeyValue{
				Key: fmt.Sprintf("traefik/http/routers/%s/tls/options",
					item.RouterName),
				Value: item.Tls.Options,
			})

		}
	}
}

func (r *TraefikConfig) middlewares() {
	for _, item := range r.RouterHttpConfig {
		for k, middleware := range item.Middlewares {
			if middleware == "" {
				continue
			}
			r.MapValue = append(r.MapValue, KeyValue{
				Key: fmt.Sprintf("traefik/http/routers/%s/middlewares/%d",
					item.RouterName, k),
				Value: middleware,
			})

		}
	}
}

func (r *TraefikConfig) router() {
	for _, item := range r.RouterHttpConfig {
		if item.RouterName == "" {
			continue
		}
		r.MapValue = append(r.MapValue, KeyValue{
			Key:   fmt.Sprintf("traefik/http/routers/%s/rule", item.RouterName),
			Value: item.Rule,
		})
	}
	return
}

func (r *TraefikConfig) entryPoints() {
	for _, item := range r.RouterHttpConfig {
		for k, it := range item.EntryPoints {
			if it == "" {
				continue
			}
			r.MapValue = append(r.MapValue, KeyValue{
				Key: fmt.Sprintf("traefik/http/routers/%s/entrypoints/%d",
					item.RouterName, k),
				Value: it,
			})

		}
	}
}

func (r *TraefikConfig) service() {
	for _, item := range r.RouterHttpConfig {
		if item.Service == "" {
			continue
		}
		r.MapValue = append(r.MapValue, KeyValue{
			Key: fmt.Sprintf("traefik/http/routers/%s/service",
				item.RouterName),
			Value: item.Service,
		})

	}
}

func (r *TraefikConfig) priority() {
	for _, item := range r.RouterHttpConfig {
		if item.Priority == 0 {
			continue
		}
		r.MapValue = append(r.MapValue, KeyValue{
			Key: fmt.Sprintf("traefik/http/routers/%s/priority",
				item.RouterName),
			Value: strconv.Itoa(item.Priority),
		})
	}
}
