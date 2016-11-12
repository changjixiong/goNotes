package main

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"goNotes/dbnotes/modeltool"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func genModelFile(render *template.Template, dbName, dbConnection, tableName string) {
	tableSchema := &[]modeltool.TABLE_SCHEMA{}
	err := dbhelper.SYSDB.Select(tableSchema,
		"SELECT COLUMN_NAME, DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT from COLUMNS where "+
			"TABLE_NAME"+"='"+tableName+"' and "+"table_schema = '"+dbName+"'")

	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := "../model/" + strings.ToLower(tableName) + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &modeltool.ModelInfo{
		PackageName:  "model",
		BDName:       dbName,
		DBConnection: dbConnection,
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

func main() {

	data, _ := ioutil.ReadFile("../modeltool/model.tpl")
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

	dbName := "dbnote"

	var tablaNames []string
	err := dbhelper.SYSDB.Select(&tablaNames,
		"SELECT table_name from tables where table_schema = '"+dbName+"'")
	if err != nil {
		fmt.Println(err)
	}

	for _, table := range tablaNames {
		genModelFile(render, dbName, "DB", table)
	}

}
