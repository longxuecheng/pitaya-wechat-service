package model

type Attribute struct {
	ID                  int64  `db:"id"`
	AttributeCategoryID int64  `db:"attribute_category_id"`
	Name                string `db:"name"`
	InputType           int    `db:"input_type"`
	Value               string `db:"value"`
	SortOrder           int    `db:"sort_order"`
}

type AttributeCategory struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	enabled int    `db:"enabled"`
}
