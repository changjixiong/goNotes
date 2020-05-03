package main

import (
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func getTargerTime(hour, minute, second int, offset int64) int64 {
	utcTime := time.Now().UTC()
	targetTime := time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(),
		hour, minute, second, 0, utcTime.Location())

	return targetTime.Unix() + offset
}

type RefreshConfig struct {
	TargetHour      int
	TargetMinute    int
	Targetsecond    int
	Offset          int64
	lastRefreshTime int64
}

var zoneToOffset = map[string]int64{
	"Z0":  0,
	"E1":  -1 * 3600,
	"E2":  -2 * 3600,
	"E3":  -3 * 3600,
	"E4":  -4 * 3600,
	"E5":  -5 * 3600,
	"E6":  -6 * 3600,
	"E7":  -7 * 3600,
	"E8":  -8 * 3600,
	"E9":  -9 * 3600,
	"E10": -10 * 3600,
	"E11": -11 * 3600,
	"E12": 12 * 3600,
	"W1":  1 * 3600,
	"W2":  2 * 3600,
	"W3":  3 * 3600,
	"W4":  4 * 3600,
	"W5":  5 * 3600,
	"W6":  6 * 3600,
	"W7":  7 * 3600,
	"W8":  8 * 3600,
	"W9":  9 * 3600,
	"W10": 10 * 3600,
	"W11": 11 * 3600,
	"W12": 12 * 3600,
}

func TimeIsUp(refreshConfig *RefreshConfig) bool {

	targetTime := getTargerTime(refreshConfig.TargetHour,
		refreshConfig.TargetMinute,
		refreshConfig.Targetsecond,
		refreshConfig.Offset)

	return refreshConfig.lastRefreshTime < targetTime &&
		time.Now().Unix() >= targetTime
}

func main() {

	refreshConfigs := []*RefreshConfig{}

	refreshConfigs = append(refreshConfigs, &RefreshConfig{TargetHour: 20,
		TargetMinute: 5,
		Targetsecond: 0,
		Offset:       zoneToOffset["E8"]})

	refreshConfigs = append(refreshConfigs, &RefreshConfig{TargetHour: 9,
		TargetMinute: 6,
		Targetsecond: 0,
		Offset:       zoneToOffset["W3"]})

	for {
		fmt.Println("server Time:", time.Now().Format(timeFormat))

		for i, r := range refreshConfigs {
			if TimeIsUp(r) {
				fmt.Println(i, "canRefresh")
				r.lastRefreshTime = time.Now().Unix()
			}
			fmt.Println("wait...")
			time.Sleep(time.Second * 5)
		}

	}

}
