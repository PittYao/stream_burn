package helper

import "time"

func GetTimeNowStr() string {
	t := time.Now()
	now := t.Format("2006.01.02-15.04.05")
	return now
}

func GetTimeNow() *time.Time {
	t := time.Now()
	return &t
}

// GetDisuseAtTime 获取文件的淘汰时间
func GetDisuseAtTime(rootTime time.Time, day int) *time.Time {
	if rootTime.IsZero() {
		rootTime = time.Now()
	}
	hour := 24 * day
	duration := time.Duration(hour)

	disuseAt := rootTime.Add(time.Hour * duration)
	return &disuseAt
}
