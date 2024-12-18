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
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/go-sql-driver/mysql"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"gorm.io/gorm"
)

//mysql 错误的码
const (
	MysqlErrorTableNotExists = 1146 //表不存在
)

type (
	Dao interface {
		// BatchAdd 批量添加数据
		BatchAdd(data *BatchAddDataParameter) (err error)

		// AddOneData 单条添加数据
		AddOneData(parameter *AddOneDataParameter) (err error)

		// ActErrorHandler 操作逻辑 如果表不存在，创建表
		ActErrorHandler(actHandler ActHandlerDao) (err error)

		// ScopesDeletedAt 添加deleted_at字段WHERE条件拼接
		ScopesDeletedAt(prefix ...string) func(db *gorm.DB) *gorm.DB

		ScopesPager(pager *response.Pager) func(db *gorm.DB) *gorm.DB

		// GetDefaultDb 获取默认的DB操作
		GetDefaultDb(modelBase ModelBase) (res CommonDb)

		// GetDefaultAddOneDataParameter 获取插入数据的对象
		GetDefaultAddOneDataParameter(modelBase ModelBase) (res *AddOneDataParameter)

		// GetDefaultBatchAddDataParameter 获取批量插入数据操作的类型
		GetDefaultBatchAddDataParameter(modelBase ...ModelBase) (res *BatchAddDataParameter, err error)

		// GetDefaultActErrorHandlerResult 获取默认对象
		GetDefaultActErrorHandlerResult(model ModelBase) (res *ActErrorHandlerResult)

		// CreateTableWithError 创建表
		CreateTableWithError(db *gorm.DB, tableName string, e, model interface{}, comment ...TableSetOption) (err error)

		RefreshCache(refreshCache ...bool) (res bool)

		// RecordLog 记录日志逻辑实现
		RecordLog(message string, logContent map[string]interface{}, err error, needRecordInfo ...bool)
	}
	ServiceDao struct {
		Context *Context
	}
)

func (r *ServiceDao) GetDefaultAddOneDataParameter(modelBase ModelBase) (res *AddOneDataParameter) {
	res = &AddOneDataParameter{Model: modelBase}
	res.CommonDb = r.GetDefaultDb(modelBase)
	return
}

func RealRandomNumber(maxWeight int64) (result int64, err error) {
	if maxWeight == 0 {
		result = 0
		return
	}
	//生成[0, maxWeight) 范围的真随机数。
	bigInt := big.NewInt(maxWeight)
	var res *big.Int
	if res, err = rand.Int(rand.Reader, bigInt); err != nil {
		return
	}
	result = res.Int64()
	//fmt.Printf("随机数 maxWeight:%d result:%d",maxWeight,result)
	return
}

//真随机数生成
//生成规则参考，https://blog.csdn.net/study_in/article/details/102919019
func (r *ServiceDao) RealRandomNumber(maxWeight int64) (result int64, err error) {
	return RealRandomNumber(maxWeight)
}

func (r *ServiceDao) GetDefaultBatchAddDataParameter(modelBase ...ModelBase) (res *BatchAddDataParameter, err error) {
	if len(modelBase) == 0 {
		err = fmt.Errorf("您没有选择要添加的数据")
		return
	}
	res = &BatchAddDataParameter{Data: modelBase,}
	res.CommonDb = r.GetDefaultDb(modelBase[0])
	return
}

