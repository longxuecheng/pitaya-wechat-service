package sys

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var easyDB *EasyDB
var dbConnection *sqlx.DB

// GetEasyDB is a method for getting a outer layer DB
func GetEasyDB() *EasyDB {
	if easyDB != nil {
		return easyDB
	}
	easyDB = new(EasyDB)
	easyDB.ctx = context.Background()
	easyDB.connection = DBConnection()
	return easyDB
}

func connectDataBase() {
	if dbConnection == nil {
		db, err := sqlx.Connect("mysql", "root:6263272lxc@tcp(localhost:3306)/mymall?allowNativePasswords=true&parseTime=true")
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

// EasyDB is a
type EasyDB struct {
	connection *sqlx.DB
	ctx        context.Context
}

// Select 只是简单套用为了能在外部用EasyDB直接执行，用的是Native SQL
func (db *EasyDB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.connection.Select(dest, query, args...)
}

// SelectOneDSL 只是简单套用为了能在外部用EasyDB直接执行，用的是Native SQL
func (db *EasyDB) SelectOneDSL(destptr interface{}, columns []string, tableName string, pred interface{}) error {
	sql, args, err := sq.Select(columns...).From(tableName).Where(pred).ToSql()
	if err != nil {
		return err
	}
	return db.SelectOne(destptr, sql, args...)
}

// SelectDSL 使用ORM的DSL进行SQL语句的初始化
func (db *EasyDB) SelectDSL(destptr interface{}, columns []string, tableName string, pred interface{}) error {
	sql, args, err := sq.Select(columns...).From(tableName).Where(pred).ToSql()
	if err != nil {
		return err
	}
	return db.Select(destptr, sql, args...)
}

// SelectOne 是对sqlx包中的查询单个的简化
func (db *EasyDB) SelectOne(target interface{}, query string, args ...interface{}) error {
	rows, err := db.connection.QueryxContext(db.ctx, query, args...)
	if err != nil {
		return err
	}
	if rows.Next() {
		err = rows.StructScan(target)
	}
	return nil
}

func (db *EasyDB) Insert(tableName string, setMap map[string]interface{}) (rowsAffected, lastInsertID int64, err error) {
	sql, args, err := sq.Insert(tableName).SetMap(setMap).ToSql()
	if err != nil {
		return
	}
	result, err := db.connection.ExecContext(db.ctx, sql, args...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	lastInsertID, err = result.LastInsertId()
	if err != nil {
		return
	}
	return
}

// Update 更新操作
func (db *EasyDB) Update(tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (rowsAffected int64, err error) {
	sql, args, err := sq.Update(tableName).SetMap(setMap).Where(pred, args...).ToSql()
	if err != nil {
		return
	}
	result, err := db.connection.ExecContext(db.ctx, sql, args...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}
