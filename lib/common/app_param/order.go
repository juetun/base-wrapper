package app_param

type (
	ArgOrderFromCartItem struct {
		SkuId  string `json:"sku_id" form:"sku_id"` //sku地址
		Num    int    `json:"num" form:"num"`       //商品数量
		ShopId int64  `json:"shop_id" form:"shop_id"`
		SpuId  string `json:"spu_id" form:"spu_id"`
	}
)
