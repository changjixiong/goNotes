package model

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"time"

	"github.com/jmoiron/sqlx"
)

type Mail struct {
	ID         int        `db:"id" json:"id"`                   //
	SenderID   int        `db:"sender_id" json:"sender_id"`     //
	ReceiverID int        `db:"receiver_id" json:"receiver_id"` //
	Title      string     `db:"title" json:"title"`             //
	Content    string     `db:"content" json:"content"`         //
	Status     int8       `db:"status" json:"status"`           //
	Createtime *time.Time `db:"createtime" json:"createtime"`   //
}

var DefaultMail = &Mail{}

func (m *Mail) GetByPK(id int) (*Mail, bool) {
	obj := &Mail{}
	sql := "select * from mail where id=?"
	err := dbhelper.DB.Get(obj, sql,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return obj, true
}

func (m *Mail) Insert() (int64, error) {
	return m.InsertTx(dbhelper.DB)
}

func (m *Mail) InsertTx(ext sqlx.Ext) (int64, error) {
	sql := "insert into mail(id,sender_id,receiver_id,title,content,status,createtime) values(?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.ID,
		m.SenderID,
		m.ReceiverID,
		m.Title,
		m.Content,
		m.Status,
		m.Createtime,
	)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (m *Mail) Delete() error {
	return m.DeleteTx(dbhelper.DB)
}

func (m *Mail) DeleteTx(ext sqlx.Ext) error {
	sql := `delete from mail where id=?`
	_, err := ext.Exec(sql,
		m.ID,
	)
	return err
}

func (m *Mail) Update() error {
	return m.UpdateTx(dbhelper.DB)
}

func (m *Mail) UpdateTx(ext sqlx.Ext) error {
	sql := `update mail set sender_id=?,receiver_id=?,title=?,content=?,status=?,createtime=? where id=?`
	_, err := ext.Exec(sql,
		m.SenderID,
		m.ReceiverID,
		m.Title,
		m.Content,
		m.Status,
		m.Createtime,
		m.ID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *Mail) QueryByMap(ma map[string]interface{}) ([]*Mail, error) {
	result := []*Mail{}
	var params []interface{}

	sql := "select * from mail where 1=1 "
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
