package jobkicker

import (
	"time"
)

func delayToDuration(delay time.Time) time.Duration {
	zeroTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	return delay.Sub(zeroTime)
}
