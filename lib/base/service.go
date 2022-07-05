package base

type (
	ServiceBase struct {
		Context *Context
	}

	//Excel导入服务需要定义的接口，对应服务上需要实现这些方法和调用接口
	ServiceExcelImport interface {
		//excel导入的header关系
		ExcelImportHeaderRelate(args *ArgExcelImportHeaderRelate) (res *ResultExcelImportHeaderRelate, err error)

		//excel导入的参数校验
		ExcelImportValidate(args *ArgExcelImportValidateAndSync) (res []*ExcelImportDataItem, err error)

		//数据同步
		ExcelImportSyncData(args *ArgExcelImportValidateAndSync) (res *ResultExcelImportSyncData, err error)
	}

	ArgExcelImportHeaderRelate struct {
		Scene string `json:"scene" form:"scene"`
	}
	ResultExcelImportHeaderRelate struct {
		Headers []*ExcelImportHeaderRelateItem `json:"headers"`
	}
	ArgExcelImportValidateAndSync struct {
		Scene string                 `json:"scene" form:"scene"`
		Data  []*ExcelImportDataItem `json:"data" form:"data"`
	}
	ExcelImportHeaderRelateItem struct {
		Label      string `json:"label"`       //列中文标题
		ColumnName string `json:"column_name"` //列英文标题
		Index      int64  `json:"index"`       //列序号 如:第一列：0, 第二列：1
	}

	ExcelImportDataItem struct {
		Id             int64  `gorm:"column:id" json:"id"`
		Line           int64  `gorm:"column:line" json:"line"`
		Data           string `gorm:"column:data" json:"data"`
		SheetName      string `gorm:"column:sheet_name" json:"sheet_name"`
		ValidateStatus uint8  `gorm:"-" json:"validate_status"` //验证状态是否通过
		ErrMsg         string `gorm:"-" json:"err_msg"`         //错误信息提示
	}
	ResultExcelImportSyncData struct {
		Result bool `json:"result"`
	}
)

const (
	ExcelImportDataValidateStatusInit     = iota + 1 //导入数据初始化
	ExcelImportDataValidateStatusOk                  //校验成功
	ExcelImportDataValidateStatusFailure             //校验失败
	ExcelImportDataValidateStatusImportOk            //导入完成
)

var (
	SliceExcelImportDataValidateStatus = ModelItemOptions{
		{
			Value: ExcelImportDataValidateStatusInit,
			Label: "初始化",
		},
		{
			Value: ExcelImportDataValidateStatusOk,
			Label: "校验成功",
		},
		{
			Value: ExcelImportDataValidateStatusFailure,
			Label: "校验失败",
		},
		{
			Value: ExcelImportDataValidateStatusImportOk,
			Label: "导入完成",
		},
	}
)

func (r *ServiceBase) SetContext(context ...*Context) (s *ServiceBase) {
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

// RecordLog 记录日志使用
func (r *ServiceBase) RecordLog(locKey string, logContent map[string]interface{}, err error, needRecordInfo ...bool) {
	if err != nil {
		logContent["err"] = err.Error()
		r.Context.Error(logContent, locKey)
		return
	}
	if len(needRecordInfo) > 0 {
		r.Context.Info(logContent, locKey)
	}
}
