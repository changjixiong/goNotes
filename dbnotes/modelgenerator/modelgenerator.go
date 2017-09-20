package main

import (
	"flag"
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"goNotes/dbnotes/modeltool"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

func genModelFile(db *sqlx.DB, render *template.Template, packageName, tableName string) {
	tableSchema := &[]modeltool.TABLE_SCHEMA{}
	err := db.Select(tableSchema,
		"SELECT COLUMN_NAME, DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT from COLUMNS where "+
			"TABLE_NAME"+"='"+tableName+"' and "+"table_schema = '"+*dbName+"'")

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(*tableSchema) <= 0 {
		fmt.Println(tableName, "tableSchema is null")
		return
	}

	fileName := *modelFolder + strings.ToLower(tableName) + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &modeltool.ModelInfo{
		PackageName:  packageName,
		BDName:       *dbName,
		DBConnection: *dbConnection,
		TableName:    tableName,
		ModelName:    tableName,
		TableSchema:  tableSchema}

	if err := render.Execute(f, model); err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileName)
	cmd := exec.Command("goimports", "-w", fileName)
	//cmd := exec.Command("gofmt", "-w", fileName)
	cmd.Run()
}

var tplFile = flag.String("tplFile", "./model.tpl", "the path of tpl file")
var modelFolder = flag.String("modelFolder", "../model/", "the path for folder of model files")
var genTable = flag.String("genTable", "", "the name of table to be generated")
var dbInstanceName = flag.String("dbInstanceName", "dbhelper.DB", "the name of db instance used in model files")
var dbConnection = flag.String("dbConnection", "", "the name of db connection instance used in model files")
var packageName = flag.String("packageName", "", "packageName")
var dbIP = flag.String("dbIP", "127.0.0.1", "the ip of db host")
var dbPort = flag.Int("dbPort", 3306, "the port of db host")
var dbName = flag.String("dbName", "dbnote", "the name of db")
var userName = flag.String("userName", "root", "the user name of db")
var pwd = flag.String("pwd", "123456", "the password of db")

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

	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"FirstCharUpper":       modeltool.FirstCharUpper,
			"TypeConvert":          modeltool.TypeConvert,
			"Tags":                 modeltool.Tags,
			"ExportColumn":         modeltool.ExportColumn,
			"Join":                 modeltool.Join,
			"MakeQuestionMarkList": modeltool.MakeQuestionMarkList,
			"ColumnAndType":        modeltool.ColumnAndType,
			"ColumnWithPostfix":    modeltool.ColumnWithPostfix,
		}).
		Parse(string(data)))

	var tablaNames []string
	sysDB := dbhelper.GetDB(*dbIP, *dbPort, "information_schema", *userName, *pwd)

	if len(*genTable) > 0 {
		tablaNames = strings.Split(*genTable, "#")
	} else {
		err = sysDB.Select(&tablaNames,
			"SELECT table_name from tables where table_schema = '"+*dbName+"'")
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, table := range tablaNames {
		genModelFile(sysDB, render, *packageName, table)
	}

}
