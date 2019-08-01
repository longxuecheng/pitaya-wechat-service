package model

type Supplier struct {
	ID      int64   `db:"id"`
	Name    *string `db:"name"`
	AdminID int64   `db:"admin_id"`
}

func (*Supplier) TableName() string {
	return "supplier"
}
