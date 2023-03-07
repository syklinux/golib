package collection

import (
	"fmt"
	"time"
)

func FormatDurationToMs(d time.Duration) string {
	return fmt.Sprintf("%.2f", float64(d.Nanoseconds())/float64(time.Millisecond))
}

func TimeBetween(current time.Time, startTime string, endTime string) bool {
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)

	if current.After(start) && current.Before(end) {
		return true
	}

	return false
}
