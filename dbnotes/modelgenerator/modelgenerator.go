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
)

func genModelFile(render *template.Template, dbName, tableName string) {
	tableSchema := &[]modeltoolv0.TABLE_SCHEMA{}
	err := dbhelper.DB.Select(tableSchema,
		"SELECT COLUMN_NAME, DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT from COLUMNS where "+
			"TABLE_NAME"+"='"+tableName+"' and "+"table_schema = '"+dbName+"'")

	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := "../model/" + tableName + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &modeltoolv0.ModelInfo{
		PackageName: "model",
		BDName:      dbName,
		TableName:   tableName,
		ModelName:   tableName,
		TableSchema: tableSchema}

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
			"FirstCharUpper":       modeltoolv0.FirstCharUpper,
			"TypeConvert":          modeltoolv0.TypeConvert,
			"Tags":                 modeltoolv0.Tags,
			"ExportColumn":         modeltoolv0.ExportColumn,
			"Join":                 modeltoolv0.Join,
			"MakeQuestionMarkList": modeltoolv0.MakeQuestionMarkList,
			"ColumnAndType":        modeltoolv0.ColumnAndType,
			"ColumnWithPostfix":    modeltoolv0.ColumnWithPostfix,
		}).
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
		genModelFile(render, dbName, table)
	}

	//genModelFile(render, "information_schema", "COLUMNS")

}
