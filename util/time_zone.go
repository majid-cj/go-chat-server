package util

import (
	// "os"
	"time"
)

// GetTimeNow ...
func GetTimeNow() time.Time {
	loc, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}

// TimeAfter ...
func TimeAfter(timeBefore time.Time, duration time.Duration) time.Time {
	return timeBefore.Add(duration)
}
