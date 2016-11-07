package model

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"time"

	"github.com/jmoiron/sqlx"
)

type Msg struct {
	Id         int        `db:"id" json:"id"`                   //
	SenderId   int        `db:"sender_id" json:"sender_id"`     //
	ReceiverId int        `db:"receiver_id" json:"receiver_id"` //
	Content    string     `db:"content" json:"content"`
	Status     int8       `db:"status" json:"status"`         //
	Createtime *time.Time `db:"createtime" json:"createtime"` //
}

var DefaultMsg = &Msg{}

func (m *Msg) GetByPK(id int) (*Msg, bool) {
	obj := &Msg{}
	sql := "select * from dbnote.msg where id=? "
	err := dbhelper.DB.Get(obj, sql,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return obj, true
}

func (m *Msg) Insert(msg *Msg) (int64, error) {
	return m.InsertTx(dbhelper.DB, msg)
}

func (m *Msg) InsertTx(ext sqlx.Ext, msg *Msg) (int64, error) {
	sql := "insert into dbnote.Msg(sender_id,receiver_id,content,status,createtime) values(?,?,?,?,?)"
	result, err := ext.Exec(sql,
		msg.SenderId,
		msg.ReceiverId,
		msg.Content,
		msg.Status,
		msg.Createtime,
	)
	if err != nil {
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (m *Msg) QueryByMap(ma map[string]interface{}) ([]*Msg, error) {
	result := []*Msg{}
	var params []interface{}

	sql := "select * from dbnote.Msg where 1=1 "
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
