// Package base
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
// @Desc 分库分表
package base

import (
	"gorm.io/gorm"
)

type (
	FetchDataParameter struct {
		CommonDb
		Id string `json:"id"`
	}
	FetchDataParameterBatch struct {
		CommonDb
		Ids []string `json:"ids"`
	}
	FetchDataParameterTimesBatch struct {
		CommonDb
		Times []TimeNormal `json:"times"`
	}
	FetchDataHandlerBatch     func(fetchDataParameter *FetchDataParameterBatch) (err error)
	FetchDataTimeHandlerBatch func(fetchDataParameter *FetchDataParameterTimesBatch) (err error)
	FetchDataHandler          func(fetchDataParameter *FetchDataParameter) (err error)
	OperateEveryDatabase      struct {
		Db     *gorm.DB `json:"-"`
		DbName string   `json:"db_name"`
		Tables []string `json:"tables"`
	}

	OperateEveryDatabaseHandler func(oed *OperateEveryDatabase) (err error)

	SubTreasury interface {
		// OperateEveryDatabase 并行处理每个数据库连接的数据（非线程安全，OperateEveryDatabaseHandler方法操作相同的数据时要加锁）
		OperateEveryDatabase(handler OperateEveryDatabaseHandler) (err error)

		// GetHashStringDb 根据字段值（当字段为字符串时使用）获取数据所在的数据
		GetHashStringDb(columnValue string) (db *gorm.DB, dbName string, err error)

		// GetHashIntegerDb 根据字段值（当字段为数字时使用）获取数据所在的数据
		GetHashIntegerDb(columnValue int64) (db *gorm.DB, dbName string, err error)

		// GetHashTable 获取数据所在的表表名
		GetHashTable(columnValue int64) (tableName string, err error)

		// GetDataByIntegerIds 根据数据获取
		GetDataByIntegerIds(ids []int64, fetchDataHandler FetchDataHandlerBatch, mapNumString ...map[int64]string) (err error)

		// GetDataByStringIds 根据数据获取
		GetDataByStringIds(ids []string, fetchDataHandler FetchDataHandlerBatch) (err error)

		GetDataByStringId(id string, fetchDataHandler FetchDataHandler) (err error)

		GetDataByIntegerId(id int64, fetchDataHandler FetchDataHandler) (err error)

		GetHashDbAndTableById(id int64) (db *gorm.DB, dbName, tableName string, err error)

		GetHashDbAndTableByTimeId(timeNormal TimeNormal) (db *gorm.DB, dbName, tableName string, err error)

		GetHashDbAndTableByStringId(id string) (db *gorm.DB, dbName, tableName string, err error)

		// GetASCII 根据字符串获取唯一的ASCII码
		GetASCII(str string) (code int64)

		GetHashNumber(columnValue int64) (dbNumber, tableNumber int64)

		GetDbByDbName(dbNameString string) (db *gorm.DB, dbName string, err error)

		TableNameString(tableIndex int64) (tableName string)
	}
	DataParameterOptions func(arg *ActErrorHandlerResult)
)

func DataParameterActErrorHandlerResult(modelBase ModelBase) (res DataParameterOptions) {
	return func(arg *ActErrorHandlerResult) {
		arg.Model = modelBase
	}
}

func (r *FetchDataParameter) ParseActErrorHandlerResult(options ...DataParameterOptions) (res *ActErrorHandlerResult) {
	res = &ActErrorHandlerResult{}
	res.CommonDb = r.CommonDb
	for _, item := range options {
		item(res)
	}
	return
}

func (r *FetchDataParameterBatch) ParseActErrorHandlerResult(options ...DataParameterOptions) (res *ActErrorHandlerResult) {
	res = &ActErrorHandlerResult{}
	res.CommonDb = r.CommonDb
	for _, item := range options {
		item(res)
	}
	return
}
func (r *FetchDataParameterTimesBatch) ParseActErrorHandlerResult() (res *ActErrorHandlerResult) {
	res = &ActErrorHandlerResult{}
	res.CommonDb = r.CommonDb
	return
}

func NewFetchDataParameterTimesBatch(dbName, tableName string, db *gorm.DB, l int) FetchDataParameterTimesBatch {
	return FetchDataParameterTimesBatch{
		CommonDb: CommonDb{
			DbName:    dbName,
			TableName: tableName,
			Db:        db,
		},
		Times: make([]TimeNormal, 0, l),
	}
}
