package app_param

import "github.com/shopspring/decimal"

type (
	ArgOrderFromCartItem struct {
		SkuId       string `json:"sku_id" form:"sku_id"`       // sku地址
		Num         int    `json:"num" form:"num"`             // 商品数量
		SkuPrice    string `json:"sku_price" form:"sku_price"` // SPU项目本次要支付的单价(定金预售定金金额或尾款金额 sku_price)
		SkuSetPrice string `json:"sku_price" form:"sku_price"` // SPU项目本的单价
		ShopId      int64  `json:"shop_id" form:"shop_id"`     // 店铺ID
		SpuId       string `json:"spu_id" form:"spu_id"`       // 商品ID
	}
)

func (r *ArgOrderFromCartItem) GetPrice() (res decimal.Decimal, err error) {
	if res, err = decimal.NewFromString(r.SkuPrice); err != nil {
		return
	}
	return
}

func (r *ArgOrderFromCartItem) GetTotalSkuPrice() (res decimal.Decimal, err error) {
	var price decimal.Decimal
	if price, err = r.GetPrice(); err != nil {
		return
	}
	res = price.Mul(decimal.NewFromInt(int64(r.Num)))
	return
}
