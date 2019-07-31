package sys

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDB(t *testing.T) {
	connectURL := fmt.Sprintf("%s:%s@tcp(%s)/mymall?allowNativePasswords=true&parseTime=true", "root", "6263272lxc", "localhost:3305")
	db, err := sqlx.Connect("mysql", connectURL)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	type u struct {
		ID   int64  `db:"id"`
		Name string `db:"nick_name"`
	}
	user := u{}
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		rows, err := db.QueryxContext(ctx, "SELECT id,nick_name FROM user")
		if err != nil {
			t.Error(err)
		}
		for rows.Next() {
			err := rows.StructScan(&user)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("> %d loop %#v\n", i, user)
		}
	}
	us := []u{}
	for i := 0; i < 100; i++ {
		err := db.Select(&us, "SELECT id,nick_name FROM user")
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("> %d loop %#v\n", i, us[0])
	}

}
