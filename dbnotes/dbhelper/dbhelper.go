package dbhelper

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/jmoiron/sqlx"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

const driverName = "mysql"

var DB *sqlx.DB
var DBCassandra *gocql.Session

const (
	default_db_max_open = 32
	default_db_max_idle = 2
)

func GetDSN(ip string, port int, dbName, userName, pwd string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		userName, pwd, ip, port, dbName, "parseTime=true&charset=utf8mb4")
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

func InitCassDB(host []string, port int, keyspace, consistencyOption, userName, password string) *gocql.Session {

	cluster := gocql.NewCluster(host...)
	cluster.Keyspace = keyspace
	if len(userName) > 0 && len(password) > 0 {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: userName,
			Password: password}
	}

	ConsistencyMap := map[string]gocql.Consistency{
		"One":         gocql.One,
		"Quorum":      gocql.Quorum,
		"LocalOne":    gocql.LocalOne,
		"LocalQuorum": gocql.LocalQuorum,
	}

	cluster.Consistency = ConsistencyMap[consistencyOption]
	cluster.Timeout = time.Duration(time.Second * 2)
	cluster.Port = port
	session, err := cluster.CreateSession()
	if err != nil {

		log.Fatal("cluster.CreateSession() err:", err)
		return nil
	}

	log.Println("initCassDB finish")
	return session
}

func Init(ip string, port int, dbName, userName, pwd string) {
	DB = GetDB(ip, port, dbName, userName, pwd)

}
