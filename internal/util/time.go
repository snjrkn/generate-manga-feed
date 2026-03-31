package util

import (
	"fmt"
	"time"
)

var jst *time.Location

func init() {
	jst = time.FixedZone("JST", 9*60*60)
}

func GetTokyoLocation() *time.Location {
	return jst
}

// parseTime は引数のlayout形式の日時データをtime.Time(JST)に変換して返す
func parseTime(date, layout string) (time.Time, error) {

	t, err := time.ParseInLocation(layout, date, jst)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse time: (Date='%v'): %w", date, err)
	}

	// 時分秒ミリ秒以下を0にする場合
	// t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, jst)

	return t, nil
}

func ItemPerSleep(itemCount, perItems, sleepSeconds int) {
	if itemCount%perItems == 0 {
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}
