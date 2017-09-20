package cassandra

import (
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"time"

	"github.com/gocql/gocql"
)

type PlayerHeroUseLog struct {
	CreateTime time.Time  `db:"create_time" json:"create_time"` //
	HeroID     int        `db:"hero_id" json:"hero_id"`         //
	ID         gocql.UUID `db:"id" json:"id"`                   //
	PlayerID   int        `db:"player_id" json:"player_id"`     //
	RoomID     int        `db:"room_id" json:"room_id"`         //
}

type playerHeroUseLogOp struct{}

var PlayerHeroUseLogOp = &playerHeroUseLogOp{}
var DefaultPlayerHeroUseLog = &PlayerHeroUseLog{}

func (op *playerHeroUseLogOp) Insert(m *PlayerHeroUseLog) (int64, error) {
	return op.InsertTx(dbhelper.DBCassandra, m)
}

func (op *playerHeroUseLogOp) InsertTx(session *gocql.Session, m *PlayerHeroUseLog) (int64, error) {
	sql := "insert into player_hero_use_log(create_time,hero_id,id,player_id,room_id) values(?,?,?,?,?)"
	if err := session.Query(
		sql,
		m.CreateTime,
		m.HeroID,
		gocql.TimeUUID(),
		m.PlayerID,
		m.RoomID,
	).Exec(); err != nil {
		return -1, err

	}

	return 0, nil
}

func (op *playerHeroUseLogOp) QueryByMap(m map[string]interface{}, options []string) ([]*PlayerHeroUseLog, error) {
	result := []*PlayerHeroUseLog{}
	var params []interface{}

	sql := "select create_time,hero_id,id,player_id,room_id from player_hero_use_log"

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

	data := &PlayerHeroUseLog{}
	for iter.Scan(
		&data.CreateTime,
		&data.HeroID,
		&data.ID,
		&data.PlayerID,
		&data.RoomID,
	) {
		result = append(result, data)

		data = &PlayerHeroUseLog{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}

func (op *playerHeroUseLogOp) QueryByMapComparison(m map[string]interface{}, options []string) ([]*PlayerHeroUseLog, error) {
	result := []*PlayerHeroUseLog{}
	var params []interface{}

	sql := "select create_time,hero_id,id,player_id,room_id from player_hero_use_log"

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

	data := &PlayerHeroUseLog{}
	for iter.Scan(
		&data.CreateTime,
		&data.HeroID,
		&data.ID,
		&data.PlayerID,
		&data.RoomID,
	) {
		result = append(result, data)

		data = &PlayerHeroUseLog{}
	}

	if err := iter.Close(); err != nil {
		fmt.Println("err:", err)
	}

	return result, nil
}
