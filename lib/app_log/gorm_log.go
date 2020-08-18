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
	"github.com/sirupsen/logrus"
)

type GOrmLog struct {
	logger *AppLog
}

func NewGOrmLog() (res *GOrmLog) {
	// app_log.GetLog().Logger
	return &GOrmLog{
		logger: GetLog(),
	}
}
func (r GOrmLog) Print(v ...interface{}) () {
	trace_id := ""
	switch v[0] {
	case "sql":
		fields := logrus.Fields{
			"trace_id":      trace_id,
			"type":          "GORM_SQL",
			"rows_returned": v[5],
			"src":           v[1],
			"values":        v[4],
			"duration":      fmt.Sprintf("%dms", v[2].(time.Duration)/1e6),
		}

		r.logger.InfoFields(fields, v[3].(string))
	case "log":
		fields := logrus.Fields{
			"trace_id": trace_id,
			"type":     "GORM_SQL",
			"src":      v[1],
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
