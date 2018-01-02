package cassandra

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"time"

	"github.com/gocql/gocql"
)

type NumLog4s struct {
	CreateTime time.Time  `db:"create_time" json:"create_time"` //
	ID         gocql.UUID `db:"id" json:"id"`                   //
	Num        int        `db:"num" json:"num"`                 //
	ServerID   int        `db:"server_id" json:"server_id"`     //
}

type numLog4sOp struct{}

var NumLog4sOp = &numLog4sOp{}
var DefaultNumLog4s = &NumLog4s{}

func (op *numLog4sOp) Insert(m *NumLog4s) (int64, error) {
	return op.InsertTx(dbhelper.DBCassandra, m)
}

func (op *numLog4sOp) InsertTx(session *gocql.Session, m *NumLog4s) (int64, error) {
	sql := "insert into num_log_4s(create_time,id,num,server_id) values(?,?,?,?)"
	if err := session.Query(
		sql,
		m.CreateTime,
		gocql.TimeUUID(),
		m.Num,
		m.ServerID,
	).Exec(); err != nil {
		return -1, err

	}

	return 0, nil
}

func (op *numLog4sOp) QueryByMap(m map[string]interface{}, options []string) ([]*NumLog4s, error) {
	result := []*NumLog4s{}
	var params []interface{}

	sql := "select create_time,id,num,server_id from num_log_4s"

	kNo := 0
	for k, v := range m {
		if kNo == 0 {
			sql += " where " + k + " = ?"
		} else {
			sql += " and " + k + " = ?"
		}

		kNo += 1

		params = append(params, v)
	}

	if len(m) > 0 {
		for _, option := range options {
			sql += " " + option
		}
	}

	iter := dbhelper.DBCassandra.Query(sql, params...).Iter()

	if nil == iter {
		return result, nil
	}

	data := &NumLog4s{}
	for iter.Scan(
		&data.CreateTime,
		&data.ID,
		&data.Num,
		&data.ServerID,
	) {
		result = append(result, data)

		data = &NumLog4s{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}

func (op *numLog4sOp) QueryByMapComparison(m map[string]interface{}, options []string) ([]*NumLog4s, error) {
	result := []*NumLog4s{}
	var params []interface{}

	sql := "select create_time,id,num,server_id from num_log_4s"

	kNo := 0
	for k, v := range m {
		if kNo == 0 {
			sql += " where " + k + " ?"
		} else {
			sql += " and " + k + " ?"
		}

		kNo += 1

		params = append(params, v)
	}

	if len(m) > 0 {
		for _, option := range options {
			sql += " " + option
		}
	}

	iter := dbhelper.DBCassandra.Query(sql, params...).Iter()

	if nil == iter {
		return result, nil
	}

	data := &NumLog4s{}
	for iter.Scan(
		&data.CreateTime,
		&data.ID,
		&data.Num,
		&data.ServerID,
	) {
		result = append(result, data)

		data = &NumLog4s{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}
