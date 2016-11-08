package main

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
select * from tables where table_schema = 'dbnote'

SELECT
        COLUMN_NAME,DATA_TYPE, COLUMN_COMMENT,
        COLUMN_DEFAULT,COLUMN_KEY,EXTRA
        FROM COLUMNS
        WHERE TABLE_NAME = 'mail'  and TABLE_SCHEMA = 'dbnote'
*/

func genModelFile(tableName string, render *template.Template) {
	os.Remove(tableName + ".go")
	f, err := os.Create(tableName + ".go")

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &ModelInfo{PackageName: "model",
		ModelName: tableName}
	if err := render.Execute(f, model); err != nil {
		log.Fatal(err)
	}
}

func firstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

type ModelInfo struct {
	PackageName string
	ModelName   string
}

func main() {

	data, _ := ioutil.ReadFile("model.tpl")
	render := template.Must(template.New("model").
		Funcs(template.FuncMap{"firstCharUpper": firstCharUpper}).
		Parse(string(data)))

	dbName := "dbnote"

	dbhelper.GetDB("127.0.0.1", 3306, "information_schema", "root", "123456")
	var tablaNames []string
	err := dbhelper.DB.Select(&tablaNames,
		"SELECT table_name from tables where table_schema = '"+dbName+"'")
	if err != nil {
		fmt.Println(err)
	}

	for _, table := range tablaNames {
		genModelFile(table, render)
	}

}
