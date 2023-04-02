package models

type PromoCode struct {
	PromoCodeId     int     `json:"promoCode_id"`
	PromoCodeName   string  `json:"promoCode_name"`
	Discount        float64 `json:"discount"`
	DiscountType    string  `json:"discount_type"`
	OrderLimitPrice float64 `json:"order_limit_price"`
}
type PromoCodePrimaryKey struct {
	PromoCodeId int `json:"promoCode_id"`
}

type CreatePromoCode struct {
	PromoCodeName   string  `json:"promoCode_name"`
	Discount        float64 `json:"discount"`
	DiscountType    string  `json:"discount_type"`
	OrderLimitPrice float64 `json:"order_limit_price"`
}

type GetListPromoCodeRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListPromoCodeResponse struct {
	Count      int          `json:"count"`
	PromoCodes []*PromoCode `json:"promoCodes"`
}
