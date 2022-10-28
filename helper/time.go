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

// TimeCompare 对比时间的大小
func TimeCompare(t1, t2 string) (error, bool, *time.Time, *time.Time) {
	time1, err := time.ParseInLocation("2006-01-02 15:04:05", t1, time.Local)
	time2, err := time.ParseInLocation("2006-01-02 15:04:05", t2, time.Local)

	now := time.Now()

	if err != nil {
		return err, false, &now, &now
	}

	if time1.Before(time2) {
		return nil, true, &time1, &time2
	}
	return nil, false, &time1, &time2
}
