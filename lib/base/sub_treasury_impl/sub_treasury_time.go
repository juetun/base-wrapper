// Package sub_treasury_impl
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package sub_treasury_impl

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	. "github.com/juetun/base-wrapper/lib/base"
	"gorm.io/gorm"
)

// SubTreasuryTime 分库分表实现(按照时间维度)
// 可通过自定义 GetDbFunc  GetTableFunc 数据数据库和表的算法实现
// 如果不传递方法，则默认使用数值(字符串取每个字符的assi码值之和)取余数作为数据库和表的定位
type SubTreasuryTime struct {

	// TablePrefix 表统一前缀
	TablePrefix string `json:"table_prefix"`

	// TableSuffixDateFormat 表后缀时间格式
	TableSuffixDateFormat string `json:"table_suffix_date_format"`

	// DbDateFormat 数据库分库时间格式
	DbDateFormat string `json:"db_date_format"`

	// dbNumber 总数据库数,通过计算DbNameSpaceList的长度获取
	dbNumber int64 `json:"-"`

	// DbNameSpaceList 数据库访问空间名
	DbNameSpaceList []string `json:"db_name_list"`

	// GetDbFuncHandler 获取数据库连接的算法
	GetDbFuncHandler GetDbWithTimeFunc

	// GetTableFuncHandler 获取数据表连接的算法
	GetTableFuncHandler GetTableWithTimeFunc

	// Context 当前配置的数据库连接
	Context *Context `json:"-"`
}

// SubTreasuryTimeOption 调用分布式数据库操作的对象结构体参数
type SubTreasuryTimeOption func(p *SubTreasuryTime)

// GetDbWithTimeFunc 自定义操作哪个数据库算法实现的操作方法(按时间)
type GetDbWithTimeFunc func(subTreasury *SubTreasuryTime, normal TimeNormal) (db *gorm.DB, dbName string, err error)

// GetTableWithTimeFunc 自定义操作哪个table实现的操作方法(按时间)
type GetTableWithTimeFunc func(subTreasury *SubTreasuryTime, normal TimeNormal) (tableName string, err error)

// HandlerGetTables 根据数据库名返回数据对应的表名
type HandlerGetTables func(dbNameSpace string, r *SubTreasuryTime) (res []string)

func (r *SubTreasuryTime) OperateEveryDatabase(handler OperateEveryDatabaseHandler, handlerGetTables HandlerGetTables) (err error) {

	var i int64
	var syncG sync.WaitGroup
	syncG.Add(int(r.dbNumber))
	for ; i < r.dbNumber; i++ {

		// 并行更新每个数据库，串行更新数据库的每张表
		go func(ind int64) {
			defer syncG.Done()
			_ = r.doEveryDb(handler, handlerGetTables, ind)
		}(i)
	}
	syncG.Wait()
	return
}

func (r *SubTreasuryTime) doEveryDb(handler OperateEveryDatabaseHandler, handlerGetTables HandlerGetTables, dbIndex int64) (err error) {
	var (
		operateEveryDatabase OperateEveryDatabase
		db                   *gorm.DB
		dbName               string
	)

	if db, dbName, err = r.getDbByIndex(dbIndex); err != nil {
		return
	}

	operateEveryDatabase = OperateEveryDatabase{DbName: dbName, Db: db, Tables: make([]string, 0)}
	operateEveryDatabase.Tables = handlerGetTables(dbName, r)

	if err = handler(&operateEveryDatabase); err != nil {
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
		}, "SubTreasuryTimeDoEveryDb")
	}
	return
}

func (r *SubTreasuryTime) GetDataByStringId(timeNormal TimeNormal, fetchDataHandler FetchDataHandler) (err error) {

	var (
		db                *gorm.DB
		tableName, dbName string
	)

	db, dbName, tableName, err = r.GetHashDbAndTableById(timeNormal)
	fetchDataParameter := FetchDataParameter{
		SourceDb:  db,
		DbName:    dbName,
		TableName: tableName,
	}
	err = fetchDataHandler(&fetchDataParameter)
	return
}

func (r *SubTreasuryTime) GetHashTable(columnValue TimeNormal) (tableName string, err error) {
	if r.GetTableFuncHandler != nil {
		tableName, err = r.GetTableFuncHandler(r, columnValue)
		return
	}
	_, tableIndex := r.GetHashNumber(columnValue)
	tableName = r.TableNameString(tableIndex)
	return
}

func (r *SubTreasuryTime) GetHashNumber(columnValue TimeNormal) (dbNumber int64, tableNumber string) {

	tmp := columnValue.Format(r.DbDateFormat)
	n, _ := strconv.ParseInt(tmp, 10, 64)
	dbNumber = int64(math.Floor(float64(n / r.dbNumber)))
	tableNumber = columnValue.Format(r.TableSuffixDateFormat)
	return
}

func (r *SubTreasuryTime) TableNameString(tableIndex string) (tableName string) {
	tableName = fmt.Sprintf("%s%s", r.TablePrefix, tableIndex)
	return
}

func (r *SubTreasuryTime) GetHashDbAndTableById(timeNormal TimeNormal) (db *gorm.DB, dbName, tableName string, err error) {

	// 获取数据库的操作源
	if db, dbName, err = r.GetHashIntegerDb(timeNormal); err != nil {
		return
	}

	// 获取数据库表名
	if tableName, err = r.GetHashTable(timeNormal); err != nil {
		return
	}
	return
}

