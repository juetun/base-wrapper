// Package sub_treasury_impl
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package sub_treasury_impl

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	. "github.com/juetun/base-wrapper/lib/base"
	"gorm.io/gorm"
)

var (
	//分库分表配置信息
	DiffDbAndTableConfig = make([]ModelDiffDbAndTable, 0, 10)
)

// SubTreasuryBase 分库分表实现
// 可通过自定义 GetDbFunc  GetTableFunc 数据数据库和表的算法实现
// 如果不传递方法，则默认使用数值(字符串取每个字符的assi码值之和)取余数作为数据库和表的定位
type (
	SubTreasuryBase struct {
		// 表统一前缀
		TablePrefix string `json:"table_prefix"`

		// 总计表数
		TableNumber int64 `json:"table_number"`

		// 总数据库数,通过计算DbNameSpaceList的长度获取
		dbNumber int64 `json:"-"`

		TableComment string `json:"table_comment"` //表的注释

		// 数据库访问空间名
		DbNameSpaceList []string `json:"db_name_list"`

		// 获取数据库连接的算法
		GetDbFuncHandler GetDbFunc

		// 获取数据表连接的算法
		GetTableFuncHandler GetTableFunc

		// 当前配置的数据库连接
		Context *base.Context `json:"-"`
	}
	ModelBaseDiffDbAndTable interface {
		ModelBase
		GetCommonOption(context ...*Context) (res []SubTreasuryBaseOption)
		GetDBAndTableNumber() (dbNameSpace []string, tableNum int64)
		GetHashIndex() (shopId int64)
	}
	ModelDiffDbAndTable struct {
		Key          string                  `json:"key"`
		TableComment string                  `json:"table_comment"`
		Struct       ModelBaseDiffDbAndTable `json:"struct"`
	}
	// GetDbFunc 自定义操作哪个数据库算法实现的操作方法
	GetDbFunc func(subTreasury *SubTreasuryBase, columnValue int64) (db *gorm.DB, dbName string, err error)

	// GetTableFunc 自定义操作哪个table实现的操作方法
	GetTableFunc func(subTreasury *SubTreasuryBase, columnValue int64) (tableName string, err error)

	// SubTreasuryBaseOption 调用分布式数据库操作的对象结构体参数
	SubTreasuryBaseOption func(p *SubTreasuryBase)
)

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
func (r *SubTreasuryBase) GetHashDbAndTableByStringId(id string) (commonDb CommonDb, err error) {
	code := r.GetASCII(id)
	commonDb, err = r.GetHashDbAndTableById(code)
	return
}

func (r *SubTreasuryBase) GetDataByStringId(id string, fetchDataHandler FetchDataHandler) (err error) {
	err = r.GetDataByIntegerId(r.GetASCII(id), fetchDataHandler)
	return
}

func (r *SubTreasuryBase) GetDataByIntegerId(id int64, fetchDataHandler FetchDataHandler) (err error) {
	var (
		commonDb CommonDb
	)

	commonDb, err = r.GetHashDbAndTableById(id)
	fetchDataParameter := FetchDataParameter{
		CommonDb: commonDb,
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
	dbNumber = r.abs(dbNumber)
	div := columnValue / r.dbNumber
	tableNumber = r.abs(div % r.TableNumber)
	return
}

func (r *SubTreasuryBase) abs(value int64) (res int64) {
	return int64(math.Abs(float64(value)))
}

func (r *SubTreasuryBase) TableNameString(tableIndex int64) (tableName string) {
	tableName = fmt.Sprintf("%s%d", r.TablePrefix, tableIndex)
	return
}

func (r *SubTreasuryBase) GetHashDbAndTableById(id int64) (commonDb CommonDb, err error) {

	if commonDb.Db, commonDb.DbName, err = r.GetHashIntegerDb(id); err != nil {
		return
	}

	if commonDb.TableName, err = r.GetHashTable(id); err != nil {
		return
	}
	return
}

//id去重
func (r *SubTreasuryBase) uniqueIds(idsArray []int64) (res []int64) {
	res = make([]int64, 0, len(idsArray))
	var mapId = make(map[int64]bool, len(idsArray))
	for _, id := range idsArray {
		if _, ok := mapId[id]; ok {
			continue
		}
		mapId[id] = true
		res = append(res, id)
	}
	return
}

// GetDataByIntegerIds mapNumString 有值表示Id为字符串格式的数据调用
func (r *SubTreasuryBase) GetDataByIntegerIds(idsArray []int64, fetchDataHandler FetchDataHandlerBatch, mapNumString ...map[int64][]string) (err error) {
	//去重ids
	var (
		ids = r.uniqueIds(idsArray)
		l   = len(ids)
	)

	if l == 0 {
		return
	}

	var (
		ok       bool
		dta      FetchDataParameterBatch
		mapDb    = make(map[string]FetchDataParameterBatch, l)
		f        bool
		v        []string
		commonDb CommonDb
	)

	if len(mapNumString) > 0 {
		f = true
	}
	for _, id := range ids {
		if commonDb, err = r.GetHashDbAndTableById(id); err != nil {
			return
		}
		uk := commonDb.DbName + commonDb.TableName
		if dta, ok = mapDb[uk]; !ok {
			dta = FetchDataParameterBatch{
				CommonDb: commonDb,
				Ids:      make([]string, 0, l),
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

func (r *SubTreasuryBase) uniqueStringIds(idsString []string) (res []string) {
	res = make([]string, 0, len(idsString))
	var mapId = make(map[string]bool, len(idsString))
	for _, id := range idsString {
		if _, ok := mapId[id]; ok {
			continue
		}
		mapId[id] = true
		res = append(res, id)
	}
	return
}

func (r *SubTreasuryBase) GetDataByStringIds(idsString []string, fetchDataHandler FetchDataHandlerBatch) (err error) {

	var (
		//去重ID
		ids = r.uniqueStringIds(idsString)
		ok  bool
		l   = len(ids)
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
	dbLen := len(r.DbNameSpaceList)
	if dbLen == 0 {
		err = fmt.Errorf("SubTreasuryBase DbNameSpaceList is nil")
		return
	}
	if index < 0 {
		err = fmt.Errorf("SubTreasuryBase index is error")
		return
	}
	if index >= int64(dbLen) {
		err = fmt.Errorf("SubTreasuryBase index is out")
		return
	}
	dbName = r.DbNameSpaceList[index]
	return r.GetDbByDbName(dbName)

}

func (r *SubTreasuryBase) GetDbByDbName(dbNameString string) (db *gorm.DB, dbName string, err error) {
	dbName = dbNameString
	s := ""
	var ctx = context.TODO()
	if r.Context != nil && nil != r.Context.GinContext {
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

func SubTreasuryTableComment(tableComment string) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.TableComment = tableComment
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
