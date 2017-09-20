{{$exportModelName := .ModelName | ExportColumn | FirstCharUpper}}

package {{.PackageName}}

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.COLUMN_NAME | ExportColumn}} {{.DATA_TYPE | TypeConvert}} {{.COLUMN_NAME | Tags}} // {{.COLUMN_COMMENT}}
{{end}}}

type {{$exportModelName | FirstCharLower}}Op struct{}

var {{$exportModelName}}Op = &{{$exportModelName | FirstCharLower}}Op{}
var Default{{$exportModelName}} = &{{$exportModelName}}{}

func (op *{{$exportModelName | FirstCharLower}}Op) Insert(m *{{$exportModelName}}) (int64, error) {
	return op.InsertTx({{.DBConnection}}, m)
}

func (op *{{$exportModelName | FirstCharLower}}Op) InsertTx(session *gocql.Session, m *{{$exportModelName}}) (int64, error) {
	sql := "insert into {{.TableName}}({{Join .ColumnNames ","}}) values({{.ColumnCount | MakeQuestionMarkList}})"
	if err := session.Query(
		sql,
		{{range .TableSchema}}{{if .DATA_TYPE | IsUUID}} gocql.TimeUUID() {{else}} m.{{.COLUMN_NAME | ExportColumn}} {{end}},
		{{end}}
	).Exec(); err != nil {
		return -1, err

	}

	return 0, nil
}

func (op *{{$exportModelName | FirstCharLower}}Op) QueryByMap(m map[string]interface{}, options []string) ([]*{{$exportModelName}}, error) {
	result := []*{{$exportModelName}}{}
	var params []interface{}

	sql := "select {{Join .ColumnNames ","}} from {{.TableName}}"

	kNo := 0
	for k,v := range m{
		if (kNo==0){
			sql += " where "+ k +" = ?"
		}else{
			sql += " and "+ k +" = ?"
		}

		kNo += 1

		params = append(params, v)
	}

	if len(m) >0 {
		for _, option := range options{
			sql += " " + option
		}
	} 

	iter := {{.DBConnection}}.Query(sql, params...).Iter()

	if nil == iter{
		return result, nil
	}

	data := &{{$exportModelName}}{}
	for iter.Scan(
	{{range .TableSchema}} &data.{{.COLUMN_NAME | ExportColumn}},
	{{end}}
	) {
		result = append(result, data)

		data = &{{$exportModelName}}{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}

func (op *{{$exportModelName | FirstCharLower}}Op) QueryByMapComparison(m map[string]interface{}, options []string) ([]*{{$exportModelName}}, error) {
	result := []*{{$exportModelName}}{}
	var params []interface{}

	sql := "select {{Join .ColumnNames ","}} from {{.TableName}}"

	kNo := 0
	for k,v := range m{
		if (kNo==0){
			sql += " where "+ k +" ?"
		}else{
			sql += " and "+ k +" ?"
		}

		kNo += 1

		params = append(params, v)
	}

	if len(m) >0 {
		for _, option := range options{
			sql += " " + option
		}
	} 

	iter := {{.DBConnection}}.Query(sql, params...).Iter()

	if nil == iter{
		return result, nil
	}

	data := &{{$exportModelName}}{}
	for iter.Scan(
	{{range .TableSchema}} &data.{{.COLUMN_NAME | ExportColumn}},
	{{end}}
	) {
		result = append(result, data)

		data = &{{$exportModelName}}{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}


