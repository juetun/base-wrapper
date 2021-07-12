package impl

import (
	  "github.com/juetun/base-wrapper/web/cons/admin"
)

type ConAdminIndexImpl struct {
}

func NewConAdminIndexImpl() admin.ConAdminIndex {
	p := &ConAdminIndexImpl{}
	return p
}
