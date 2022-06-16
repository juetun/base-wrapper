package base

import (
	"fmt"
	"gorm.io/gorm"
)

type (
	CommonDb struct {
		Db           *gorm.DB `json:"-"`
		DbName       string   `json:"db_name"`
		TableName    string   `json:"table_name"` // 添加数据对应的表名
		TableComment string   `json:"table_comment"`
	}
	ModelBase interface {
		TableName() string
		GetTableComment() (res string)
	}

	DaoBatchAdd interface {
		BatchAdd(data *BatchAddDataParameter) (err error)
	}

	DaoOneAdd interface {
		AddOneData(parameter *AddOneDataParameter) (err error)
	}

	AddOneDataParameter struct {
		CommonDb
		Model         ModelBase `json:"model"`
		IgnoreColumn  []string  `json:"ignore_column"`   // replace 忽略字段,添加到此字段中的字段不会出现在SQL执行中
		RuleOutColumn []string  `json:"rule_out_column"` // nil时使用默认值，当数据表中存在唯一数据时，此字段的值不会被新的数据替换
	}
	AddOneDataParameterOption func(addOneDataParameter *AddOneDataParameter)

	BatchAddDataParameterOption func(addOneDataParameter *BatchAddDataParameter)

	BatchAddDataParameter struct {
		CommonDb
		IgnoreColumn  []string    `json:"ignore_column"`   // replace 忽略字段,添加到此字段中的字段不会出现在SQL执行中
		RuleOutColumn []string    `json:"rule_out_column"` // nil时使用默认值，当数据表中存在唯一数据时，此字段的值不会被新的数据替换
		Data          []ModelBase `json:"data"`            // 添加的数据
	}
	ActErrorHandlerResult struct {
		CommonDb
		Err   error     `json:"err"`
		Model ModelBase `json:"model"`
	}
	TableSetOption map[string]string
	ActHandlerDao  func() (actRes *ActErrorHandlerResult)
)

func AddOneDataParameterIgnoreColumn(ignoreColumn []string) AddOneDataParameterOption {
	return func(addOneDataParameter *AddOneDataParameter) {
		addOneDataParameter.IgnoreColumn = ignoreColumn
	}
}

func AddOneDataParameterRuleOutColumn(ruleOutColumn []string) AddOneDataParameterOption {
	return func(addOneDataParameter *AddOneDataParameter) {
		addOneDataParameter.RuleOutColumn = ruleOutColumn
	}
}

func (r *ActErrorHandlerResult) GetDbWithTableName(tableAsName ...string) (db *gorm.DB) {
	var tableAsNames = r.TableName
	if len(tableAsName) > 0 {
		tableAsNames = fmt.Sprintf("%s AS %s", r.TableName, tableAsName[0])
	}
	db = r.Db.Table(tableAsNames)
	return
}

func (r *ActErrorHandlerResult) ParseBatchAddDataParameter(options ...BatchAddDataParameterOption) (res *BatchAddDataParameter) {
	res = &BatchAddDataParameter{CommonDb: r.CommonDb,}
	for _, handler := range options {
		handler(res)
	}
	return
}

func BatchAddDataParameterIgnoreColumn(ignoreColumn []string) BatchAddDataParameterOption {
	return func(batchAddDataParameter *BatchAddDataParameter) {
		batchAddDataParameter.IgnoreColumn = ignoreColumn
	}
}

func BatchAddDataParameterData(data []ModelBase) BatchAddDataParameterOption {
	return func(batchAddDataParameter *BatchAddDataParameter) {
		batchAddDataParameter.Data = data
	}
}

func BatchAddDataParameterRuleOutColumn(ruleOutColumn []string) BatchAddDataParameterOption {
	return func(batchAddDataParameter *BatchAddDataParameter) {
		batchAddDataParameter.RuleOutColumn = ruleOutColumn
	}
}
