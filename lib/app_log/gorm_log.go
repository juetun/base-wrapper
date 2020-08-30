/**
* @Author:changjiang
* @Description:
* @File:gorm_log
* @Version: 1.0.0
* @Date 2020/8/18 7:49 下午
 */
package app_log

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/sirupsen/logrus"
)

type GOrmLog struct {
	logger *AppLog
	Db     *gorm.DB
}

func NewGOrmLog(db *gorm.DB) (res *GOrmLog) {
	// app_log.GetLog().Logger
	return &GOrmLog{
		Db:     db,
		logger: GetLog(),
	}
}
func (r GOrmLog) Print(v ...interface{}) () {

	traceId := ""
	if a, ok := r.Db.Get(app_obj.TRACE_ID); ok {
		traceId = fmt.Sprintf("%v", a)
	}
	switch v[0] {
	case "sql":
		fields := logrus.Fields{
			app_obj.TRACE_ID: traceId,
			"type":           "GORM_SQL",
			"rows_returned":  v[5],
			"src":            v[1],
			"values":         v[4],
			"duration":       float64(v[2].(time.Duration) / 1e3), // 时长单位微秒
		}
		r.logger.InfoFields(fields, v[3].(string))
	case "log":
		fields := logrus.Fields{
			app_obj.TRACE_ID: traceId,
			"type":           "GORM_SQL",
			"src":            v[1],
		}
		for _, value := range v[2:] {
			switch value.(type) {
			case *mysql.MySQLError:
				tmp := value.(*mysql.MySQLError)
				r.logger.ErrorFields(fields, fmt.Sprintf("%+v", *tmp))
			default:
				r.logger.InfoFields(fields, fmt.Sprintf("%+v", v[2]))
			}
		}

		// r.logger.InfoFields(fields, v[2].(string))
	}
}
