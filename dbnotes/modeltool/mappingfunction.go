package main

import (
	"html/template"
	"strings"
)

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

func join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func columnAndType(table_schema []TABLE_SCHEMA) string {
	result := make([]string, 0, len(table_schema))
	for _, t := range table_schema {
		result = append(result, t.COLUMN_NAME+" "+typeConvert(t.DATA_TYPE))
	}
	return strings.Join(result, ",")
}

func columnWithPostfix(columns []string, Postfix, sep string) string {
	result := make([]string, 0, len(columns))
	for _, t := range columns {
		result = append(result, t+Postfix)
	}
	return strings.Join(result, sep)
}

func makeQuestionMarkList(num int) string {
	a := strings.Repeat("?,", num)
	return a[:len(a)-1]
}

/*
func joinQuestionMarkByComma(tableSchema *[]TABLE_SCHEMA) string {
	columns := make([]string, 0, len(*tableSchema))
	for _, _ = range *tableSchema {
		columns = append(columns, "?")
	}

	return strings.Join(columns, ",")
}
*/
