// Package app_obj
/**
* @Author:changjiang
* @Description:
* @File:db
* @Version: 1.0.0
* @Date 2020/3/27 10:39 下午
 */
package app_obj

import (
	"gorm.io/gorm"
)

const (
	DefaultDbNameSpace = "default"
)

var (
	DbMysql                  = make(map[string]*gorm.DB, 2)
	DistributedMysqlConnects = make([]string, 0) //当前应用支持的分布式数据库连接名
	DistributedRedisConnects = make([]string, 0) //当前应用支持的分布式缓存库连接名
)
