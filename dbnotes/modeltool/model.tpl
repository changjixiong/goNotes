package {{.PackageName}}

type {{.ModelName | firstCharUpper}} struct {
{{range .TableSchema}} {{.COLUMN_NAME | exportColumn}} {{.DATA_TYPE | typeConvert}} {{.COLUMN_NAME | tags}} // {{.COLUMN_COMMENT}}
{{end}}}

