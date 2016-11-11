package model

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"time"

	"github.com/jmoiron/sqlx"
)

type Notice struct {
	ID         int        `db:"id" json:"id"`                   //
	No         int        `db:"No" json:"No"`                   //
	SenderID   int        `db:"sender_id" json:"sender_id"`     // 发送者
	ReceiverID int        `db:"receiver_id" json:"receiver_id"` // 接收者
	Content    string     `db:"content" json:"content"`         // 内容
	Status     int8       `db:"status" json:"status"`           //
	Createtime *time.Time `db:"createtime" json:"createtime"`   //
}

var DefaultNotice = &Notice{}

func (m *Notice) GetByPK(id int, No int) (*Notice, bool) {
	obj := &Notice{}
	sql := "select * from dbnote.notice where id=? and No=?"
	err := dbhelper.DB.Get(obj, sql,
		id,
		No,
	)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return obj, true
}

func (m *Notice) Insert() (int64, error) {
	return m.InsertTx(dbhelper.DB)
}

func (m *Notice) InsertTx(ext sqlx.Ext) (int64, error) {
	sql := "insert into dbnote.notice(id,No,sender_id,receiver_id,content,status,createtime) values(?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.ID,
		m.No,
		m.SenderID,
		m.ReceiverID,
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

func (m *Notice) Delete() error {
	return m.DeleteTx(dbhelper.DB)
}

func (m *Notice) DeleteTx(ext sqlx.Ext) error {
	sql := `delete from dbnote.notice where id=? and No=?`
	_, err := ext.Exec(sql,
		m.ID,
		m.No,
	)
	return err
}

func (m *Notice) Update() error {
	return m.UpdateTx(dbhelper.DB)
}

func (m *Notice) UpdateTx(ext sqlx.Ext) error {
	sql := `update dbnote.notice set sender_id=?,receiver_id=?,content=?,status=?,createtime=? where id=? and No=?`
	_, err := ext.Exec(sql,
		m.SenderID,
		m.ReceiverID,
		m.Content,
		m.Status,
		m.Createtime,
		m.ID,
		m.No,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *Notice) QueryByMap(ma map[string]interface{}) ([]*Notice, error) {
	result := []*Notice{}
	var params []interface{}

	sql := "select * from dbnote.notice where 1=1 "
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
