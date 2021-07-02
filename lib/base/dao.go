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

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"gorm.io/gorm"
)

type ServiceDao struct {
	Context *Context
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

type ModelBase interface {
	TableName() string
}
type DaoBatchAdd interface {
	BatchAdd(data *BatchAddDataParameter) (err error)
}
type AddOneDataParameter struct {
	DbName        string    `json:"db_name"`
	Db            *gorm.DB  `json:"-"`
	Model         ModelBase `json:"model"`      // 添加的数据
	TableName     string    `json:"table_name"` // 添加数据对应的表名
	RuleOutColumn []string  `json:"rule_out_column"`
}

type BatchAddDataParameter struct {
	DbName        string      `json:"db_name"`
	Db            *gorm.DB    `json:"-"`
	TableName     string      `json:"table_name"` // 添加数据对应的表名
	RuleOutColumn []string    `json:"rule_out_column"`
	Data          []ModelBase `json:"data"` // 添加的数据
}

func (r *ServiceDao) BatchAdd(data *BatchAddDataParameter) (err error) {
	if len(data.Data) == 0 {
		return
	}
	var (
		columns, replaceKeys, vv []string
		val                      = make([]interface{}, 0, len(data.Data)*20)
	)

	for k, item := range data.Data {
		if k == 0 {
			if data.TableName == "" {
				data.TableName = item.TableName()
			}
			columns, replaceKeys = r.getItemColumns(data.Db, item, val, vv)
		} else {
			r.getItemColumns(data.Db, item, val, vv)
		}
	}
	sql := fmt.Sprintf("INSERT INTO `%s`(`"+strings.Join(columns, "`,`")+"`) VALUES ("+strings.Join(vv, ",")+
		") ON DUPLICATE KEY UPDATE "+strings.Join(replaceKeys, ","), data.TableName)

	if err = data.Db.Exec(sql, val...).Error; err != nil {
		logContent := map[string]interface{}{
			"sql":  sql,
			"data": data,
			"val":  val,
			"err":  err,
		}
		if data.DbName != "" {
			logContent[app_obj.DbNameKey] = data.DbName
		}
		r.Context.Error(logContent, "ServiceDaoBatchAdd")
		return
	}
	return
}

func (r *ServiceDao) getItemColumns(db *gorm.DB, item ModelBase, val []interface{}, vv []string) (columns, replaceKeys []string) {
	types := reflect.TypeOf(item)
	values := reflect.ValueOf(item)
	fieldNum := types.Elem().NumField()
	columns = make([]string, 0, fieldNum)
	replaceKeys = make([]string, 0, fieldNum)
	val = make([]interface{}, 0, fieldNum)
	var tag string
	var tagValue = "gorm"
	for i := 0; i < fieldNum; i++ {
		tag = r.GetColumnName(types.Elem().Field(i).Tag.Get(tagValue))
		columns = append(columns, tag)
		if tag == "id" || tag == "created_at" {
			continue
		}

		val = append(val, r.formatValue(db, values.Elem().Field(i)))
		replaceKeys = append(replaceKeys, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))
		vv = append(vv, "?")
	}
	return
}

func (r *ServiceDao) AddOneData(parameter *AddOneDataParameter) (err error) {
	if parameter.TableName == "" {
		parameter.TableName = parameter.Model.TableName()
	}
	if parameter.Db == nil {
		parameter.Db = r.Context.Db
		parameter.DbName = r.Context.DbName
	}
	values := reflect.ValueOf(parameter.Model)
	tagValue := "gorm"
	types := reflect.TypeOf(parameter.Model)
	fieldNum := types.Elem().NumField()
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
		val = append(val, r.formatValue(parameter.Db, valueStruct))
		columns = append(columns, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))
		vv = append(vv, "?")
	}
	sql := fmt.Sprintf("INSERT INTO `%s`(`"+strings.Join(keys, "`,`")+"`) VALUES ("+strings.Join(vv, ",")+
		") ON DUPLICATE KEY UPDATE "+strings.Join(columns, ","), parameter.TableName)

	if err = parameter.Db.Exec(sql, val...).Error; err != nil {
		logContent := map[string]interface{}{
			"sql":  sql,
			"data": parameter,
			"val":  val,
			"err":  err,
		}
		if parameter.DbName != "" {
			logContent[app_obj.DbNameKey] = parameter.DbName
		}
		r.Context.Error(logContent, "DaoUserEmailImplAdd")
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
