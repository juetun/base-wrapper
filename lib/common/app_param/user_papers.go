package app_param

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
