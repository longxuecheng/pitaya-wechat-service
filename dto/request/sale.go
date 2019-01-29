package request

type SaleOrderAddRequest struct {
	UserID    int64 `json:"-"`
	AddressID int64 `json:"addressId"`
	CouponID  int64 `json:"couponId"`
}
