// Package response
/**
* @Author:ChangJiang
* @Description:
* @File:pager
* @Version: 1.0.0
* @Date 2020/9/20 5:59 下午
 */
package response

import (
	"strconv"
)

const (
	DefaultPageSize = 15
	DefaultPageNo   = 1
)

type PagerParameter struct {
	PageNo    int    `form:"page_no" json:"page_no,omitempty"`
	PageSize  int    `form:"page_size" json:"page_size,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}
type PageQuery struct {
	Order  string `form:"order" json:"order,omitempty"`
	Select string `form:"select" json:"select,omitempty"`
	IsDel  int    `form:"is_del" json:"is_del,omitempty"`
	PagerParameter
}

func (r *PageQuery) GetOffset() (offset int) {
	if r.PageNo < 1 {
		r.PageNo = DefaultPageNo
	}
	offset = (r.PageNo - 1) * r.PageSize
	return
}
func (r *PageQuery) DefaultPage() {
	if r.PageNo < 1 {
		r.PageNo = 1
	}
	if r.PageSize == 0 {
		r.PageSize = DefaultPageSize
	}
}

type Pager struct {
	List       interface{} `json:"list"`
	TotalCount int         `json:"total_count,omitempty"`
	IsNext     bool        `json:"is_next,omitempty"` // [bool] 是否有下一页，true=有下一页；false=无下页，可关闭列表
	PagerParameter
}
type PageHandler func(*Pager)
type PageOption PageHandler

func PagerList(list interface{}) PageOption {
	return func(pager *Pager) {
		pager.List = list
	}
}
func PagerBaseQuery(baseQuery PageQuery) PageOption {
	if baseQuery.PageSize == 0 {
		baseQuery.PageSize = DefaultPageSize
	}
	return func(pager *Pager) {
		pager.PagerParameter = baseQuery.PagerParameter
	}
}
func PagerTotalCount(totalCount int) PageOption {
	return func(pager *Pager) {
		pager.TotalCount = totalCount
	}
}
func PagerPageNo(pageNo int) PageOption {
	return func(pager *Pager) {
		pager.PageNo = pageNo
	}
}
func PagerPageSize(pageSize int) PageOption {
	return func(pager *Pager) {
		pager.PageSize = pageSize
	}
}

// NewPager 初始化分页对象
func NewPager(option ...PageOption) *Pager {
	r := &Pager{
		TotalCount: 0,
		PagerParameter: PagerParameter{
			PageNo:   1,
			PageSize: DefaultPageSize,
		},
		List: []interface{}{},
	}
	for _, item := range option {
		item(r)
	}
	return r
}
func NewPagerAndDefault(arg *PageQuery) (pager *Pager) {
	pager = NewPager()
	pager.InitPager(arg)
	return
}

func (p *Pager) InitPager(arg *PageQuery) *Pager {
	if arg.PageNo == 0 {
		arg.PageNo = 1
	}
	p.PageNo = arg.PageNo
	if arg.PageSize == 0 {
		arg.PageSize = DefaultPageSize
	}
	p.PageSize = arg.PageSize
	return p
}

// InitPageNoAndPageSize 初始化PageNo 和PageSize
func (p *Pager) InitPageNoAndPageSize(params *map[string]string) error {
	var err error
	var pageNo, pageSize string
	if _, ok := (*params)["page_no"]; ok {
		pageNo = (*params)["page_no"]
	}
	if pageNo == "" {
		pageNo = "0"
	}
	err = p.SetPageNo(pageNo)
	if err != nil {
		return err
	}
	if _, ok := (*params)["page_size"]; ok {
		pageSize = (*params)["page_size"]
	}
	if pageSize == "" {
		pageSize = "0"
	}
	err = p.SetPageSize(pageSize)
	if err != nil {
		return err
	}
	return err
}

// FetchCount 获取数量的方法
type FetchCount func(pagerObject *Pager) (err error)

// FetchData 获取数据得方法
type FetchData func(pagerObject *Pager) (err error)

// CallGetPagerData 获取分页数据方法
// @params fetchCount 获取总条数调用方法
// @params fetchData 获取数据列表调用方法
func (p *Pager) CallGetPagerData(fetchCount FetchCount, fetchData FetchData) (err error) {

	// 获取总条数
	if err = fetchCount(p); err != nil {
		return
	}

	// 如果总条数大于0,获取数据列表
	if p.TotalCount > 0 {
		err = fetchData(p)
	}
	return
}

// SetPageNoAndSize 初始化分页
func (p *Pager) SetPageNoAndSize(pageNo string, pageSize string) (err error) {

	if err = p.SetPageNo(pageNo); err != nil {
		return err
	}
	err = p.SetPageSize(pageSize)
	return
}

func (p *Pager) GetFromAndLimit() int {
	return (p.PageNo - 1) * p.PageSize
}

func (p *Pager) SetPageNo(pageNo string) error {
	pageNumber, err := strconv.Atoi(pageNo)
	if err != nil {
		return err
	}
	p.PageNo = pageNumber
	if p.PageNo < 1 {
		p.PageNo = 1
	}
	return nil
}

func (p *Pager) SetPageSize(pageSize string) error {
	pageSizeNumber, err := strconv.Atoi(pageSize)
	if err != nil {
		return err
	}
	p.PageSize = pageSizeNumber
	if p.PageSize < 1 {
		p.PageSize = 15
	}
	return nil
}
