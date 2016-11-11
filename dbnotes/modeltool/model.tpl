{{$exportModelName := .ModelName | firstCharUpper}}

package {{.PackageName}}

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.COLUMN_NAME | exportColumn}} {{.DATA_TYPE | typeConvert}} {{.COLUMN_NAME | tags}} // {{.COLUMN_COMMENT}}
{{end}}}

var Default{{$exportModelName}} = &{{$exportModelName}}{}


func (m *{{$exportModelName}}) GetByPK({{.PkColumnsSchema | pkWithType}}) (*{{$exportModelName}}, bool) {
	obj := &{{$exportModelName}}{}
	sql := "select * from {{.BDName}}.{{.TableName}} where {{pkWithPostfix .PkColumnsSchema "=?" " and "}}"
	err := dbhelper.DB.Get(obj, sql,
		{{range $K:=.PkColumns}}{{$K}},
		{{end}}
	)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return obj, true
}

func (m *{{$exportModelName}}) Insert() (int64, error) {
	return m.InsertTx(dbhelper.DB)
}

func (m *{{$exportModelName}}) InsertTx(ext sqlx.Ext) (int64, error) {
	sql := "insert into {{.BDName}}.{{.TableName}}({{join .ColumnNames ","}}) values({{.ColumnCount | makeQuestionMarkList}})"
	result, err := ext.Exec(sql,
		{{range .TableSchema}}m.{{.COLUMN_NAME | exportColumn}},
		{{end}}
	)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (m *{{$exportModelName}}) Update() error {
	return m.UpdateTx(dbhelper.DB)
}

func (m *{{$exportModelName}}) UpdateTx(ext sqlx.Ext) error {
	sql := `update {{.BDName}}.{{.TableName}} set {{pkWithPostfix .NoPkColumnsSchema "=?" ","}} where {{pkWithPostfix .PkColumnsSchema "=?" " and "}}`
	_, err := ext.Exec(sql,
		{{range .NoPkColumns}}m.{{. | exportColumn}},
		{{end}}{{range .PkColumns}}m.{{. | exportColumn}},
		{{end}}
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *{{$exportModelName}}) QueryByMap(ma map[string]interface{}) ([]*{{$exportModelName}}, error) {
	result := []*{{$exportModelName}}{}
	var params []interface{}

	sql := "select * from {{.BDName}}.{{.TableName}} where 1=1 "
	for k, v := range ma {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := dbhelper.DB.Select(&result, sql, params...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}