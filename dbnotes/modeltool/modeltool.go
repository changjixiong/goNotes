package main

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

/*
select * from tables where table_schema = 'dbnote'

SELECT
        COLUMN_NAME,DATA_TYPE, COLUMN_COMMENT,
        COLUMN_DEFAULT,COLUMN_KEY,EXTRA
        FROM COLUMNS
        WHERE TABLE_NAME = 'mail'  and TABLE_SCHEMA = 'dbnote'
*/
type ModelInfo struct {
	BDName      string
	TableName   string
	PackageName string
	ModelName   string
	TableSchema *[]TABLE_SCHEMA
}

type TABLE_SCHEMA struct {
	COLUMN_NAME    string `db:"COLUMN_NAME" json:"column_name"`
	DATA_TYPE      string `db:"DATA_TYPE" json:"data_type"`
	COLUMN_KEY     string `db:"COLUMN_KEY" json:"column_key"`
	COLUMN_COMMENT string `db:"COLUMN_COMMENT" json:"COLUMN_COMMENT"`
}

func (m *ModelInfo) ColumnWithModelName() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		result = append(result, m.ModelName+"."+exportColumn(t.COLUMN_NAME))
	}

	return result
}

func (m *ModelInfo) ColumnNames() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {

		result = append(result, t.COLUMN_NAME)

	}
	return result
}

func (m *ModelInfo) PkWithType() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.COLUMN_KEY == "PRI" {
			result = append(result, t.COLUMN_NAME+" "+typeConvert(t.DATA_TYPE))
		}
	}
	return result
}

func genModelFile(render *template.Template, dbName, tableName string) {
	tableSchema := &[]TABLE_SCHEMA{}
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

	model := &ModelInfo{
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

	data, _ := ioutil.ReadFile("model.tpl")
	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"firstCharUpper":          firstCharUpper,
			"typeConvert":             typeConvert,
			"tags":                    tags,
			"exportColumn":            exportColumn,
			"joinByComma":             joinByComma,
			"joinQuestionMarkByComma": joinQuestionMarkByComma}).
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
