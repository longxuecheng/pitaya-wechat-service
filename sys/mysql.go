package sys

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

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var easyDB *EasyDB
var dbConnection *sqlx.DB

const dsnTmplate = "%s:%s@tcp(%s)/mymall?charset=utf8mb4&allowNativePasswords=true&parseTime=true&loc=%s"

// EasyDB is a
type EasyDB struct {
	connection *sqlx.DB
	ctx        context.Context
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
		ctx:        context.Background(),
		connection: DBConnection(),
	}
	return easyDB
}

func connect() {
	if dbConnection == nil {
		var dbHost string
		env := os.Getenv("ENV")
		if env == settings.EnvProd {
			dbHost = settings.ProdDBHost
		} else {
			dbHost = settings.DevDBHost
		}
		log.Println(fmt.Sprintf("connecting to database %s", dbHost))
		connectURL := fmt.Sprintf(dsnTmplate, "root", "6263272lxc", dbHost, url.QueryEscape("Asia/Shanghai"))
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
		dbConnection = db
	}
}

func DBConnection() *sqlx.DB {
	connect()
	return dbConnection
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
func (db *EasyDB) SelectDSL(destptr interface{}, columns []string, tableName string, pred interface{}, orderBys ...string) error {
	sql, args, err := sq.Select(columns...).From(tableName).Where(pred).OrderBy(orderBys...).ToSql()
	if err != nil {
		return err
	}
	// fmt.Println(sql)
	return db.Select(destptr, sql, args...)
}

type PaginationCondition struct {
	Columns   []string
	TableName string
	Offset    uint64
	Limit     uint64
	WherePred interface{}
}

func (db *EasyDB) SelectPagination(destptr interface{}, condition PaginationCondition) (int64, error) {
	// columns = append(columns, "count(1) over() as count")
	qSQL, queryArgs, err := sq.Select(condition.Columns...).From(condition.TableName).Where(condition.WherePred).Offset(condition.Offset).Limit(condition.Limit).ToSql()
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

// SelectOne 是对sqlx包中的查询单个的简化
func (db *EasyDB) SelectOne(target interface{}, query string, args ...interface{}) error {
	rows, err := db.connection.QueryxContext(db.ctx, query, args...)
	if err != nil {
		return err
	}
	if rows.Next() {
		err = rows.StructScan(target)
		if err != nil {
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	defer rows.Close()
	return nil
}

func (db *EasyDB) Insert(tableName string, setMap map[string]interface{}, tx *sql.Tx) (rowsAffected, lastInsertID int64, err error) {
	query, args, err := sq.Insert(tableName).SetMap(setMap).ToSql()
	if err != nil {
		return
	}
	var result sql.Result
	// if tx == nil || len(tx) == 0 {
	// 	result, err = db.connection.ExecContext(db.ctx, query, args...)
	// } else {
	// 	result, err = tx[0].ExecContext(db.ctx, query, args...)
	// }
	if tx == nil {
		result, err = db.connection.ExecContext(db.ctx, query, args...)
	} else {
		result, err = tx.ExecContext(db.ctx, query, args...)
	}
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

func (db *EasyDB) UpdateTx(tx *sql.Tx, tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (rowsAffected int64, err error) {
	query, args, err := db.buildUpdtaSQL(tableName, setMap, pred, args...)
	if err != nil {
		return
	}
	var result sql.Result
	if tx != nil {
		result, err = tx.ExecContext(db.ctx, query, args...)
		if err != nil {
			return
		}
		rowsAffected, err = result.RowsAffected()
	} else {
		rowsAffected, err = db.Update(tableName, setMap, pred, args)
	}
	return
}

// Update 更新操作
func (db *EasyDB) Update(tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (rowsAffected int64, err error) {
	sql, args, err := db.buildUpdtaSQL(tableName, setMap, pred, args)
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

func (db *EasyDB) buildUpdtaSQL(tableName string, setMap map[string]interface{}, pred interface{}, args ...interface{}) (sql string, args1 []interface{}, err error) {
	sql, args1, err = sq.Update(tableName).SetMap(setMap).Where(pred, args...).ToSql()
	return
}
