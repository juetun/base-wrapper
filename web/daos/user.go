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
	"github.com/juetun/base-wrapper/web/pojos"
)

type DaoUser interface {
	GetUser(arg *pojos.ArgumentDefault) (res []models.User, err error)
	TestOrm(arg *pojos.ArgumentDefault) (res []models.User, err error)
}
