package modeltoolv0

import (
	"html/template"
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

func (m *ModelInfo) ColumnNames() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {

		result = append(result, t.COLUMN_NAME)

	}
	return result
}

func (m *ModelInfo) ColumnCount() int {
	return len(*m.TableSchema)
}

func (m *ModelInfo) PkColumnsSchema() []TABLE_SCHEMA {
	result := make([]TABLE_SCHEMA, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.COLUMN_KEY == "PRI" {
			result = append(result, t)
		}
	}
	return result
}

func (m *ModelInfo) NoPkColumnsSchema() []TABLE_SCHEMA {
	result := make([]TABLE_SCHEMA, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.COLUMN_KEY != "PRI" {
			result = append(result, t)
		}
	}
	return result
}

func (m *ModelInfo) NoPkColumns() []string {
	noPkColumnsSchema := m.NoPkColumnsSchema()
	result := make([]string, 0, len(noPkColumnsSchema))
	for _, t := range noPkColumnsSchema {
		result = append(result, t.COLUMN_NAME)
	}
	return result
}

func (m *ModelInfo) PkColumns() []string {
	pkColumnsSchema := m.PkColumnsSchema()
	result := make([]string, 0, len(pkColumnsSchema))
	for _, t := range pkColumnsSchema {
		result = append(result, t.COLUMN_NAME)
	}
	return result
}

func FirstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

func Tags(columnName string) template.HTML {

	return template.HTML("`db:" + `"` + columnName + `"` +
		" json:" + `"` + columnName + "\"`")
}

func ExportColumn(columnName string) string {
	columnItems := strings.Split(columnName, "_")
	columnItems[0] = FirstCharUpper(columnItems[0])
	for i := 0; i < len(columnItems); i++ {
		if strings.ToUpper(columnItems[i]) == "ID" {
			columnItems[i] = "ID"
		}
	}

	return strings.Join(columnItems, "")

}

func TypeConvert(str string) string {

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

func Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func ColumnAndType(table_schema []TABLE_SCHEMA) string {
	result := make([]string, 0, len(table_schema))
	for _, t := range table_schema {
		result = append(result, t.COLUMN_NAME+" "+TypeConvert(t.DATA_TYPE))
	}
	return strings.Join(result, ",")
}

func ColumnWithPostfix(columns []string, Postfix, sep string) string {
	result := make([]string, 0, len(columns))
	for _, t := range columns {
		result = append(result, t+Postfix)
	}
	return strings.Join(result, sep)
}

func MakeQuestionMarkList(num int) string {
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
