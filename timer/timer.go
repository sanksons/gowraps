package timer

import "time"

func GetHttpTime(t time.Time) string {
	return t.Format(time.RFC1123)
}
