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

	//表注释
	TableComment string `json:"table_comment"`

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
	Ctx     context.Context
}

// SubTreasuryTimeOption 调用分布式数据库操作的对象结构体参数
type SubTreasuryTimeOption func(p *SubTreasuryTime)

// GetDbWithTimeFunc 自定义操作哪个数据库算法实现的操作方法(按时间)
type GetDbWithTimeFunc func(subTreasury *SubTreasuryTime, normal int64) (db *gorm.DB, dbName string, err error)

// GetTableWithTimeFunc 自定义操作哪个table实现的操作方法(按时间)
type GetTableWithTimeFunc func(subTreasury *SubTreasuryTime, normal int64) (tableName string, err error)

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
		CommonDb: CommonDb{
			Db:        db,
			DbName:    dbName,
			TableName: tableName,
		},
	}
	err = fetchDataHandler(&fetchDataParameter)
	return
}

func (r *SubTreasuryTime) GetHashTable(tableIndex int64) (tableName string, err error) {
	if r.GetTableFuncHandler != nil {
		if tableName, err = r.GetTableFuncHandler(r, tableIndex); err != nil {
			return
		}
		return
	}
	tableName = r.TableNameString(tableIndex)
	return
}

//func (r *SubTreasuryTime) getHashNumber(columnValue TimeNormal) (dbNumber int64, tableNumber string) {
//	dbNumber = r.getDbName(columnValue)
//	tableNumber = columnValue.Format(r.TableSuffixDateFormat)
//	return
//}

func (r *SubTreasuryTime) getDbName(dbIndex int64) (dbNumber int64) {
	dbNumber = int64(math.Floor(float64(dbIndex / r.dbNumber)))
	return
}
func (r *SubTreasuryTime) TableNameString(tableIndex int64) (tableName string) {
	tableName = fmt.Sprintf("%s%d", r.TablePrefix, tableIndex)
	return
}

func (r *SubTreasuryTime) GetHashDbAndTableById(timeNormal TimeNormal) (db *gorm.DB, dbName, tableName string, err error) {

	// 获取数据库的操作源
	if db, dbName, err = r.GetHashIntegerDb(timeNormal); err != nil {
		return
	}
	var (
		tableIndex       int64
		tableIndexString = timeNormal.Format(r.TableSuffixDateFormat)
	)

	if tableIndex, err = strconv.ParseInt(tableIndexString, 10, 64); err != nil {
		return
	}
	// 获取数据库表名
	if tableName, err = r.GetHashTable(tableIndex); err != nil {
		return
	}
	return
}

func (r *SubTreasuryTime) getStartAndOverTimeNumber(startTime, overTime TimeNormal) (overTimeIndexNumber, startTimeIndexNumber int64, err error) {

	if overTimeIndexNumber, err = strconv.ParseInt(overTime.Format(r.TableSuffixDateFormat), 10, 64); err != nil {
		return
	}
	if startTimeIndexNumber, err = strconv.ParseInt(startTime.Format(r.TableSuffixDateFormat), 10, 64); err != nil {
		return
	}
	return
}

// 按照时间开始时间 和 结束时间跨表查询数据使用
func (r *SubTreasuryTime) GetDataByTimeStartAndOverTime(startTime, overTime TimeNormal, fetchDataHandler OperateEveryDatabaseHandler) (err error) {

	var (
		overTimeIndexNumber, startTimeIndexNumber int64
		mapDb                                     = make(map[string]OperateEveryDatabase, r.dbNumber)
		db                                        *gorm.DB
		ok                                        bool
		dbName, tableName                         string
		dta                                       OperateEveryDatabase
	)

	//时间交换保证开始时间一定小于结束时间
	if startTime.After(overTime.Time) {
		overTime, startTime = startTime, overTime
	}

	if overTimeIndexNumber, startTimeIndexNumber, err = r.getStartAndOverTimeNumber(startTime, overTime); err != nil {
		return
	}

	for {
		if db, dbName, tableName, err = r.GetHashDbAndTableByIndexId(startTimeIndexNumber); err != nil {
			return
		}
		if dta, ok = mapDb[dbName]; !ok {
			dta = OperateEveryDatabase{Db: db, DbName: dbName, Tables: []string{tableName}}
		}
		if startTimeIndexNumber > overTimeIndexNumber {
			break
		}
		dta.Tables = append(dta.Tables, tableName)
		startTimeIndexNumber += 1
		mapDb[dbName] = dta
	}
	r.getDataByTimeStartAndOverTimeItemDb(mapDb, fetchDataHandler)
	return
}

