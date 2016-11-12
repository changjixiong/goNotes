package dbhelper

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

const driverName = "mysql"

var DB *sqlx.DB
var SYSDB *sqlx.DB

const (
	default_db_max_open    = 32
	default_db_max_idle    = 2
	default_redis_max_open = 1
)

func GetDSN(ip string, port int, dbName, userName, pwd string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		userName, pwd, ip, port, dbName, "parseTime=true")
}

func GetDB(ip string, port int, dbName, userName, pwd string) *sqlx.DB {
	return initSqlxDB(GetDSN(ip, port, dbName, userName, pwd),
		default_db_max_open, default_db_max_idle)
}

func initSqlxDB(dbDSN string, maxOpen, maxIdle int) *sqlx.DB {
	db := sqlx.MustConnect(driverName, dbDSN)
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	return db
}

func init() {

	DB = GetDB("127.0.0.1", 3306, "dbnote", "root", "123456")
	SYSDB = GetDB("127.0.0.1", 3306, "information_schema", "root", "123456")
}
