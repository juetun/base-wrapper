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
	"runtime"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
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
	Model         ModelBase `json:"model"`           // 添加的数据
	TableName     string    `json:"table_name"`      // 添加数据对应的表名
	IgnoreColumn  []string  `json:"ignore_column"`   // replace 忽略字段,添加到此字段中的字段不会出现在SQL执行中
	RuleOutColumn []string  `json:"rule_out_column"` // nil时使用默认值，当数据表中存在唯一数据时，此字段的值不会被新的数据替换
}

type BatchAddDataParameter struct {
	DbName        string      `json:"db_name"`
	Db            *gorm.DB    `json:"-"`
	TableName     string      `json:"table_name"`      // 添加数据对应的表名
	IgnoreColumn  []string    `json:"ignore_column"`   // replace 忽略字段,添加到此字段中的字段不会出现在SQL执行中
	RuleOutColumn []string    `json:"rule_out_column"` // nil时使用默认值，当数据表中存在唯一数据时，此字段的值不会被新的数据替换
	Data          []ModelBase `json:"data"`            // 添加的数据
}

// CreateTableWithError 判断SQL err 如果表不存在，则创建表
func (r *ServiceDao) CreateTableWithError(db *gorm.DB, tableName string, e, model interface{}) (err error) {

	var file string
	// 获取上层调用者PC，文件名，所在行	// 拼接文件名与所在行
	if _, codePath, codeLine, ok := runtime.Caller(1); ok {
		file = fmt.Sprintf("%s(l:%d)", codePath, codeLine) // runtime.FuncForPC(pc).Name(),
	}
	r.Context.Error(map[string]interface{}{
		"err": e,
		"src": file,
	}, "DaoUserImplCreateUserTable")
	// 延迟处理的函数
	switch e.(type) {
	case *mysql.MySQLError: // 运行时错误
		me := e.(*mysql.MySQLError)
		if me.Number == 1146 {
			if err = db.Table(tableName).Migrator().
				CreateTable(model); err != nil {
				return
			}
			return
		}
		err = fmt.Errorf(me.Error())
	default:
		err = fmt.Errorf("数据异常,请重试(102)")
		return
	}
	return
}
func (r *ServiceDao) BatchAdd(data *BatchAddDataParameter) (err error) {
	if len(data.Data) == 0 {
		return
	}
	r.defaultBatchAddDataParameter(data)
	var (
		columns, replaceKeys []string
		l                    = len(data.Data) * 20
		vv                   = make([]string, 0, l)
		val                  = make([]interface{}, 0, l)
	)

	for k, item := range data.Data {
		if k == 0 {
			if data.TableName == "" {
				data.TableName = item.TableName()
			}
			columns, replaceKeys = r.getItemColumns(data.IgnoreColumn, data.RuleOutColumn, data.Db, item, &val, &vv)
		} else {
			_, _ = r.getItemColumns(data.IgnoreColumn, data.RuleOutColumn, data.Db, item, &val, &vv)
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

func (r *ServiceDao) InStringSlice(s string, slice []string) (res bool) {
	if s == "" {
		return
	}
	for _, s2 := range slice {
		if s2 == s {
			res = true
			return
		}
	}
	return
}
func (r *ServiceDao) getItemColumns(ignoreColumn, ruleOutColumn []string, db *gorm.DB, item ModelBase, val *[]interface{}, vv *[]string) (columns, replaceKeys []string) {
	types := reflect.TypeOf(item)
	values := reflect.ValueOf(item)
	fieldNum := types.Elem().NumField()
	columns = make([]string, 0, fieldNum)
	replaceKeys = make([]string, 0, fieldNum)
	var tag string
	var tagValue = "gorm"
	for i := 0; i < fieldNum; i++ {
		tag = r.GetColumnName(types.Elem().Field(i).Tag.Get(tagValue))
		if r.InStringSlice(tag, ignoreColumn) {
			continue
		}
		columns = append(columns, tag)
		*val = append(*val, r.formatValue(db, values.Elem().Field(i)))
		*vv = append(*vv, "?")
		if r.InStringSlice(tag, ruleOutColumn) {
			continue
		}
		replaceKeys = append(replaceKeys, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))

	}
	return
}

func (r *ServiceDao) defaultRuleOutColumn() []string {
	return []string{"created_at"}
}

func (r *ServiceDao) defaultIgnoreColumn() []string {
	return []string{"id"}
}

func (r *ServiceDao) defaultBatchAddDataParameter(parameter *BatchAddDataParameter) {
	if parameter.TableName == "" {
		if len(parameter.Data) > 0 {
			parameter.TableName = parameter.Data[0].TableName()
		}
	}
	if parameter.Db == nil {
		parameter.Db = r.Context.Db
		parameter.DbName = r.Context.DbName
	}

	if parameter.RuleOutColumn == nil {
		parameter.RuleOutColumn = r.defaultRuleOutColumn()
	}
	if parameter.IgnoreColumn == nil {
		parameter.IgnoreColumn = r.defaultIgnoreColumn()
	}
}

func (r *ServiceDao) defaultAddOneDataParameter(parameter *AddOneDataParameter) {
	if parameter.TableName == "" {
		parameter.TableName = parameter.Model.TableName()
	}
	if parameter.Db == nil {
		parameter.Db = r.Context.Db
		parameter.DbName = r.Context.DbName
	}
	if parameter.RuleOutColumn == nil {
		parameter.RuleOutColumn = r.defaultRuleOutColumn()
	}
	if parameter.IgnoreColumn == nil {
		parameter.IgnoreColumn = r.defaultIgnoreColumn()
	}
}

func (r *ServiceDao) AddOneData(parameter *AddOneDataParameter) (err error) {
	r.defaultAddOneDataParameter(parameter)
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
		if r.InStringSlice(tag, parameter.IgnoreColumn) {
			continue
		}

		keys = append(keys, tag)
		valueStruct = values.Elem().Field(i)
		val = append(val, r.formatValue(parameter.Db, valueStruct))
		vv = append(vv, "?")
		if r.InStringSlice(tag, parameter.RuleOutColumn) {
			continue
		}
		columns = append(columns, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))
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
