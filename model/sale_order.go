package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

var Table_Sale_Order string = "sale_order"

type SaleOrder struct {
	ID            int64           `db:"id" exclude:"true"`
	OrderNo       sql.NullString  `db:"order_no"`
	CreateTime    time.Time       `db:"create_time" exclude:"true"`
	UserID        int64           `db:"user_id"`
	Status        string          `db:"status" exclude:"true"`
	Receiver      string          `db:"receiver"`
	ProvinceID    int             `db:"province_id"`
	CityID        int             `db:"city_id"`
	DistrictID    int             `db:"district_id"`
	Address       string          `db:"address"`
	PhoneNo       string          `db:"phone_no"`
	SupplierID    int64           `db:"supplier_id"`
	OrderAmt      decimal.Decimal `db:"order_amt"`
	GoodsAmt      decimal.Decimal `db:"goods_amt"`
	ExpressMethod *string         `db:"express_method"`
	ExpressNo     *string         `db:"express_order_no"`
	ExpressFee    decimal.Decimal `db:"express_fee"`
	Count         int64           `db:"count" exclude:"true"`
}
