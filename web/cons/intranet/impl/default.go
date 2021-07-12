package impl

import (
	"github.com/juetun/base-wrapper/web/cons/intranet"
)

type ConIntranetIndexImpl struct {
}

func NewConIntranetIndexImpl() intranet.ConIntranetIndex {
	p := &ConIntranetIndexImpl{}
	return p
}
