/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:18 下午
 */
package wrapper

import (
	"github.com/juetun/base-wrapper/web/models"
)

type (
	ArgumentDefault struct {
		IdKey string `json:"id_key" form:"id_key"`
	}
	ResultDefault struct {
		Users []models.User
	}
)

func (r *ArgumentDefault) SetPathParam(hid string) {

	return
}
