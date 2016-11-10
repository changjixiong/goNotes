{{$modelName := .ModelName}}

package {{.PackageName}}

type {{.ModelName | firstCharUpper}} struct {
{{range .TableSchema}} {{.COLUMN_NAME | exportColumn}} {{.DATA_TYPE | typeConvert}} {{.COLUMN_NAME | tags}} // {{.COLUMN_COMMENT}}
{{end}}}

var Default{{.ModelName | firstCharUpper}} = &{{.ModelName | firstCharUpper}}{}

func (m *{{.ModelName | firstCharUpper}}) GetByPK(id int) (*{{.ModelName | firstCharUpper}}, bool) {
	obj := &{{.ModelName | firstCharUpper}}{}
	sql := "select * from {{.BDName}}.{{.TableName}} where id=? "
	err := dbhelper.DB.Get(obj, sql,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return obj, true
}

func (m *{{.ModelName | firstCharUpper}}) Insert({{$modelName}} *{{.ModelName | firstCharUpper}}) (int64, error) {
	return m.InsertTx(dbhelper.DB, {{$modelName}})
}

func (m *{{.ModelName | firstCharUpper}}) InsertTx(ext sqlx.Ext, {{$modelName}} *{{.ModelName | firstCharUpper}}) (int64, error) {
	sql := "insert into {{.BDName}}.{{.TableName}}({{.TableSchema | joinByComma}}) values({{.TableSchema | joinQuestionMarkByComma}})"
	result, err := ext.Exec(sql,
		{{range .TableSchema}}{{$modelName}}.{{.COLUMN_NAME | exportColumn}},
		{{end}}
	)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (m *{{.ModelName | firstCharUpper}}) QueryByMap(ma map[string]interface{}) ([]*{{.ModelName | firstCharUpper}}, error) {
	result := []*{{.ModelName | firstCharUpper}}{}
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