{{$exportModelName := .ModelName | FirstCharUpper}}

package {{.PackageName}}

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.COLUMN_NAME | ExportColumn}} {{.DATA_TYPE | TypeConvert}} {{.COLUMN_NAME | Tags}} // {{.COLUMN_COMMENT}}
{{end}}}

var Default{{$exportModelName}} = &{{$exportModelName}}{}

{{if .HavePk}}
func (m *{{$exportModelName}}) GetByPK({{.PkColumnsSchema | ColumnAndType}}) (*{{$exportModelName}}, bool) {
	obj := &{{$exportModelName}}{}
	sql := "select * from {{.TableName}} where {{ColumnWithPostfix .PkColumns "=?" " and "}}"
	err := {{.DBConnection}}.Get(obj, sql,
		{{range $K:=.PkColumns}}{{$K}},
		{{end}}
	)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return obj, true
}
{{end}}

func (m *{{$exportModelName}}) Insert() (int64, error) {
	return m.InsertTx({{.DBConnection}})
}

func (m *{{$exportModelName}}) InsertTx(ext sqlx.Ext) (int64, error) {
	sql := "insert into {{.TableName}}({{Join .ColumnNames ","}}) values({{.ColumnCount | MakeQuestionMarkList}})"
	result, err := ext.Exec(sql,
		{{range .TableSchema}}m.{{.COLUMN_NAME | ExportColumn}},
		{{end}}
	)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

{{if .HavePk}}
func (m *{{$exportModelName}}) Delete() error {
	return m.DeleteTx({{.DBConnection}})
}

func (m *{{$exportModelName}}) DeleteTx(ext sqlx.Ext) error {
	sql := `delete from {{.TableName}} where {{ColumnWithPostfix .PkColumns "=?" " and "}}`
	_, err := ext.Exec(sql,
		{{range .PkColumns}}m.{{. | ExportColumn}},
		{{end}}
	)
	return err
}
{{end}}

{{if .HavePk}}
func (m *{{$exportModelName}}) Update() error {
	return m.UpdateTx({{.DBConnection}})
}

func (m *{{$exportModelName}}) UpdateTx(ext sqlx.Ext) error {
	sql := `update {{.TableName}} set {{ColumnWithPostfix .NoPkColumns "=?" ","}} where {{ColumnWithPostfix .PkColumns "=?" " and "}}`
	_, err := ext.Exec(sql,
		{{range .NoPkColumns}}m.{{. | ExportColumn}},
		{{end}}{{range .PkColumns}}m.{{. | ExportColumn}},
		{{end}}
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
{{end}}

func (m *{{$exportModelName}}) QueryByMap(ma map[string]interface{}) ([]*{{$exportModelName}}, error) {
	result := []*{{$exportModelName}}{}
	var params []interface{}

	sql := "select * from {{.TableName}} where 1=1 "
	for k, v := range ma {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := {{.DBConnection}}.Select(&result, sql, params...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}