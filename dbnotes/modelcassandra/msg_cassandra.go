package cassandra

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"time"

	"github.com/gocql/gocql"
)

type Msg struct {
	Content    string     `db:"content" json:"content"`         //
	CreateTime time.Time  `db:"create_time" json:"create_time"` //
	ID         gocql.UUID `db:"id" json:"id"`                   //
}

type msgOp struct{}

var MsgOp = &msgOp{}
var DefaultMsg = &Msg{}

func (op *msgOp) Insert(m *Msg) (int64, error) {

	return op.InsertTx(dbhelper.DBCassandra, m)
}

func (op *msgOp) InsertTx(session *gocql.Session, m *Msg) (int64, error) {
	sql := "insert into msg(content,create_time,id) values(?,?,?)"
	if err := session.Query(
		sql,
		m.Content,
		m.CreateTime,
		gocql.TimeUUID(),
	).Exec(); err != nil {
		fmt.Println("InsertTx", err)
		return -1, err

	}

	return 0, nil
}

func (op *msgOp) QueryByMap(m map[string]interface{}, options []string) ([]*Msg, error) {
	result := []*Msg{}
	var params []interface{}

	sql := "select content,create_time,id from msg"

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

	data := &Msg{}
	for iter.Scan(
		&data.Content,
		&data.CreateTime,
		&data.ID,
	) {
		result = append(result, data)

		data = &Msg{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}

func (op *msgOp) QueryByMapComparison(m map[string]interface{}, options []string) ([]*Msg, error) {
	result := []*Msg{}
	var params []interface{}

	sql := "select content,create_time,id from msg"

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

	data := &Msg{}
	for iter.Scan(
		&data.Content,
		&data.CreateTime,
		&data.ID,
	) {
		result = append(result, data)

		data = &Msg{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}
