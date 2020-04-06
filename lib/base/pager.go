package base

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPageSize = 15
	DefaultPageNo   = 1
)

type Pager struct {
	PageNo     int         `json:"page_no"`
	PageSize   int         `json:"page_size"`
	List       interface{} `json:"list"`
	TotalCount int         `json:"total_count"`
}

func NewPager() *Pager {
	return &Pager{
		PageNo:   1,
		PageSize: 15,
	}
}
func (r *Pager) SetPageNo(pageNo int) *Pager {
	r.PageNo = pageNo
	return r
}
func (r *Pager) SetPageSize(pageSize int) *Pager {
	r.PageSize = pageSize
	return r
}
func (r *Pager) SetList(list interface{}) *Pager {
	r.List = list
	return r
}

func (r *Pager) Offset(page string, limit string) (limitInt int, offset int) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limitInt, err = strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}
	return limitInt, (pageInt - 1) * limitInt
}
func (r *Pager) InitPageBy(c *gin.Context, method string) (limit, offset int) {
	queryPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNo))
	queryLimit := c.DefaultQuery("limit", strconv.Itoa(DefaultPageSize))
	limit, offset = r.Offset(queryPage, queryLimit)
	return
}
