package sys

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var dbConnection *sqlx.DB

func connectDataBase() *sqlx.DB{
	if dbConnection != nil {
		return dbConnection
	}
	db, err := sqlx.Connect("mysql", "root:6263272lxc@tcp(localhost:3306)/mymall")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic("ping to database maybe some problems")
	}
	return db
}

func DBConnection() *sqlx.DB {
	return connectDataBase()
}
