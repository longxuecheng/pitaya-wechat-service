package sys

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var dbConnection *sqlx.DB

func connectDataBase() {
	if dbConnection == nil {
		db, err := sqlx.Connect("mysql", "root:6263272lxc@tcp(localhost:3306)/mymall?allowNativePasswords=true")
		if err != nil {
			log.Fatalln(err)
		}
		err = db.Ping()
		if err != nil {
			log.Panic("ping to database maybe some problems")
		}
		dbConnection = db
	}
}

func DBConnection() *sqlx.DB {
	connectDataBase()
	return dbConnection
}
