package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gotrue/facility/utils"
	"gotrue/settings"
	"log"
	"net/url"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var easyDB *EasyDB

const dsnTemplate = "%s:%s@tcp(%s)/mymall?charset=utf8mb4&allowNativePasswords=true&parseTime=true&loc=%s"

// EasyDB is a
type EasyDB struct {
	connection *sqlx.DB
}

type count struct {
	Count int64 `db:"count"`
}

// GetEasyDB is a method for getting a outer layer DB
func GetEasyDB() *EasyDB {
	if easyDB != nil {
		return easyDB
	}
	easyDB = &EasyDB{
		connection: NewConn(),
	}
	return easyDB
}

func NewConn() *sqlx.DB {
	var dbHost string
	env := os.Getenv("ENV")
	if env == settings.EnvProd {
		dbHost = settings.ProdDBHost
	} else {
		dbHost = settings.DevDBHost
	}
	log.Printf("connecting to database %s \n", dbHost)
	connectURL := fmt.Sprintf(dsnTemplate, "root", "6263272lxc", dbHost, url.QueryEscape("Asia/Shanghai"))
	db, err := sqlx.Connect("mysql", connectURL)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	err = db.Ping()
	if err != nil {
		log.Panic("ping to database maybe some problems")
	}
	return db
}

func (db *EasyDB) NewTxExecutor() (*TxExecutor, error) {
	tx, err := db.connection.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &TxExecutor{
		tx:    tx,
		chain: []TxFunc{},
	}, nil
}

func (db *EasyDB) Stats() string {
	stats := db.connection.Stats()
	return fmt.Sprintf("Idle %d InUse %d Open %d wait %d", stats.Idle, stats.InUse, stats.OpenConnections, stats.WaitCount)
}

// ExecTx exec a pipline operation in transation
func (db *EasyDB) ExecTx(execFunc func(tx *sql.Tx) error) {
	tx, err := db.connection.Begin()
	utils.CheckAndPanic(err)
	err = execFunc(tx)
	defer func() {
		if err == nil {
			err = tx.Commit()
			utils.CheckAndPanic(err)
		} else {
			if tx == nil {
				panic(errors.New("transaction not exists"))
			} else {
				tx.Rollback()
				panic(err)
			}
		}
	}()
}

func (db *EasyDB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.connection.Select(dest, query, args...)
}

func (db *EasyDB) SelectOneDSL(destptr interface{}, columns []string, tableName string, pred interface{}) error {
	sql, args, err := sq.Select(columns...).From(tableName).Where(pred).ToSql()
	if err != nil {
		return err
	}
	return db.SelectOne(destptr, sql, args...)
}

func (db *EasyDB) SelectDSL(destptr interface{}, columns []string, tableName string, pred interface{}, orderBys ...string) error {
	sql, args, err := sq.Select(columns...).From(tableName).Where(pred).OrderBy(orderBys...).ToSql()
	if err != nil {
		return err
	}
	return db.Select(destptr, sql, args...)
}

type PaginationCondition struct {
	Columns     []string
	OrderbyPred string
	TableName   string
	Offset      uint64
	Limit       uint64
	WherePred   interface{}
}

func (db *EasyDB) SelectPagination(destptr interface{}, condition PaginationCondition) (int64, error) {
	qSQL, queryArgs, err := sq.Select(condition.Columns...).From(condition.TableName).
		Where(condition.WherePred).
		OrderBy(condition.OrderbyPred).
		Offset(condition.Offset).
		Limit(condition.Limit).ToSql()
	if err != nil {
		return 0, err
	}
	countSQL, countArgs, err := sq.Select("count(1) as count").From(condition.TableName).Where(condition.WherePred).ToSql()
	count := new(count)
	err = db.SelectOne(count, countSQL, countArgs...)
	if err != nil {
		return 0, err
	}
	return count.Count, db.Select(destptr, qSQL, queryArgs...)
}

func (db *EasyDB) SelectOne(target interface{}, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	e := make(chan error)
	go func() {
		rows, err := db.connection.QueryxContext(ctx, query, args...)
		if err != nil {
			e <- err
			return
		}
		defer rows.Close()
		if err != nil {
			e <- err
			return
		}
		if rows.Next() {
			err = rows.StructScan(target)
			if err != nil {
				e <- err
			}
		} else {
			e <- sql.ErrNoRows
		}
		e <- nil
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-e:
		return err
	}
}

func (db *EasyDB) Insert(tableName string, setMap map[string]interface{}, tx *sql.Tx) (rowsAffected, lastInsertID int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query, args, err := sq.Insert(tableName).SetMap(setMap).ToSql()
	if err != nil {
		return
	}
	e := make(chan error)
	var result sql.Result
	go func() {
		if tx == nil {
			result, err = db.connection.ExecContext(ctx, query, args...)
		} else {
			result, err = tx.ExecContext(ctx, query, args...)
		}
		if err != nil {
			e <- err
			return
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			e <- err
			return
		}
		lastInsertID, err = result.LastInsertId()
		e <- err
	}()
	select {
	case <-ctx.Done():
		return rowsAffected, lastInsertID, ctx.Err()
	case err := <-e:
		return rowsAffected, lastInsertID, err
	}
}

var ErrorNowRowAffected = errors.New("No rows affected")

func (db *EasyDB) ExecuteSQLTx(tx *sql.Tx, query string, args ...interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	txCopy := tx
	if txCopy == nil {
		newTx, err := db.connection.Begin()
		if err != nil {
			return err
		}
		txCopy = newTx
	}
	var err error
	defer func() {
		// if tx is nil,commit as a separated transaction
		// if tx is not nil, the control right will give to existed tx
		if tx == nil {
			if err != nil {
				txCopy.Rollback()
			}
			txCopy.Commit()
		}

	}()
	result, err := txCopy.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrorNowRowAffected
	}
	return nil
}

func (db *EasyDB) UpdateTx(tx *sql.Tx, tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query, args, err := db.buildUpdateSQL(tableName, setMap, pred, args...)
	if err != nil {
		return
	}
	e := make(chan error)
	var affectedRows int64
	go func() {
		var result sql.Result
		var err error
		if tx != nil {
			result, err = tx.ExecContext(ctx, query, args...)
			if err != nil {
				e <- err
				return
			}
			affectedRows, err = result.RowsAffected()
		} else {
			affectedRows, err = db.Update(tableName, setMap, pred, args)
		}
		e <- err
	}()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err := <-e:
		return affectedRows, err
	}
}

func (db *EasyDB) Update(tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (affectedRows int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	sql, args, err := db.buildUpdateSQL(tableName, setMap, pred, args)
	if err != nil {
		return
	}
	e := make(chan error)
	go func() {
		result, err := db.connection.ExecContext(ctx, sql, args...)
		if err != nil {
			e <- err
		}
		affectedRows, err = result.RowsAffected()
		e <- err
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err := <-e:
		return affectedRows, err
	}
}

func (db *EasyDB) buildUpdateSQL(tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (sql string, args1 []interface{}, err error) {
	sql, args1, err = sq.Update(tableName).SetMap(setMap).Where(pred, args...).ToSql()
	return
}
