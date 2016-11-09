package main

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
type ModelInfo struct {
	PackageName string
	ModelName   string
	TableSchema *[]TABLE_SCHEMA
}

type TABLE_SCHEMA struct {
	COLUMN_NAME    string `db:"COLUMN_NAME" json:"column_name"`
	DATA_TYPE      string `db:"DATA_TYPE" json:"data_type"`
	COLUMN_COMMENT string `db:"COLUMN_COMMENT" json:"COLUMN_COMMENT"`
}

func genModelFile(render *template.Template, dbName, tableName string) {
	tableSchema := &[]TABLE_SCHEMA{}
	err := dbhelper.DB.Select(tableSchema,
		"SELECT COLUMN_NAME, DATA_TYPE,COLUMN_COMMENT from COLUMNS where "+
			"TABLE_NAME"+"='"+tableName+"' and "+"table_schema = '"+dbName+"'")
	//fmt.Println(tableSchema)

	if err != nil {
		fmt.Println(err)
		return
	}

	//return

	fileName := tableName + ".go"

	os.Remove(fileName)
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &ModelInfo{PackageName: "model",
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

func firstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

func tags(columnName string) template.HTML {

	return template.HTML("`db:" + `"` + columnName + `"` +
		" json:" + `"` + columnName + "\"`")
}

func exportColumn(columnName string) string {
	columnItems := strings.Split(columnName, "_")
	columnItems[0] = firstCharUpper(columnItems[0])
	for i := 0; i < len(columnItems); i++ {
		if strings.ToUpper(columnItems[i]) == "ID" {
			columnItems[i] = "ID"
		}
	}

	return strings.Join(columnItems, "")

}

func typeConvert(str string) string {

	switch str {
	case "smallint", "tinyint":
		return "int8"

	case "varchar", "text", "longtext", "char":
		return "string"

	case "date":
		return "string"

	case "int":
		return "int"

	case "timestamp":
		return "*time.Time"

	case "bigint":
		return "int64"

	case "float", "double", "decimal":
		return "float64"

	default:
		return str
	}
}

func main() {

	data, _ := ioutil.ReadFile("model.tpl")
	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"firstCharUpper": firstCharUpper,
			"typeConvert":    typeConvert,
			"tags":           tags,
			"exportColumn":   exportColumn}).
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
