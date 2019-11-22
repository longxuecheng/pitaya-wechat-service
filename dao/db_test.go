package dao

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func TestDBStats(t *testing.T) {
	for i := 0; i < 100; i++ {
		total, err := CartDao.SelectCountByUserID(0)
		fmt.Println(GetEasyDB().Stats())
		if err != nil {
			t.Error(err)
		}
		fmt.Println(total)
	}
}

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

func TestSessionID(t *testing.T) {
	connectURL := fmt.Sprintf("%s:%s@tcp(%s)/mymall?allowNativePasswords=true&parseTime=true", "root", "6263272lxc", "localhost:3305")
	db1, err := sqlx.Connect("mysql", connectURL)
	if err != nil {
		log.Fatalln(err)
	}
	db1.SetMaxIdleConns(5)
	db1.SetMaxOpenConns(5)
	// u := &model.User{}
	db2, err := sqlx.Connect("mysql", connectURL)
	if err != nil {
		log.Fatalln(err)
	}
	db2.SetMaxIdleConns(5)
	db2.SetMaxOpenConns(5)
	ctx := context.Background()
	for i := 0; i < 1000; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("%+v\n", err)
				}
			}()

			// 在多个事务并发执行
			txx, err := db1.BeginTxx(ctx, &sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
			})
			if err != nil {
				fmt.Printf("start tx err %+v\n", err)
			}
			result, err := txx.ExecContext(ctx, "update user set channel_user_id = channel_user_id + 1 where id = 11")
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
			txx.Commit()
			fmt.Printf("DB1 Affected rows number is %d\n", count)
		}()

	}

	for i := 0; i < 1000; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("%+v\n", err)
				}
			}()

			// 在多个事务并发执行
			txx, err := db2.BeginTxx(ctx, &sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
			})
			if err != nil {
				fmt.Printf("start tx err %+v\n", err)
			}
			result, err := txx.ExecContext(ctx, "update user set channel_user_id = channel_user_id + 1 where id = 11")
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
			txx.Commit()
			fmt.Printf("DB2 Affected rows number is %d\n", count)
		}()

	}

	time.Sleep(10 * time.Second)
}
