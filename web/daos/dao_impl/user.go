/**
* @Author:changjiang
* @Description:
* @File:user
* @Version: 1.0.0
* @Date 2020/8/18 6:52 下午
 */

// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package dao_impl

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/daos"
	"github.com/juetun/base-wrapper/web/models"
	"github.com/juetun/base-wrapper/web/wrapper"
)

type DaoUserImpl struct {
	base.ServiceDao
}

func NewDaoUserImpl(ctx ...*base.Context) (res daos.DaoUser) {
	p := &DaoUserImpl{}
	p.SetContext(ctx...)
	return p
}

func (r *DaoUserImpl) BatchData() (err error) {
	data := make([]base.ModelBase, 0, 1)
	dt := &models.DataChildren{}
	dt.UserId = "1"
	dt.UserHid = "user_id"
	data = append(data, dt)
	param := &base.BatchAddDataParameter{}
	param.Data = data
	err = r.BatchAdd(param)
	return
}

func (r *DaoUserImpl) GetUser(arg *wrapper.ArgumentDefault) (res []models.User, err error) {
	err = r.Context.Db.
		Where("id=?", 1).
		Find(&res).
		Error
	return
}
func (r *DaoUserImpl) TestOrm(arg *wrapper.ArgumentDefault) (res []models.User, err error) {
	err = r.Context.Db.
		Where("key=?", 1).
		Find(&res).
		Error
	return
}
