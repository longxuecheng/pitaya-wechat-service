package model

var Table_Specification = "specification"

type Specification struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
