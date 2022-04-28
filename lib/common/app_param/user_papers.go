package app_param

import (
	"encoding/json"
	"fmt"
)

//用户资质是否需要填写类型
const (
	PaperMustDateNotNeed = iota //不需要时间
	PaperMustDateYes            //必须填写时间
	PaperMustDateNo             //可不填时间
)

var MapMustDate = map[uint8]string{
	PaperMustDateNotNeed: "不填",
	PaperMustDateYes:     "必填",
	PaperMustDateNo:      "可不填",
}

const (
	DataPapersGroupShopPropertyRadio    = "radio"    //单选
	DataPapersGroupShopPropertyCheckbox = "checkbox" //多选
	DataPapersGroupShopPropertySelect   = "select"   //下拉菜单
)

var MapDataPapersGroupShopProperty = map[string]string{
	DataPapersGroupShopPropertyRadio:    "单选",
	DataPapersGroupShopPropertyCheckbox: "多选",
	DataPapersGroupShopPropertySelect:   "下拉菜单",
}


type (
	DataPapersGroupShopProperty struct {
		ShowType     string `json:"show_type"` //DataPapersGroupShopPropertyRadio DataPapersGroupShopPropertyCheckbox  DataPapersGroupShopPropertySelect
		ShowTypeName string `json:"show_type_name"`
	}
)

func (r *DataPapersGroupShopProperty) ParseShowType() (res string) {
	var ok bool
	if res, ok = MapDataPapersGroupShopProperty[r.ShowType]; ok {
		return
	}
	res = fmt.Sprintf("未知类型(%s)", r.ShowType)
	return
}

func (r *DataPapersGroupShopProperty) ToString() (res string, err error) {
	var bt []byte
	bt, err = json.Marshal(r)
	res = string(bt)
	return
}