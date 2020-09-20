/**
* @Author:changjiang
* @Description:
* @File:pager
* @Version: 1.0.0
* @Date 2020/9/20 5:59 下午
 */
package response

import (
	"strconv"
)

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

type Pager struct {
	List       interface{} `json:"list"`
	TotalCount int         `json:"total_count"`
	PageNo     int         `json:"page_no"`
	PageSize   int         `json:"page_size"`
}

// NewPager
func NewPager() *Pager {
	return &Pager{
		TotalCount: 0,
		PageSize:   10,
		PageNo:     1,
		List:       []interface{}{},
	}
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

// 获取数量的方法
type FetchCount func(pagerObject *Pager) (err error)

// 获取数据得方法
type FetchData func(pagerObject *Pager) (err error)

// 获取分页数据方法
// @params fetchCount 获取总条数调用方法
// @params fetchData 获取数据列表调用方法
func (p *Pager) CallGetPagerData(fetchCount FetchCount, fetchData FetchData) {

	// 获取总条数
	fetchCount(p)

	// 如果总条数大于0,获取数据列表
	if p.TotalCount > 0 {
		fetchData(p)
	}
}

//
func (p *Pager) SetPageNoAndSize(pageNo string, pageSize string) error {
	err := p.SetPageNo(pageNo)
	if err != nil {
		return err
	}
	err = p.SetPageSize(pageSize)
	return err
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
