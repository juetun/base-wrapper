/**
* @Author:changjiang
* @Description:
* @File:gorm_log
* @Version: 1.0.0
* @Date 2020/8/18 7:49 下午
 */
package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/sirupsen/logrus"
)

type GOrmLog struct {
	logger *app_obj.AppLog
	Db     *gorm.DB
	GoPath string
}

func NewGOrmLog(db *gorm.DB) (res *GOrmLog) {
	// app_log.GetLog().Logger
	return &GOrmLog{
		Db:     db,
		logger: app_obj.GetLog(),
		GoPath: os.Getenv("GOPATH"),
	}
}
func (r GOrmLog) Print(v ...interface{}) () {

	traceId := ""
	if a, ok := r.Db.Get(app_obj.TRACE_ID); ok {
		traceId = fmt.Sprintf("%v", a)
	}
	fields := logrus.Fields{
		app_obj.TRACE_ID:      traceId,
		app_obj.APP_FIELD_KEY: "GORMSQL",
		app_obj.APP_LOG_LOC:   v[1].(string),
	}
	if r.GoPath != "" {
		fields[app_obj.APP_LOG_LOC] = "$GOPATH" + strings.TrimPrefix(v[1].(string), r.GoPath)
	}

	switch v[0] {
	case "sql":
		fields["rows_returned"] = v[5]
		fields["duration"] = float64(v[2].(time.Duration) / 1e3) // 时长单位微秒
		r.logger.Info(nil, fields, fmt.Sprintf("SQL:%s [VAL]:%#v", v[3].(string), v[4]))
	case "log":
		for _, value := range v[2:] {
			switch value.(type) {
			case *mysql.MySQLError:
				tmp := value.(*mysql.MySQLError)
				r.logger.Error(nil, fields, fmt.Sprintf("%#v", *tmp))
			default:
				r.logger.Info(nil, fields, fmt.Sprintf("%#v", v[2]))
			}
		}
	}
}
