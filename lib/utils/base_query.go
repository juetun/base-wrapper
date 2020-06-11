package utils

type BaseQuery struct {
	Page     int    `form:"page_no" json:"page_no"`
	PageSize int    `form:"page_size" json:"page_size"`
	Order    string `form:"order" json:"order"`
	Select   string `form:"select" json:"select"`
	IsDel    int    `form:"is_del" json:"is_del"`
}

func (r *BaseQuery) DefaultPage() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}
}
