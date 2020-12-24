/**
* @Author:changjiang
* @Description:
* @File:user
* @Version: 1.0.0
* @Date 2020/8/18 6:52 下午
 */
package daos

import (
	"github.com/juetun/base-wrapper/web/models"
	"github.com/juetun/base-wrapper/web/wrapper"
)

type DaoUser interface {
	GetUser(arg *wrapper.ArgumentDefault) (res []models.User, err error)
	TestOrm(arg *wrapper.ArgumentDefault) (res []models.User, err error)
}
