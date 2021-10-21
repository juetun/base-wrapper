// Package base
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
// @Desc 分库分表
package base

import (
	"gorm.io/gorm"
)

type FetchDataParameter struct {
	SourceDb  *gorm.DB `json:"-"`
	DbName    string   `json:"db_name"`
	TableName string   `json:"table_name"`
	Id        string   `json:"id"`
}
type FetchDataParameterBatch struct {
	SourceDb  *gorm.DB `json:"-"`
	DbName    string   `json:"db_name"`
	TableName string   `json:"table_name"`
	Ids       []string `json:"ids"`
}


type FetchDataParameterTimesBatch struct {
	SourceDb  *gorm.DB     `json:"-"`
	DbName    string       `json:"db_name"`
	TableName string       `json:"table_name"`
	Times     []TimeNormal `json:"times"`
}

func NewFetchDataParameterTimesBatch(dbName, tableName string, db *gorm.DB, l int) FetchDataParameterTimesBatch {
	return FetchDataParameterTimesBatch{
		DbName:    dbName,
		TableName: tableName,
		SourceDb:  db,
		Times:     make([]TimeNormal, 0, l),
	}
}

type FetchDataHandlerBatch func(fetchDataParameter *FetchDataParameterBatch) (err error)

type FetchDataTimeHandlerBatch func(fetchDataParameter *FetchDataParameterTimesBatch) (err error)

type FetchDataHandler func(fetchDataParameter *FetchDataParameter) (err error)

type OperateEveryDatabase struct {
	Db     *gorm.DB `json:"-"`
	DbName string   `json:"db_name"`
	Tables []string `json:"tables"`
}

type OperateEveryDatabaseHandler func(oed *OperateEveryDatabase) (err error)

type SubTreasury interface {
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
