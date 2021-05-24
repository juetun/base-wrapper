package base

// import (
// 	"strconv"
//
// 	"github.com/gin-gonic/gin"
// )
//
// const (
// 	DefaultPageSize = 15
// 	DefaultPageNo   = 1
// )
//
// type ReqPager struct {
// 	PageNo   int `json:"page_no" form:"page_no"`
// 	PageSize int `json:"page_size" form:"page_size"`
// }
//
// func (r *ReqPager) DefaultPager() {
// 	if r.PageNo < 1 {
// 		r.PageNo = DefaultPageNo
// 	}
// 	if r.PageSize == 0 {
// 		r.PageSize = DefaultPageSize
// 	}
// }
// func (r *ReqPager) GetOffset() (offset int) {
// 	r.DefaultPager()
// 	offset = (r.PageNo - 1) * r.PageSize
// 	return
// }
//
// type Pager struct {
// 	ReqPager
// 	List       interface{} `json:"list"`
// 	TotalCount int         `json:"total_count"`
// }
//
// func (r *Pager) DefaultPager() {
// 	if r.PageNo < 1 {
// 		r.PageNo = DefaultPageNo
// 	}
// 	if r.PageSize < 1 {
// 		r.PageSize = DefaultPageSize
// 	}
// }
// func NewPager() *Pager {
// 	return &Pager{
// 		ReqPager: ReqPager{
// 			PageNo:   DefaultPageNo,
// 			PageSize: DefaultPageSize,
// 		},
// 	}
// }
// func (r *Pager) SetPageNo(pageNo int) *Pager {
// 	if pageNo == 0 {
// 		r.PageNo = DefaultPageNo
// 		return r
// 	}
// 	r.PageNo = pageNo
// 	return r
// }
// func (r *Pager) SetPageSize(pageSize int) *Pager {
// 	if pageSize == 0 {
// 		r.PageSize = DefaultPageSize
// 		return r
// 	}
// 	r.PageSize = pageSize
// 	return r
// }
// func (r *Pager) SetList(list interface{}) *Pager {
// 	r.List = list
// 	return r
// }
//
// //计算偏移量
// func (r *Pager) Offset(page string, limit string) (limitInt int, offset int) {
// 	var err error
// 	if r.PageNo, err = strconv.Atoi(page); err != nil {
// 		r.PageNo = DefaultPageNo
// 	}
// 	if r.PageSize, err = strconv.Atoi(limit); err != nil {
// 		r.PageSize = DefaultPageSize
// 	}
// 	return r.PageSize, (r.PageNo - 1) * r.PageSize
// }
//
// //计算偏移量
// func (r *Pager) GetOffset() (offset int) {
// 	if r.PageNo < 1 {
// 		r.PageNo = DefaultPageNo
// 	}
// 	offset = (r.PageNo - 1) * r.PageSize
// 	return
// }
// func (r *Pager) InitPageBy(c *gin.Context, method string) (limit, offset int) {
// 	queryPage := c.DefaultQuery("page", strconv.Itoa(DefaultPageNo))
// 	queryLimit := c.DefaultQuery("limit", strconv.Itoa(DefaultPageSize))
// 	limit, offset = r.Offset(queryPage, queryLimit)
// 	return
// }
