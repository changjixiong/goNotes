package rankservice

import (
	"goNotes/redisnote/playercache"
	"log"

	"time"

	"strconv"

	redis "gopkg.in/redis.v4"
)

var Rds *redis.Client

const (
	PlayerLvRankKey = "Rank:PlayerLv"
)

type RankService struct {
}

var DefaultRankService = &RankService{}

func lvScoreWithTime(lv int, timeStamp int64) float64 {
	if 0 != timeStamp {
		return float64(lv) + (1<<14)/(float64(timeStamp)/float64(3600*24))

	} else {
		return float64(lv)
	}

}

func (rankService *RankService) GetPlayerByLvRank(start, count int64) []*playercache.PlayerInfo {

	playerInfos := []*playercache.PlayerInfo{}

	ids, err := Rds.ZRevRange(PlayerLvRankKey, start, start+count-1).Result()

	if nil != err {
		log.Println("RankService: GetPlayerByLvRank:", err)
		return playerInfos
	}

	for _, idstr := range ids {
		id, err := strconv.Atoi(idstr)

		if nil != err {
			log.Println("RankService: GetPlayerByLvRank:", err)
		} else {
			playerInfo := playercache.LoadPlayerInfo(id)

			if nil != playerInfos {
				playerInfos = append(playerInfos, playerInfo)
			}
		}
	}

	return playerInfos

}

func (rankService *RankService) SetPlayerLvRank(playerInfo *playercache.PlayerInfo) bool {

	if nil == playerInfo {
		return false
	}

	err := Rds.ZAdd(
		PlayerLvRankKey,
		redis.Z{
			Score:  lvScoreWithTime(playerInfo.Lv, time.Now().Unix()),
			Member: playerInfo.PlayerID,
		},
	).Err()

	if nil != err {
		log.Println("RankService: SetPlayerLvRank:", err)
		return false
	}

	return true
}

func (rankService *RankService) AddPlayerExp(playerID, exp int) bool {

	player := playercache.GetPlayerInfo(playerID)
	if nil == player {
		return false
	}

	player.Exp += exp
	// 固定经验升级，可以按需要修改
	if player.Exp >= playercache.LvUpExp {
		player.Lv += 1
		player.Exp = player.Exp - playercache.LvUpExp
		rankService.SetPlayerLvRank(player)
	}

	playercache.SetPlayerInfo(player)

	return true
}