func (r *SubTreasuryTime) getDataByTimeStartAndOverTimeItemDb(mapDb map[string]OperateEveryDatabase, fetchDataHandler OperateEveryDatabaseHandler) {
	var syncG sync.WaitGroup
	syncG.Add(len(mapDb))
	var actHandler = func(k string, item OperateEveryDatabase) {
		defer syncG.Done()
		//执行每个数据库操作动作
		if err := fetchDataHandler(&item); err != nil {
			r.Context.Error(map[string]interface{}{
				"err": err.Error(),
				"k":   k,
			}, "SubTreasuryTimeGetDataByTimeStartAndOverTimeItemDb")
		}
	}
	for k, itemDb := range mapDb {
		go actHandler(k, itemDb)
	}
	syncG.Wait()
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
	if nil != r.Context.GinContext {
		if tp, ok := r.Context.GinContext.Get(app_obj.TraceId); ok {
			s = fmt.Sprintf("%v", tp)
		}
	}
	db, dbName, err = GetDbClient(&GetDbClientData{
		Context:     r.Context,
		DbNameSpace: dbName,
		CallBack: func(db *gorm.DB, dbName string) (dba *gorm.DB, err error) {
			dba = db.WithContext(context.WithValue(r.Ctx, app_obj.DbContextValueKey, DbContextValue{
				TraceId: s,
				DbName:  dbName,
			}))
			return
		},
	})
	return
}

func (r *SubTreasuryTime) getHashIntegerDbWithNumber(dbIndexNum int64) (db *gorm.DB, dbName string, err error) {

	// 如果自定义了获取数据连接的算法
	if r.GetDbFuncHandler != nil {
		return r.GetDbFuncHandler(r, dbIndexNum)
	}
	configDatabaseIndex := r.getDbName(dbIndexNum)
	if db, dbName, err = r.getDbByIndex(configDatabaseIndex); err != nil {
		return
	}
	return
}

func (r *SubTreasuryTime) GetHashIntegerDb(timeNormal TimeNormal) (db *gorm.DB, dbName string, err error) {
	var (
		dbIndexNum int64
		dbIndex    = timeNormal.Format(r.DbDateFormat)
	)

	if dbIndexNum, err = strconv.ParseInt(dbIndex, 10, 64); err != nil {
		return
	}

	db, dbName, err = r.getHashIntegerDbWithNumber(dbIndexNum)
	return
}

func (r *SubTreasuryTime) GetHashDbAndTableByTimeId(timeNormal TimeNormal) (db *gorm.DB, dbName, tableName string, err error) {
	if db, dbName, err = r.GetHashIntegerDb(timeNormal); err != nil {
		return
	}
	var tableIndex int64
	if tableIndex, err = r.getTableIndexWith(timeNormal); err != nil {
		return
	}
	tableName, err = r.GetHashTable(tableIndex)
	return
}
func (r *SubTreasuryTime) getTableIndexWith(timeNormal TimeNormal) (tableIndex int64, err error) {
	timeNormal.Format(r.TableSuffixDateFormat)
	return
}

func (r *SubTreasuryTime) GetHashDbAndTableByIndexId(indexId int64) (db *gorm.DB, dbName, tableName string, err error) {

	if db, dbName, err = r.getHashIntegerDbWithNumber(indexId); err != nil {
		return
	}
	tableName, err = r.GetHashTable(indexId)
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

func SubTreasuryTimeTableComment(tableComment string) SubTreasuryTimeOption {
	return func(p *SubTreasuryTime) {
		p.TableComment = tableComment
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
	if p.Ctx == nil {
		p.Ctx = context.TODO()
	}
	return p
}
