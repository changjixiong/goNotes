package playercache

import (
	"fmt"
	"goNotes/utils"

	"sync"

	redis "gopkg.in/redis.v4"
)

type PlayerInfo struct {
	PlayerID   int    `db:"player_id" json:"player_id"`
	PlayerName string `db:"player_name" json:"player_name"`
	Exp        int    `db:"exp" json:"exp"`
	Lv         int    `db:"lv" json:"lv"`
	Online     bool   `db:"online" json:"online"`
}

var Rds *redis.Client
var playerMapMutex sync.Mutex
var playerMap = map[int]*PlayerInfo{}

const (
	LvUpExp = 100
)

type setPlayerValueFunc func(ma map[string]string)

func GetPlayerInfo(playerID int) *PlayerInfo {
	playerMapMutex.Lock()
	defer playerMapMutex.Unlock()
	return playerMap[playerID]
}

func SetPlayerInfo(player *PlayerInfo) {
	if nil == player {
		return
	}
	playerMapMutex.Lock()
	defer playerMapMutex.Unlock()
	playerMap[player.PlayerID] = player

	SetPlayerAllValue(player)
}

func setPlayerValue(playerID int, f setPlayerValueFunc) {

	ma := map[string]string{}

	f(ma)

	if len(ma) > 0 {
		err := Rds.HMSet(fmt.Sprintf("playerInfo:%d", playerID), ma).Err()

		if nil != err {
			fmt.Println(err)
		}
	}

}

func SetPlayerOnline(playerID int, online bool) {

	f := func(ma map[string]string) {
		ma["online"] = utils.Any(online)
	}

	setPlayerValue(playerID, f)
}

func SetPlayerLv(playerID int, lv int) {

	f := func(ma map[string]string) {
		ma["lv"] = utils.Any(lv)
	}

	setPlayerValue(playerID, f)
}

func SetPlayerAllValue(p *PlayerInfo) {

	f := func(ma map[string]string) {

		for k, v := range utils.Struct2MapString(p) {
			ma[k] = v
		}
	}

	setPlayerValue(p.PlayerID, f)

}

func LoadPlayerInfo(playerID int) *PlayerInfo {
	mapobj, err := Rds.HGetAll(fmt.Sprintf("playerInfo:%d", playerID)).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(mapobj) < 1 {
		//redis中没有缓存，从库中获取并写如redis
		//这里略过
		return nil
	}

	p := &PlayerInfo{}

	if utils.MapString2Struct(mapobj, p) {
		return p
	} else {
		return nil
	}

}
