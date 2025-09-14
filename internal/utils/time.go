package utils

import (
	"fmt"
	"time"
)

var tokyoLocation *time.Location

func init() {
	var err error
	tokyoLocation, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(fmt.Sprintf("failed to load time location: %v", err))
	}
}

func GetTokyoLocation() *time.Location {
	return tokyoLocation
}

// parseTime は引数のlayout形式の日時データをtime.Time(JST)に変換して返す
func parseTime(date, layout string) (time.Time, error) {

	t, err := time.ParseInLocation(layout, date, tokyoLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse time: (Date='%v'): %w", date, err)
	}

	// 時分秒ミリ秒以下を0にする場合
	// t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tokyoLocation)

	return t, nil
}

func ItemPerSleep(itemCount, perItems, sleepSeconds int) {
	if itemCount%perItems == 0 {
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}
