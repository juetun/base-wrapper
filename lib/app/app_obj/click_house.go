package app_obj

import (
	"database/sql"
	"fmt"
)

var DbClickHouse = make(map[string]*ClickHouseClient)

// GetClickHouseClient 获取ClickHouseClient操作实例
func GetClickHouseClient(nameSpace ...string) (client *ClickHouseClient, nameKey string) {
	if len(DbClickHouse) == 0 {
		return
	}
	switch l := len(nameSpace); l {
	case 0:
		nameKey = "default"
	case 1:
		nameKey = nameSpace[0]
	default:
		panic("nameSpace receive at most one parameter")
	}
	if _, ok := DbClickHouse[nameKey]; ok {
		client = DbClickHouse[nameKey]
		return
	}
	panic(fmt.Sprintf("the clickhouse connect(%s) is not exist", nameKey))
}

type ClickHouseClient struct {
	*sql.DB
}
