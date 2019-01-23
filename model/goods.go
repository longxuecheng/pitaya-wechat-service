package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Goods struct {
	ID               int64           `db:"id"`
	Name             string          `db:"name"`
	CategoryID       int             `db:"category_id"`
	BriefDescription sql.NullString  `db:"brief_description"`
	Description      sql.NullString  `db:"description"`
	Status           string          `db:"status"`
	IsDelete         int8            `db:"is_delete"`
	SupplierID       int64           `db:"supplier_id"`
	CreateTime       *time.Time      `db:"create_time"`
	ListPicURL       sql.NullString  `db:"list_pic_url"`
	RetailPrice      decimal.Decimal `db:"retail_price"`
}