// GetDataByTimes 使用时间维度分库分表
func (r *SubTreasuryTime) GetDataByTimes(times []TimeNormal, fetchDataHandler FetchDataTimeHandlerBatch) (err error) {
	var l int
	if l = len(times); l == 0 {
		return
	}
	var (
		db                *gorm.DB
		dbName, tableName string
		ok                bool
		dta               FetchDataParameterTimesBatch
		mapDb             = make(map[string]FetchDataParameterTimesBatch, l)
	)

	for _, timeItem := range times {
		if db, dbName, tableName, err = r.GetHashDbAndTableByTimeId(timeItem); err != nil {
			return
		}
		uk := dbName + tableName
		if dta, ok = mapDb[uk]; !ok {
			dta = NewFetchDataParameterTimesBatch(dbName, tableName, db, l)
		}
		dta.Times = append(dta.Times, timeItem)
		mapDb[uk] = dta
	}

	var syncG sync.WaitGroup

	syncG.Add(len(mapDb))
	for _, itemDb := range mapDb {
		go func(item FetchDataParameterTimesBatch) {
			defer syncG.Done()
			r.getByTimeIdBatch(item, fetchDataHandler)
		}(itemDb)
	}
	syncG.Wait()
	return
}

func (r *SubTreasuryTime) getByTimeIdBatch(it FetchDataParameterTimesBatch, handler FetchDataTimeHandlerBatch) {
	var err error
	if err = handler(&it); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"id":        it.Times,
			"dbName":    it.DbName,
			"tableName": it.TableName,
		}, "subTreasuryBaseGetByTimeIdBatch")
		return
	}
}

func (r *SubTreasuryTime) getByIdBatch(it FetchDataParameterBatch, fetchDataHandler FetchDataHandlerBatch) {

	var err error
	if err = fetchDataHandler(&it); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"id":        it.Ids,
			"dbName":    it.DbName,
			"tableName": it.TableName,
		}, "subTreasuryBaseGetByIdBatch")
		return
	}
}
func (r *SubTreasuryTime) getById(it FetchDataParameter, fetchDataHandler FetchDataHandler) {

	var err error
	if err = fetchDataHandler(&it); err != nil {
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
			"id":  it,
		}, "subTreasuryBaseGetById")
		return
	}
}

func (r *SubTreasuryTime) getDbByIndex(index int64) (db *gorm.DB, dbName string, err error) {

	dbName = r.DbNameSpaceList[index]
	return r.GetDbByDbName(dbName)

}

func (r *SubTreasuryTime) GetDbByDbName(dbNameString string) (db *gorm.DB, dbName string, err error) {
	dbName = dbNameString
	s := ""
	var ctx = context.TODO()
	if nil != r.Context.GinContext {
		if tp, ok := r.Context.GinContext.Get(app_obj.TraceId); ok {
			s = fmt.Sprintf("%v", tp)
		}
		// ctx = r.Context.GinContext.Request.Context()
	}
	db, dbName, err = GetDbClient(&GetDbClientData{
		Context:     r.Context,
		DbNameSpace: dbName,
		CallBack: func(db *gorm.DB, dbName string) (dba *gorm.DB, err error) {
			dba = db.WithContext(context.WithValue(ctx, app_obj.DbContextValueKey, DbContextValue{
				TraceId: s,
				DbName:  dbName,
			}))
			return
		},
	})
	return
}

func (r *SubTreasuryTime) GetHashIntegerDb(timeNormal TimeNormal) (db *gorm.DB, dbName string, err error) {
	// 如果自定义了获取数据连接的算法
	if r.GetDbFuncHandler != nil {
		return r.GetDbFuncHandler(r, timeNormal)
	}

	dbNumber, _ := r.GetHashNumber(timeNormal)
	db, dbName, err = r.getDbByIndex(dbNumber)
	return
}

func (r *SubTreasuryTime) GetHashDbAndTableByTimeId(timeNormal TimeNormal) (db *gorm.DB, dbName, tableName string, err error) {
	if db, dbName, err = r.GetHashIntegerDb(timeNormal); err != nil {
		return
	}
	tableName, err = r.GetHashTable(timeNormal)
	return
}

func (r *SubTreasuryTime) GetHashTimeDb(columnValue TimeNormal) (db *gorm.DB, dbName string, err error) {
	return r.GetHashIntegerDb(columnValue)
}

func SubTreasuryTimeTablePrefix(tablePrefix string) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.TablePrefix = tablePrefix
	}
}
func SubTreasuryTimeDbDateFormat(dbDateFormat string) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.DbDateFormat = dbDateFormat
	}
}
func SubTreasuryTimeDbNameSpace(dbNameSpace ...string) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.DbNameSpaceList = dbNameSpace
	}
}

// SubTreasuryTimeTableNumber 表后缀格式
func SubTreasuryTimeTableNumber(tableSuffix string) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.TableSuffixDateFormat = tableSuffix
	}
}

func SubTreasuryTimeGetTableFuncHandler(getTableFuncHandler GetTableWithTimeFunc) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.GetTableFuncHandler = getTableFuncHandler
	}
}
func SubTreasuryTimeGetDbFunc(getDbFunc GetDbWithTimeFunc) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.GetDbFuncHandler = getDbFunc
	}
}
func SubTreasuryTimeContext(ctx *Context) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.Context = ctx
	}
}

// NewSubTreasuryTime 初始化数据模型
func NewSubTreasuryTime(options ...SubTreasuryTimeOption) (res *SubTreasuryTime) {

	p := &SubTreasuryTime{} // Dbs: app_obj.DbMysql,

	for _, option := range options {
		option(p)
	}
	p.dbNumber = int64(len(p.DbNameSpaceList))
	if p.dbNumber == 0 {
		p.dbNumber = 1
		p.DbNameSpaceList = []string{app_obj.DefaultDbNameSpace}
	}

	if p.DbDateFormat == "" {
		p.DbDateFormat = "2006"
	}
	if p.TableSuffixDateFormat == "" {
		p.TableSuffixDateFormat = "200601" // 默认按照月分表
	}
	return p
}
