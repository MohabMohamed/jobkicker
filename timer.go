package jobkicker

import "time"

type ITimer interface {
	Stop()
	GetWaiter() <-chan time.Time
}

type KickerTimer struct {
	Timer *time.Timer
}

type KickerTicker struct {
	Ticker *time.Ticker
}

func (kickerTimer *KickerTimer) Stop() {
	kickerTimer.Timer.Stop()
}
func (kickerTimer *KickerTimer) GetWaiter() <-chan time.Time {
	return kickerTimer.Timer.C
}

func InitiateNewKickerTimer(d time.Duration) *KickerTimer {
	return &KickerTimer{Timer: time.NewTimer(d)}
}

func (kickerTicker *KickerTicker) Stop() {
	kickerTicker.Ticker.Stop()
}

func (kickerTicker *KickerTicker) GetWaiter() <-chan time.Time {
	return kickerTicker.Ticker.C
}

func InitiateNewKickerTicker(d time.Duration) *KickerTicker {
	return &KickerTicker{Ticker: time.NewTicker(d)}
}
