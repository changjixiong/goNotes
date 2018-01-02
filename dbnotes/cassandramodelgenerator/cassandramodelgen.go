package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goNotes/dbnotes/modeltool"

	"github.com/gocql/gocql"
)

var StatsDBCassandra *gocql.Session

func initCassDB(host string, port int, keyspace, consistencyOption, userName, password string) *gocql.Session {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Port = port
	if len(userName) > 0 && len(password) > 0 {
		cluster.Authenticator = gocql.PasswordAuthenticator{Username: userName, Password: password}
	}

	ConsistencyMap := map[string]gocql.Consistency{
		"One":         gocql.One,
		"Quorum":      gocql.Quorum,
		"LocalOne":    gocql.LocalOne,
		"LocalQuorum": gocql.LocalQuorum,
	}

	cluster.Consistency = ConsistencyMap[consistencyOption]
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("cluster.CreateSession() err:", err)
		return nil
	}

	return session
}

func genModelFile(render *template.Template, session *gocql.Session, packageName, dbConnection, keyspaceName, tableName string) {
	tableSchema := &[]modeltool.TABLE_SCHEMA{}

	sql := "select column_name, type, kind from system_schema.columns where keyspace_name= ? and table_name= ?"
	iter := session.Query(
		sql,
		keyspaceName,
		tableName,
	).Iter()

	columnName := ""
	typeName := ""
	kind := ""
	for iter.Scan(&columnName, &typeName, &kind) {
		*tableSchema = append(*tableSchema, modeltool.TABLE_SCHEMA{
			COLUMN_NAME: columnName,
			COLUMN_KEY:  kind,
			DATA_TYPE:   typeName,
		})
	}

	if len(*tableSchema) <= 0 {
		fmt.Println(tableName, "tableSchema is null")
		return
	}

	model := &modeltool.ModelInfo{
		PackageName:  packageName,
		DBConnection: dbConnection,
		TableName:    tableName,
		ModelName:    tableName,
		TableSchema:  tableSchema}

	fileName := *modelFolder + strings.ToLower(tableName) + *fileTail + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	if err := render.Execute(f, model); err != nil {

		log.Fatal("", err)
	}
	fmt.Println("fileName", fileName)
	cmd := exec.Command("goimports", "-w", fileName)
	cmd.Run()

}

func getTableNames(session *gocql.Session, keyspaceName string) (tableNames []string) {
	sql := "select table_name from system_schema.tables where keyspace_name= ?"
	iter := session.Query(
		sql,
		keyspaceName,
	).Iter()

	tableName := ""
	for iter.Scan(&tableName) {
		tableNames = append(tableNames, tableName)
	}

	return tableNames
}

var tplFile = flag.String("tplFile", "", "the path of tpl file")
var modelFolder = flag.String("modelFolder", "", "the path for folder of model files")
var packageName = flag.String("packageName", "", "packageName")
var fileTail = flag.String("fileTail", "_cassandra", "fileTail")

var dbConnection = flag.String("dbConnection", "db.DBCassandra", "the name of db instance used in model files")
var dbIP = flag.String("dbIP", "", "the ip of db host")
var dbPort = flag.Int("dbPort", 9042, "the port of db host")

var dbName = flag.String("dbName", "", "the name of db")
var userName = flag.String("userName", "", "the user name of db")
var pwd = flag.String("pwd", "", "the password of db")
var genTable = flag.String("genTable", "", "the name of table to be generated")

func main() {

	flag.Parse()

	logDir, _ := filepath.Abs(*modelFolder)
	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}

	data, err := ioutil.ReadFile(*tplFile)
	if nil != err {
		fmt.Println("read tplFile err:", err)
		return
	}

	StatsDBCassandra = initCassDB(*dbIP, *dbPort, *dbName, "Quorum", *userName, *pwd)

	if nil == StatsDBCassandra {
		fmt.Println("连接失败")
		return
	}
	render := template.Must(template.New("statscassandra").
		Funcs(template.FuncMap{
			"FirstCharUpper":       modeltool.FirstCharUpper,
			"FirstCharLower":       modeltool.FirstCharLower,
			"TypeConvert":          modeltool.TypeConvert,
			"Tags":                 modeltool.Tags,
			"ExportColumn":         modeltool.ExportColumn,
			"Join":                 modeltool.Join,
			"MakeQuestionMarkList": modeltool.MakeQuestionMarkList,
			"ColumnAndType":        modeltool.ColumnAndType,
			"ColumnWithPostfix":    modeltool.ColumnWithPostfix,
			"IsUUID":               modeltool.IsUUID,
		}).
		Parse(string(data)))

	tables := make([]string, 0)

	if len(*genTable) > 0 {
		tables = strings.Split(*genTable, "#")
	} else {
		tables = getTableNames(StatsDBCassandra, *dbName)
	}

	for _, table := range tables {
		genModelFile(render, StatsDBCassandra, *packageName, *dbConnection, *dbName, table)
	}
}
