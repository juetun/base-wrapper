// Package sub_treasury_impl
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package sub_treasury_impl

import (
	"context"
	"fmt"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	. "github.com/juetun/base-wrapper/lib/base"
	"gorm.io/gorm"
)

// SubTreasuryBase 分库分表实现
// 可通过自定义 GetDbFunc  GetTableFunc 数据数据库和表的算法实现
// 如果不传递方法，则默认使用数值(字符串取每个字符的assi码值之和)取余数作为数据库和表的定位
type SubTreasuryBase struct {

	// 表统一前缀
	TablePrefix string `json:"table_prefix"`

	// 总计表数
	TableNumber int64 `json:"table_number"`

	// 总数据库数,通过计算DbNameSpaceList的长度获取
	dbNumber int64 `json:"-"`

	// 数据库访问空间名
	DbNameSpaceList []string `json:"db_name_list"`

	// 获取数据库连接的算法
	GetDbFuncHandler GetDbFunc

	// 获取数据表连接的算法
	GetTableFuncHandler GetTableFunc

	// 当前配置的数据库连接
	Context *base.Context `json:"-"`
}

// GetDbFunc 自定义操作哪个数据库算法实现的操作方法
type GetDbFunc func(subTreasury *SubTreasuryBase, columnValue int64) (db *gorm.DB, dbName string, err error)

// GetTableFunc 自定义操作哪个table实现的操作方法
type GetTableFunc func(subTreasury *SubTreasuryBase, columnValue int64) (tableName string, err error)

// SubTreasuryBaseOption 调用分布式数据库操作的对象结构体参数
type SubTreasuryBaseOption func(p *SubTreasuryBase)

func (r *SubTreasuryBase) OperateEveryDatabase(handler OperateEveryDatabaseHandler) (err error) {

	var i int64
	var syncG sync.WaitGroup
	syncG.Add(int(r.dbNumber))
	for ; i < r.dbNumber; i++ {

		// 并行更新每个数据库，串行更新数据库的每张表
		go func(ind int64) {
			defer syncG.Done()
			_ = r.doEveryDb(handler, ind)
		}(i)
	}
	syncG.Wait()
	return
}
func (r *SubTreasuryBase) doEveryDb(handler OperateEveryDatabaseHandler, i int64) (err error) {
	var (
		operateEveryDatabase OperateEveryDatabase
		db                   *gorm.DB
		dbName               string
		j                    int64
	)

	if db, dbName, err = r.getDbByIndex(i); err != nil {
		return
	}

	operateEveryDatabase = OperateEveryDatabase{DbName: dbName, Db: db, Tables: make([]string, 0, r.TableNumber)}

	for ; j < r.TableNumber; j++ {
		operateEveryDatabase.Tables = append(operateEveryDatabase.Tables, r.TableNameString(j))
	}

	if err = handler(&operateEveryDatabase); err != nil {
		r.Context.Error(map[string]interface{}{"err": err.Error()}, "SubTreasuryBaseOperateEveryDatabase")
	}
	return
}
func (r *SubTreasuryBase) GetHashDbAndTableByStringId(id string) (db *gorm.DB, dbName, tableName string, err error) {
	code := r.GetASCII(id)
	db, dbName, tableName, err = r.GetHashDbAndTableById(code)
	return
}

func (r *SubTreasuryBase) GetDataByStringId(id string, fetchDataHandler FetchDataHandler) (err error) {
	err = r.GetDataByIntegerId(r.GetASCII(id), fetchDataHandler)
	return
}

func (r *SubTreasuryBase) GetDataByIntegerId(id int64, fetchDataHandler FetchDataHandler) (err error) {
	var (
		db                *gorm.DB
		tableName, dbName string
	)

	db, dbName, tableName, err = r.GetHashDbAndTableById(id)
	fetchDataParameter := FetchDataParameter{
		CommonDb:CommonDb{
			Db:  db,
			DbName:    dbName,
			TableName: tableName,
		},
	}
	err = fetchDataHandler(&fetchDataParameter)
	return
}
func (r *SubTreasuryBase) GetHashTable(columnValue int64) (tableName string, err error) {
	if r.GetTableFuncHandler != nil {
		tableName, err = r.GetTableFuncHandler(r, columnValue)
		return
	}
	_, tableIndex := r.GetHashNumber(columnValue)
	tableName = r.TableNameString(tableIndex)
	return
}
func (r *SubTreasuryBase) GetHashNumber(columnValue int64) (dbNumber, tableNumber int64) {

	dbNumber = columnValue % r.dbNumber

	div := columnValue / r.dbNumber
	tableNumber = div % r.TableNumber
	return
}

func (r *SubTreasuryBase) TableNameString(tableIndex int64) (tableName string) {
	tableName = fmt.Sprintf("%s%d", r.TablePrefix, tableIndex)
	return
}

func (r *SubTreasuryBase) GetHashDbAndTableById(id int64) (db *gorm.DB, dbName, tableName string, err error) {

	if db, dbName, err = r.GetHashIntegerDb(id); err != nil {
		return
	}

	if tableName, err = r.GetHashTable(id); err != nil {
		return
	}
	return
}

