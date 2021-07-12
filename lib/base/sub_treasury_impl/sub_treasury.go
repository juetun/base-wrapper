// Package sub_treasury_impl
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package sub_treasury_impl

import (
	"fmt"
	"sync"

	. "github.com/juetun/app-api-user/pkg/lib/base"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"gorm.io/gorm"
)

// FetchDataHandler 调用数据库查询数据操作

type SubTreasuryBaseOption func(p *SubTreasuryBase)

// NewSubTreasuryBase 初始化数据模型
func NewSubTreasuryBase(options ...SubTreasuryBaseOption) (res SubTreasury) {
	p := &SubTreasuryBase{
		Dbs: app_obj.DbMysql,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

type SubTreasuryBase struct {
	// 数据库前缀
	DbPrefix string `json:"db_prefix"`

	// 表统一前缀
	TablePrefix string `json:"table_prefix"`

	// 总数据库数
	DbNumber int64 `json:"db_number"`

	// 总计表数
	TableNumber int64 `json:"table_number"`
	// 当前配置的数据库连接
	Dbs     map[string]*gorm.DB `json:"-"`
	Context *base.Context       `json:"-"`
}

func (r *SubTreasuryBase) OperateEveryDatabase(handler OperateEveryDatabaseHandler) (err error) {

	var i int64
	var sync sync.WaitGroup
	sync.Add(int(r.DbNumber))
	for ; i < r.DbNumber; i++ {
		go func(ind int64) { // 并行更新每个数据库，串行更新数据库的每张表
			defer sync.Done()
			_ = r.doEveryDb(handler, ind)
		}(i)
	}
	sync.Wait()
	return
}
func (r *SubTreasuryBase) doEveryDb(handler OperateEveryDatabaseHandler, i int64) (err error) {
	var (
		operateEveryDatabase OperateEveryDatabase
		db                   *gorm.DB
		dbName               string
		j                    int64
	)
	db, dbName, err = r.getDbByIndex(i)
	operateEveryDatabase = OperateEveryDatabase{
		DbName: dbName,
		Db:     db,
		Tables: make([]string, 0, r.TableNumber),
	}
	for ; j < r.TableNumber; j++ {
		operateEveryDatabase.Tables = append(operateEveryDatabase.Tables, r.tableNameString(j))
	}
	if err = handler(&operateEveryDatabase); err != nil {
		r.Context.Error(map[string]interface{}{"err": err.Error()}, "SubTreasuryBaseOperateEveryDatabase")
	}
	return
}
func (r *SubTreasuryBase) GetHashDbAndTableByStringId(id string) (db *gorm.DB, dbName, tableName string, err error) {
	code := r.GetASCII(id)
	return r.GetHashDbAndTableById(code)
}

func (r *SubTreasuryBase) GetDataByStringId(id string, fetchDataHandler FetchDataHandler) (err error) {
	err = r.GetDataByIntegerId(r.GetASCII(id), fetchDataHandler)
	return
}

func (r *SubTreasuryBase) GetHashNumber(columnValue int64) (dbNumber, tableNumber int64) {
	if r.DbNumber == 0 {
		r.DbNumber = 1
	}
	if r.TableNumber == 0 {
		r.TableNumber = 1
	}
	dbNumber = columnValue % r.DbNumber

	div := columnValue / r.DbNumber
	tableNumber = div % r.TableNumber
	return
}

func (r *SubTreasuryBase) GetDataByIntegerId(id int64, fetchDataHandler FetchDataHandler) (err error) {
	var (
		db                *gorm.DB
		tableName, dbName string
	)

	db, dbName, tableName, err = r.GetHashDbAndTableById(id)
	fetchDataParameter := FetchDataParameter{
		SourceDb:  db,
		DbName:    dbName,
		TableName: tableName,
	}
	err = fetchDataHandler(&fetchDataParameter)
	return
}
func (r *SubTreasuryBase) GetHashTable(columnValue int64) (tableName string, err error) {
	_, tableIndex := r.GetHashNumber(columnValue)
	tableName = r.tableNameString(tableIndex)
	return
}

func (r *SubTreasuryBase) tableNameString(tableIndex int64) (tableName string) {
	tableName = fmt.Sprintf("%s%d", r.TablePrefix, tableIndex)
	return
}

type ItemDb struct {
	Db        *gorm.DB
	Ids       []int64
	DbName    string `json:"db_name"`
	TableName string `json:"table_name"`
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
func (r *SubTreasuryBase) GetDataByIntegerIds(ids []int64, fetchDataHandler FetchDataHandler) (err error) {
	l := len(ids)
	if l == 0 {
		return
	}

	var (
		db                *gorm.DB
		dbName, tableName string
		ok                bool
		dta               ItemDb
		mapDb             = make(map[string]ItemDb, l)
	)

	for _, id := range ids {
		if db, dbName, tableName, err = r.GetHashDbAndTableById(id); err != nil {
			return
		}
		uk := dbName + tableName
		if dta, ok = mapDb[uk]; !ok {
			dta = ItemDb{
				DbName:    dbName,
				TableName: tableName,
				Db:        db,
				Ids:       make([]int64, 0, l),
			}
		}
		dta.Ids = append(dta.Ids, id)
		mapDb[uk] = dta
	}

	var syncG sync.WaitGroup

	syncG.Add(len(mapDb))
	for _, itemDb := range mapDb {
		go func(item ItemDb) {
			defer syncG.Done()
			r.getById(item, fetchDataHandler)
		}(itemDb)
	}
	syncG.Wait()

	return
}
func (r *SubTreasuryBase) GetDataByStringIds(ids []string, fetchDataHandler FetchDataHandler) (err error) {
	idInt := make([]int64, 0, len(ids))
	for _, id := range ids {
		idInt = append(idInt, r.GetASCII(id))
	}
	err = r.GetDataByIntegerIds(idInt, fetchDataHandler)
	return
}

func (r *SubTreasuryBase) getById(it ItemDb, fetchDataHandler FetchDataHandler) {

	var err error
	fetchDataParameter := FetchDataParameter{
		SourceDb:  it.Db,
		DbName:    it.DbName,
		TableName: it.TableName,
	}
	if err = fetchDataHandler(&fetchDataParameter); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"id":        it.Ids,
			"dbName":    it.DbName,
			"tableName": it.TableName,
		}, "subTreasuryBaseError")
		return
	}
}

func (r *SubTreasuryBase) getDbByIndex(index int64) (db *gorm.DB, dbName string, err error) {
	dbName = fmt.Sprintf("%s%d", r.DbPrefix, index)
	var ok bool
	if db, ok = r.Dbs[dbName]; !ok {
		err = fmt.Errorf("the database (%s) what you config is not exists", dbName)
		return
	}
	return
}

func (r *SubTreasuryBase) GetHashIntegerDb(columnValue int64) (db *gorm.DB, dbName string, err error) {
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
func SubTreasuryDbNumber(dbNumber int64) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.DbNumber = dbNumber
	}
}
func SubTreasuryTableNumber(tableNumber int64) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.TableNumber = tableNumber
	}
}
func SubTreasuryContext(ctx *base.Context) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.Context = ctx
	}
}
func SubTreasuryDbPrefix(dbPrefix string) SubTreasuryBaseOption {
	return func(p *SubTreasuryBase) {
		p.DbPrefix = dbPrefix
	}
}
