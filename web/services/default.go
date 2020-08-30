/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:13 下午
 */
package services

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/daos"
	"github.com/juetun/base-wrapper/web/pojos"
)

type ServiceDefault struct {
	base.ServiceBase
}

func NewServiceDefault(context ...*base.Context) (p *ServiceDefault) {
	p = &ServiceDefault{}
	p.SetContext(context...)
	return
}
func (r *ServiceDefault) Index(arg *pojos.ArgumentDefault) (res *pojos.ResultDefault, err error) {
	res = &pojos.ResultDefault{}
	dao := daos.NewDaoUser(r.Context)
	res.Users, err = dao.GetUser(arg)
	return
}
