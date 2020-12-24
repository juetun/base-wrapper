/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:13 下午
 */
package srvs

import (
	"github.com/juetun/base-wrapper/web/wrapper"
)

type ServiceDefault interface {
	Index(arg *wrapper.ArgumentDefault) (res *wrapper.ResultDefault, err error)
	TestEs(arg *wrapper.ArgumentDefault) (result interface{}, err error)
	Tmain(arg *wrapper.ArgumentDefault) (result interface{}, err error)
}
