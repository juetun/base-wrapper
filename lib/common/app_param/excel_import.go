package app_param

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
)

type (
	ControllerExcelImport interface {
		ExcelImportHeaderRelate(c *gin.Context)

		//excel导入的参数校验
		ExcelImportValidate(c *gin.Context)

		//数据同步
		ExcelImportSyncData(c *gin.Context)
	}

	//Excel导入服务需要定义的接口，对应服务上需要实现这些方法和调用接口
	ServiceExcelImport interface {
		//excel导入的header关系
		ExcelImportHeaderRelate(args *ArgExcelImportHeaderRelate) (res *ResultExcelImportHeaderRelate, err error)

		//excel导入的参数校验
		ExcelImportValidate(args *ArgExcelImportValidateAndSync) (res []ExcelImportDataItem, err error)

		//数据同步
		ExcelImportSyncData(args *ArgExcelImportValidateAndSync) (res []ExcelImportDataItem, err error)
	}

	ArgExcelImportHeaderRelate struct {
		Scene string `json:"scene" form:"scene"`
	}
	ResultExcelImportHeaderRelate struct {
		Headers map[string]ExcelImportHeaderRelateItem `json:"headers"`
	}
	ArgExcelImportValidateAndSync struct {
		Scene string                `json:"scene" form:"scene"`
		Data  []ExcelImportDataItem `json:"data" form:"data"`
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

)

const (
	ExcelImportDataValidateStatusInit     = iota + 1 //导入数据初始化
	ExcelImportDataValidateStatusOk                  //校验成功
	ExcelImportDataValidateStatusFailure             //校验失败
	ExcelImportDataValidateStatusImportOk            //导入完成
)

var (
	SliceExcelImportDataValidateStatus = base.ModelItemOptions{
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

func (r *ArgExcelImportValidateAndSync) ToJson() (res []byte, err error) {
	if r == nil {
		r = &ArgExcelImportValidateAndSync{}
	}
	res, err = json.Marshal(r)
	return
}

func (r *ExcelImportDataItem) GetId() (res int64) {
	return r.Id
}

func ExcelImportHeaderRelate(c *gin.Context, srv ServiceExcelImport) (data *ResultExcelImportHeaderRelate, err error) {
	var (
		arg ArgExcelImportHeaderRelate
	)
	data = &ResultExcelImportHeaderRelate{}
	if err = c.Bind(&arg); err != nil {
		return
	}
	if data, err = srv.ExcelImportHeaderRelate(&arg); err != nil {
		return
	}
	return
}

func ExcelImportValidate(c *gin.Context, srv ServiceExcelImport) (data []ExcelImportDataItem, err error) {

	var (
		arg ArgExcelImportValidateAndSync
	)
	data = make([]ExcelImportDataItem, 0)
	if err = c.Bind(&arg); err != nil {
		return
	}
	if data, err = srv.ExcelImportValidate(&arg); err != nil {
		return
	}
	return
}

func ExcelImportSyncData(c *gin.Context, srv ServiceExcelImport) (data []ExcelImportDataItem, err error) {
	var (
		arg ArgExcelImportValidateAndSync
	)
	data = []ExcelImportDataItem{}
	if err = c.Bind(&arg); err != nil {
		return
	}
	if data, err = srv.ExcelImportSyncData(&arg); err != nil {
		return
	}
	return
}
