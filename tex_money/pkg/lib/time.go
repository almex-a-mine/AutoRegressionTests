package lib

import (
	"time"
)

// GeDateTime 現在の日付と時刻をstring型で取得
func GeDateTime() (string, string, error) {
	// 日本時間に設定
	time.Local = time.FixedZone("Local", 9*60*60)
	jst, err := time.LoadLocation("Local")
	if err != nil {
		return "", "", err
	}

	// 現在日時を取得
	dateNow := time.Now().In(jst).Format("20060102")
	timeNow := time.Now().In(jst).Format("150405070")

	return dateNow, timeNow, nil
}
