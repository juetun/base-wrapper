/**
* @Author:changjiang
* @Description:
* @File:user
* @Version: 1.0.0
* @Date 2020/8/18 6:52 下午
 */
package daos

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/models"
	"github.com/juetun/base-wrapper/web/pojos"
)

type DaoUser struct {
	base.ServiceDao
}

func NewDaoUser(context ...*base.Context) (p *DaoUser) {
	p = &DaoUser{}
	p.SetContext(context...)
	return
}

func (r *DaoUser) GetUser(arg *pojos.ArgumentDefault) (res []models.User, err error) {
	err = r.Context.Db.
		Where("key=?", 1).
		Find(&res).
		Error
	return
}
func (r *DaoUser) TestOrm(arg *pojos.ArgumentDefault) (res []models.User, err error) {
	err = r.Context.Db.
		Where("key=?", 1).
		Find(&res).
		Error
	return
}
