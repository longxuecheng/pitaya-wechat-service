package model

type User struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	PhoneNo string `db:"phone_no"`
	Email   string `db:"email"`
}
