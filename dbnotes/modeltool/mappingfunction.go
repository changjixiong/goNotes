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

func joinByComma(tableSchema *[]TABLE_SCHEMA) string {
	columns := make([]string, 0, len(*tableSchema))
	for _, t := range *tableSchema {
		columns = append(columns, t.COLUMN_NAME)
	}

	return strings.Join(columns, ",")
}

func joinByComma2(src []string) string {

	return strings.Join(src, ",\n")
}

func joinQuestionMarkByComma(tableSchema *[]TABLE_SCHEMA) string {
	columns := make([]string, 0, len(*tableSchema))
	for _, _ = range *tableSchema {
		columns = append(columns, "?")
	}

	return strings.Join(columns, ",")
}
