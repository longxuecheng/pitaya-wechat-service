package dto

type AttributeDTO struct {
	ID    int64  `json:"-"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