// ActErrorHandler 操作(当前实现逻辑 如果报指定状态，则创建表)
func (r *ServiceDao) ActErrorHandler(actHandler ActHandlerDao) (err error) {
	var (
		res            *ActErrorHandlerResult
		tableSetOption = TableSetOption{}
	)
	if res = actHandler(); res.Err != nil {
		if res.Model == nil {
			err = fmt.Errorf("您没有选择创建表的结构体")
			return
		}
		if res.TableName == "" {
			res.TableName = res.Model.TableName()
		}
		if res.TableComment == "" {
			res.TableComment = res.Model.GetTableComment()
		}

		if res.TableComment != "" {
			tableSetOption["COMMENT"] = res.TableComment
		}
		if len(tableSetOption) > 0 {
			err = r.CreateTableWithError(res.Db, res.TableName, res.Err, res.Model, tableSetOption)
		} else {
			err = r.CreateTableWithError(res.Db, res.TableName, res.Err, res.Model)
		}
		if err != nil {
			return
		}
		if res = actHandler(); res.Err != nil {
			err = res.Err
			return
		}
	}

	return
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

func (r *ActErrorHandlerResult) ParseAddOneDataParameter(options ...AddOneDataParameterOption) (res *AddOneDataParameter) {
	res = &AddOneDataParameter{}
	res.CommonDb = r.CommonDb
	res.Model = r.Model
	for _, handler := range options {
		handler(res)
	}
	return
}

func (r *ServiceDao) GetDefaultActErrorHandlerResult(model ModelBase) (res *ActErrorHandlerResult) {
	res = &ActErrorHandlerResult{
		CommonDb: r.GetDefaultDb(model),
		Model:    model,
	}
	return
}

func (r *ServiceDao) GetDefaultDb(modelBase ModelBase) (res CommonDb) {
	res = CommonDb{
		Db:        r.Context.Db,
		DbName:    r.Context.DbName,
		TableName: modelBase.TableName(),
	}
	return
}

// InitFetchParameters 初始化FetchDataParameter
func (r *ServiceDao) InitFetchParameters(model ModelBase) (fetchData *FetchDataParameter) {
	fetchData = &FetchDataParameter{
		CommonDb: CommonDb{
			Db:        r.Context.Db,
			DbName:    r.Context.DbName,
			TableName: model.TableName(),
		},
	}
	return
}

func (r *ServiceDao) formatValue(valueStruct reflect.Value) (res interface{}) {

	switch valueStruct.Kind() {

	case reflect.Interface:
		res = valueStruct.Interface()
	case reflect.Ptr:
		if valueStruct.IsNil() {
			res = nil
			return
		}
		return r.formatValue(valueStruct.Elem())
	case reflect.Bool:
		res = valueStruct.Bool()
	case reflect.String:
		res = valueStruct.String()
	case reflect.Int8, reflect.Int, reflect.Int32, reflect.Int64:
		res = valueStruct.Int()
	case reflect.Uint8, reflect.Uint, reflect.Uint32, reflect.Uint64:
		res = valueStruct.Uint()
	case reflect.Float32, reflect.Float64:
		res = valueStruct.Float()
	default:
		switch valueStruct.Type().String() {
		case "base.TimeNormal":
			dt := valueStruct.Interface().(TimeNormal)
			res = dt.Format(utils.DateTimeGeneral)
		case "time.Time":
			res = valueStruct.Interface().(time.Time).Format(utils.DateTimeGeneral)
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

// ScopesPager 分页操作条件
func (r *ServiceDao) ScopesPager(pager *response.Pager) func(db *gorm.DB) *gorm.DB {
	return ScopesPager(pager)
}

func ScopesPager(pager *response.Pager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pager.PagerParameter.GetOffset()).Limit(pager.PageSize)
	}
}

// ScopesDeletedAt 查询条件添加删除条件为空
func (r *ServiceDao) ScopesDeletedAt(prefix ...string) func(db *gorm.DB) *gorm.DB {
	return ScopesDeletedAt(prefix...)
}

func ScopesDeletedAt(prefix ...string) func(db *gorm.DB) *gorm.DB {
	pre := ""
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	return func(db *gorm.DB) *gorm.DB {
		if pre != "" {
			db = db.Where(fmt.Sprintf("%s.`deleted_at` IS NULL", pre))
		} else {
			db = db.Where("`deleted_at` IS NULL")
		}

		return db
	}

}

func (r *ServiceDao) mysqlCreateError(e error) (err error) {
	err = e
	switch e.(type) {
	case *mysql.MySQLError: // 运行时错误
		me := e.(*mysql.MySQLError)
		if me.Number == MysqlErrorTableNotExists {
			err = nil
			return
		}
	}
	return
}

func (r *ServiceDao) createTableAct(db *gorm.DB, tableName string, model interface{}, comment ...TableSetOption) (err error) {
	dba := db.Table(tableName)
	for _, option := range comment {
		for s, s2 := range option {
			dba = dba.Set("gorm:table_options", fmt.Sprintf("%s='%s'", s, s2))
		}
	}

	//执行创建表
	if err = dba.Migrator().CreateTable(model); err != nil {
		//如果是创建表操作,如果返回的错误还是表不存在，则跳过
		err = r.mysqlCreateError(err)
		return
	}
	return
}

// CreateTableWithError 判断SQL err 如果表不存在，则创建表
func (r *ServiceDao) CreateTableWithError(db *gorm.DB, tableName string, e, model interface{}, comment ...TableSetOption) (err error) {

	var file string
	// 获取上层调用者PC，文件名，所在行	// 拼接文件名与所在行
	if _, codePath, codeLine, ok := runtime.Caller(1); ok {
		file = fmt.Sprintf("%s(l:%d)", codePath, codeLine) // runtime.FuncForPC(pc).Name(),
	}
	logContent := map[string]interface{}{"src": file, "e": fmt.Sprintf("%+v", e)}
	defer func() {
		if err != nil {
			logContent["err"] = err.Error()
			r.Context.Error(logContent, "ServiceDaoCreateTableWithError")
		}
	}()
	// 延迟处理的函数
	switch e.(type) {
	case *mysql.MySQLError: // 运行时错误
		me := e.(*mysql.MySQLError)
		switch me.Number {
		case MysqlErrorTableNotExists:
			err = r.createTableAct(db, tableName, model, comment...)
			return
		}

		err = fmt.Errorf(me.Error())
	case *ErrorRuntimeStruct:
		err = e.(*ErrorRuntimeStruct)
		return
	default:
		err = fmt.Errorf("数据异常,请重试(102)")
		return
	}
	return
}

func newDataModal(data *BatchAddDataParameter) (res *dataModal) {
	l := len(data.Data) * defaultMaxColumn
	res = &dataModal{
		val:         make([]interface{}, 0, l),
		columns:     make([]string, 0, defaultMaxColumn),
		replaceKeys: make([]string, 0, defaultMaxColumn),
	}
	return
}

func (r *ServiceDao) getItemColumns(dataModal *dataModal) (err error) {

	dataModalReflect := &dataModalReflect{
		Types:  reflect.TypeOf(dataModal.modal),
		Values: reflect.ValueOf(dataModal.modal),
	}
	if err = r.getKind(dataModal, dataModalReflect); err != nil {
		return
	}
	return
}

func (r *ServiceDao) batchAddAct(data *BatchAddDataParameter, logContent *map[string]interface{}) (err error) {

	var (
		l         = len(data.Data) * 20
		vl        = make([]string, 0, len(data.Data))
		dataModal *dataModal
	)

	dataModal = newDataModal(data)
	dataModal.ignoreColumn = data.IgnoreColumn
	dataModal.ruleOutColumn = data.RuleOutColumn

	for k, item := range data.Data {
		if k == 0 && data.TableName == "" {
			data.TableName = item.TableName()
		}
		dataModal.vv = make([]string, 0, l)
		dataModal.modal = item
		dataModal.ind = k
		if err = r.getItemColumns(dataModal); err != nil {
			return
		}
		vvs := fmt.Sprintf("(%s)", strings.Join(dataModal.vv, ","))
		vl = append(vl, vvs)
	}
	err = r.execNotNeedReturn(dataModal, data, logContent, vl)

	return
}

func (r *ServiceDao) execNotNeedReturn(dataModal *dataModal, data *BatchAddDataParameter, logContent *map[string]interface{}, vl []string) (err error) {
	//sql := fmt.Sprintf("INSERT INTO `%s`(`"+strings.Join(dataModal.columns, "`,`")+"`) VALUES "+strings.Join(vl, ",")+
	//	" ON DUPLICATE KEY UPDATE "+strings.Join(dataModal.replaceKeys, ","), data.TableName)
	sql := fmt.Sprintf("INSERT INTO `%s`(`%s`) VALUES %s ON DUPLICATE KEY UPDATE %s",
		data.TableName,
		strings.Join(dataModal.columns, "`,`"),
		strings.Join(vl, ","),
		strings.Join(dataModal.replaceKeys, ","),
	)
	(*logContent)["sql"] = sql
	(*logContent)["val"] = dataModal.val
	if err = data.Db.Exec(sql, dataModal.val...).Error; err != nil {
		if data.DbName != "" {
			(*logContent)[app_obj.DbNameKey] = data.DbName
		}
		return
	}
	return
}
func (r *ServiceDao) BatchAdd(data *BatchAddDataParameter) (err error) {
	if len(data.Data) == 0 {
		return
	}
	logContent := map[string]interface{}{"data": data,}
	defer func() {
		if err == nil || r.Context == nil {
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, "ServiceDaoBatchAdd")
	}()
	r.defaultAddDataParameter(&data.AddDataParameter, data.Data[0])
	err = r.batchAddAct(data, &logContent)
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

const GORMTAG = "gorm"

type (
	dataModal struct {
		val           []interface{}
		vv            []string //？占位符
		ignoreColumn  []string
		ruleOutColumn []string
		columns       []string
		replaceKeys   []string
		fieldNum      int
		ind           int
		modal         ModelBase
	}
	dataModalReflect struct {
		Types  reflect.Type
		Values reflect.Value
	}
)

const (
	defaultMaxColumn = 80
)

//获取ModelBase对象的所有gorm 支持字段
func (r *ServiceDao) GetAllColumn(modal ModelBase) (columns []string, err error) {
	dataModalReflect := &dataModalReflect{
		Types:  reflect.TypeOf(modal),
		Values: reflect.ValueOf(modal),
	}
	dataModal := &dataModal{
		val:         make([]interface{}, 0, defaultMaxColumn),
		columns:     make([]string, 0, defaultMaxColumn),
		replaceKeys: make([]string, 0, defaultMaxColumn),
	}
	if err = r.getKind(dataModal, dataModalReflect); err != nil {
		return
	}
	columns = dataModal.columns
	return
}

func (r *ServiceDao) getKind(dataModal *dataModal, dataModalReflectObj *dataModalReflect) (err error) {
	//if dataModalReflectObj.Types.IsVariadic() {
	//	err = fmt.Errorf("参数格式错误")
	//	return
	//}
	kind := dataModalReflectObj.Types.Kind()

	switch kind {
	case reflect.Struct: //如果是结构体
		var (
			ignoreColumnFlag bool
			tag              string
		)
		for i := 0; i < dataModalReflectObj.Types.NumField(); i++ {
			field := dataModalReflectObj.Types.Field(i)
			value := dataModalReflectObj.Values.Field(i)

			if ignoreColumnFlag, tag = r.GetColumnName(field.Tag.Get(GORMTAG), field.Name); ignoreColumnFlag { // 如果是tag标记忽略字段
				continue
			}
			dataModalReflectTmp := &dataModalReflect{Types: field.Type, Values: value,}
			if tag == "" {
				err = r.getKind(dataModal, dataModalReflectTmp)
				continue
			}

			if r.InStringSlice(tag, dataModal.ignoreColumn) {
				continue
			}

			dataModal.val = append(dataModal.val, r.formatValue(dataModalReflectTmp.Values))
			dataModal.vv = append(dataModal.vv, "?")
			if dataModal.ind == 0 { //如果是第一条数据获取字段信息
				dataModal.columns = append(dataModal.columns, tag)
				if !r.InStringSlice(tag, dataModal.ruleOutColumn) {
					dataModal.replaceKeys = append(dataModal.replaceKeys, fmt.Sprintf("`%s`=VALUES(`%s`)", tag, tag))
				}
			}
		}
	case reflect.Ptr: //如果是指针
		dataModalReflectTmp := &dataModalReflect{
			Types:  dataModalReflectObj.Types.Elem(),
			Values: dataModalReflectObj.Values.Elem(),
		}
		err = r.getKind(dataModal, dataModalReflectTmp)
	default:
		dataModal.val = append(dataModal.val, r.formatValue(dataModalReflectObj.Values))
	}
	return
}

func (r *ServiceDao) defaultRuleOutColumn() []string {
	return []string{"created_at"}
}

func (r *ServiceDao) defaultIgnoreColumn() []string {
	return []string{"id"}
}

func (r *ServiceDao) defaultAddDataParameter(parameter *AddDataParameter, model ModelBase) {

	if parameter.TableName == "" {
		parameter.TableName = model.TableName()

	}

	if parameter.Db == nil && r.Context != nil {
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
	logContent := map[string]interface{}{"parameter": parameter,}
	defer func() {
		if err == nil || r.Context == nil {
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, "ServiceAddOneData")
	}()
	if parameter.Model == nil {
		err = fmt.Errorf("您要添加的数据为空")
		return
	}

	r.defaultAddDataParameter(&parameter.AddDataParameter, parameter.Model)
	batchAddDataParameter := &BatchAddDataParameter{}
	batchAddDataParameter.AddDataParameter = parameter.AddDataParameter
	batchAddDataParameter.Data = []ModelBase{parameter.Model}
	err = r.batchAddAct(batchAddDataParameter, &logContent)
	return
}

// GetColumnName 获取字段对应的 key
// param s string 			对象的tag值
// param fieldName string  	对象的字段名
func (r *ServiceDao) GetColumnName(s, fieldName string) (ignoreColumn bool, res string) {
	if s == "" {
		return
	}
	li := strings.Split(s, ";")
	for _, s2 := range li {
		if s2 == "" {
			return
		}
		if s2 == "-" { // 如果是忽略字段
			ignoreColumn = true
			res = r.getDefaultColumnValue(fieldName)
			return
		}
		li1 := strings.Split(s2, ":")
		if len(li1) > 1 && li1[0] == "column" {
			if li1[1] != "-" { // "column:-"非格式的获取
				res = li1[1]
				return
			} else { // 如果是忽略字段
				ignoreColumn = true
				return
			}
			// res = r.getDefaultColumnValue(fieldName)
			// return
		}
	}
	if res == "" {
		res = r.getDefaultColumnValue(fieldName)
		return
	}
	res = s
	return
}

func (r *ServiceDao) getDefaultColumnValue(name string) (res string) {
	if name == "ID" {
		res = "id"
		return
	}
	var buffer = make([]rune, 0, len(name)+50)
	for i, bt := range name {
		if unicode.IsUpper(bt) {
			if i != 0 {
				buffer = append(buffer, '_')
			}
			buffer = append(buffer, unicode.ToLower(bt))
		} else {
			buffer = append(buffer, bt)
		}
	}
	return string(buffer)
}

// RecordLog 记录日志使用
func (r *ServiceDao) RecordLog(message string, logContent map[string]interface{}, err error, needRecordInfo ...bool) {
	if err != nil {
		logContent["err"] = err.Error()
		r.Context.Error(logContent, message)
		return
	}
	if len(needRecordInfo) > 0 {
		r.Context.Info(logContent, message)
	}
}

// RefreshCache 封装是否刷新缓存逻辑
func (r *ServiceDao) RefreshCache(refreshCache ...bool) (res bool) {
	if len(refreshCache) > 0 {
		res = refreshCache[0]
	}
	return
}