// GetDataByIntegerIds mapNumString 有值表示Id为字符串格式的数据调用
func (r *SubTreasuryBase) GetDataByIntegerIds(ids []int64, fetchDataHandler FetchDataHandlerBatch, mapNumString ...map[int64][]string) (err error) {
	l := len(ids)
	if l == 0 {
		return
	}

	var (
		db                *gorm.DB
		dbName, tableName string
		ok                bool
		dta               FetchDataParameterBatch
		mapDb             = make(map[string]FetchDataParameterBatch, l)
		f                 bool
		v                 []string
	)

	if len(mapNumString) > 0 {
		f = true
	}
	for _, id := range ids {
		if db, dbName, tableName, err = r.GetHashDbAndTableById(id); err != nil {
			return
		}
		uk := dbName + tableName
		if dta, ok = mapDb[uk]; !ok {
			dta = FetchDataParameterBatch{
				CommonDb:CommonDb{
					Db:  db,
					DbName:    dbName,
					TableName: tableName,
				},
				Ids:       make([]string, 0, l),
			}
		}

		if f {
			v = mapNumString[0][id]
		} else {
			v = []string{fmt.Sprintf("%d", id)}
		}
		dta.Ids = append(dta.Ids, v...)
		mapDb[uk] = dta
	}

	var syncG sync.WaitGroup

	syncG.Add(len(mapDb))
	for _, itemDb := range mapDb {
		go func(item FetchDataParameterBatch) {
			defer syncG.Done()
			r.getByIdBatch(item, fetchDataHandler)
		}(itemDb)
	}
	syncG.Wait()

	return
}
func (r *SubTreasuryBase) GetDataByStringIds(ids []string, fetchDataHandler FetchDataHandlerBatch) (err error) {
	var (
		ok bool
		l  = len(ids)
	)
	idInt := make([]int64, 0, l)
	var mapId = make(map[int64][]string, l)
	for _, id := range ids {
		idNum := r.GetASCII(id)
		idInt = append(idInt, idNum)
		if _, ok = mapId[idNum]; !ok {
			mapId[idNum] = make([]string, 0, l)
		}
		mapId[idNum] = append(mapId[idNum], id)
	}
	err = r.GetDataByIntegerIds(idInt, fetchDataHandler, mapId)
	return
}

func (r *SubTreasuryBase) getByIdBatch(it FetchDataParameterBatch, fetchDataHandler FetchDataHandlerBatch) {

	var err error
	if err = fetchDataHandler(&it); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"id":        it.Ids,
			"dbName":    it.DbName,
			"tableName": it.TableName,
		}, "subTreasuryBaseError")
		return
	}
}
func (r *SubTreasuryBase) getById(it FetchDataParameter, fetchDataHandler FetchDataHandler) {

	var err error
	if err = fetchDataHandler(&it); err != nil {
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
			"id":  it,
		}, "subTreasuryBaseError")
		return
	}
}

func (r *SubTreasuryBase) getDbByIndex(index int64) (db *gorm.DB, dbName string, err error) {

	dbName = r.DbNameSpaceList[index]
	return r.GetDbByDbName(dbName)

}

func (r *SubTreasuryBase) GetDbByDbName(dbNameString string) (db *gorm.DB, dbName string, err error) {
	dbName = dbNameString
	s := ""
	var ctx = context.TODO()
	if nil != r.Context.GinContext {
		if tp, ok := r.Context.GinContext.Get(app_obj.TraceId); ok {
			s = fmt.Sprintf("%v", tp)
		}
		// ctx = r.Context.GinContext.Request.Context()
	}
	db, dbName, err = base.GetDbClient(&base.GetDbClientData{
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

func (r *SubTreasuryBase) GetHashIntegerDb(columnValue int64) (db *gorm.DB, dbName string, err error) {
	// 如果自定义了获取数据连接的算法
	if r.GetDbFuncHandler != nil {
		return r.GetDbFuncHandler(r, columnValue)
	}

	dbNumber, _ := r.GetHashNumber(columnValue)
	db, dbName, err = r.getDbByIndex(dbNumber)
	return
}

// GetASCII 将字符串的ASCII码总和拼接成一个数字，作为唯一的key
func (r *SubTreasuryBase) GetASCII(str string) (code int64) {
	runeS := []rune(str)
	for _, r2 := range runeS {
		code += int64(r2)
	}
	return
}
func (r *SubTreasuryBase) GetHashStringDb(columnValue string) (db *gorm.DB, dbName string, err error) {
	return r.GetHashIntegerDb(r.GetASCII(columnValue))
}

func SubTreasuryTablePrefix(tablePrefix string) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.TablePrefix = tablePrefix
	}
}

func SubTreasuryDbNameSpace(dbNameSpace ...string) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.DbNameSpaceList = dbNameSpace
	}
}
func SubTreasuryTableNumber(tableNumber int64) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.TableNumber = tableNumber
	}
}

func SubTreasuryGetTableFuncHandler(getTableFuncHandler GetTableFunc) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.GetTableFuncHandler = getTableFuncHandler
	}
}
func SubTreasuryGetDbFunc(getDbFunc GetDbFunc) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.GetDbFuncHandler = getDbFunc
	}
}
func SubTreasuryContext(ctx *base.Context) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.Context = ctx
	}
}

// NewSubTreasuryBase 初始化数据模型
func NewSubTreasuryBase(options ...SubTreasuryBaseOption) (res *SubTreasuryBase) {

	p := &SubTreasuryBase{} // Dbs: app_obj.DbMysql,

	for _, option := range options {
		option(p)
	}
	p.dbNumber = int64(len(p.DbNameSpaceList))
	if p.dbNumber == 0 {
		p.dbNumber = 1
		p.DbNameSpaceList = []string{app_obj.DefaultDbNameSpace}
	}
	if p.TableNumber == 0 {
		p.TableNumber = 1
	}
	return p
}
