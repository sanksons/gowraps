package timer

import "time"

func GetHttpTime(t time.Time) string {
	return t.Format(time.RFC1123)
}

func GetCurrentTime(zone string) (time.Time, error) {
	if zone == "" {
		zone = "UTC"
	}
	cTime := time.Now()
	location, locErr := time.LoadLocation(zone)
	if locErr != nil {
		return cTime, locErr
	}
	return cTime.In(location), nil
}
