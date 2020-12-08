/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:13 下午
 */
package srvs

import (
	"github.com/juetun/base-wrapper/web/pojos"
)

type ServiceDefault interface {
	Index(arg *pojos.ArgumentDefault) (res *pojos.ResultDefault, err error)
	TestEs(arg *pojos.ArgumentDefault) (result interface{}, err error)
	Tmain(arg *pojos.ArgumentDefault) (result interface{}, err error)
}
