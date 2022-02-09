package jobkicker

import (
	"time"
)

// delayToDuration turns time like time object with all fields as zero expect
// the seconds is 2 turned to `time.Duration(2 * time.Second)`
func delayToDuration(delay time.Time) time.Duration {
	zeroTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	return delay.Sub(zeroTime)
}
