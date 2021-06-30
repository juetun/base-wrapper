// Package base
/**
* @Author:ChangJiang
* @Description:
* @File:dao
* @Version: 1.0.0
* @Date 2020/4/5 8:22 下午
 */
package base

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ServiceDao struct {
	Context *Context
}
type ModelBase interface {
	TableName() string
}

func (r *ServiceDao) SetContext(context ...*Context) (s *ServiceDao) {
	for _, cont := range context {
		cont.InitContext()
	}
	switch len(context) {
	case 0:
		r.Context = NewContext()
		break
	case 1:
		r.Context = context[0]
		break
	default:
		panic("你传递的参数当前不支持")
	}

	return r
}
func (r *ServiceDao) formatValue(db *gorm.DB, valueStruct reflect.Value) (res interface{}) {

	switch valueStruct.Kind() {

	case reflect.Interface:
		res = valueStruct.Interface()
	case reflect.Ptr:
		if valueStruct.IsNil() {
			res = nil
			return
		}
		return r.formatValue(db, valueStruct.Elem())
	case reflect.Bool:
		res = valueStruct.Bool()
	case reflect.String:
		res = valueStruct.String()
	case reflect.Int, reflect.Int32, reflect.Int64:
		res = valueStruct.Int()
	case reflect.Float32, reflect.Float64:
		res = valueStruct.Float()
	default:
		switch valueStruct.Type().String() {
		case "base.TimeNormal":
			dt := valueStruct.Interface().(TimeNormal)
			res = dt.Format("2006-01-02 15:04:05")
		case "time.Time":
			res = valueStruct.Interface().(time.Time).Format("2006-01-02 15:04:05")
		case "time.Duration":
			res = valueStruct.Interface().(time.Duration).String()
		case "int":
			res = valueStruct.Int()
		default:
			res = valueStruct.String()
		}
	}
	return
}
func (r *ServiceDao) AddOneData(db *gorm.DB, data ModelBase, tableName ...string) (err error) {
	tbName := ""
	if len(tableName) > 0 {
		tbName = tableName[0]
	} else {
		tbName = data.TableName()
	}
	values := reflect.ValueOf(data)
	tagValue := "gorm"
	types := reflect.TypeOf(data)
	var fieldNum int
	fieldNum = types.Elem().NumField()
	var valueStruct reflect.Value

	keys := make([]string, 0, fieldNum)
	columns := make([]string, 0, fieldNum)
	val := make([]interface{}, 0, fieldNum)
	vv := make([]string, 0, fieldNum)
	var tag string
	for i := 0; i < fieldNum; i++ {
		tag = r.GetColumnName(types.Elem().Field(i).Tag.Get(tagValue))
		if tag == "id" || tag == "created_at" {
			continue
		}
		keys = append(keys, tag)
		valueStruct = values.Elem().Field(i)
		val = append(val, r.formatValue(db, valueStruct))
		columns = append(columns, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))
		vv = append(vv, "?")
	}
	sql := fmt.Sprintf("INSERT INTO `%s`(`"+strings.Join(keys, "`,`")+"`) VALUES ("+strings.Join(vv, ",")+
		") ON DUPLICATE KEY UPDATE "+strings.Join(columns, ","), tbName)

	if err = db.Exec(sql, val...).Error; err != nil {
		r.Context.Error(map[string]interface{}{
			"sql":  sql,
			"data": data,
			"val":  val,
			"err":  err,
		}, "DaoUserEmailImplAdd")
		return
	}
	return
}

func (r *ServiceDao) GetColumnName(s string) (res string) {
	li := strings.Split(s, ";")
	res = s
	for _, s2 := range li {
		if s2 == "" {
			return
		}
		li1 := strings.Split(s2, ":")
		if len(li1) > 1 && li1[0] == "column" {
			res = li1[1]
		}
	}
	return
}
