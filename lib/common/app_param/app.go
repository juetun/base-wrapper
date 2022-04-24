package app_param

const (
	AppNameAdmin   = "admin-main"
	AppNameUpload  = "api-upload"
	AppNameExport  = "api-export"
	AppNameUser    = "api-user"
	AppNameTag     = "api-tag"
	AppNameComment = "api-comment"
	AppNameChat    = "api-chat"

	AppNameCart = "api-car"

	AppNameMall             = "api-mall"
	AppNameMallOrder        = "api-order"
	AppNameMallOrderComment = "api-ordercomment"
	AppNameMallActivity     = "api-activity"
)

//标签类型定义
const (
	DataPapersGroupCategoryTag          = "user_tag"           // 用户标签
	DataPapersGroupCategoryMallCategory = "mall_category"      // 电商类目
	DataPapersGroupCategoryMallBrand    = "mall_brand_quality" // 电商品牌类目
)

var MapDataPapersGroupCategory = map[string]string{
	DataPapersGroupCategoryTag:          "用户标签",
	DataPapersGroupCategoryMallCategory: "电商类目",
	DataPapersGroupCategoryMallBrand:    "电商品牌",
}
