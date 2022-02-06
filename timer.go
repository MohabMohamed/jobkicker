package jobkicker

import "time"

type ITimer interface {
	Stop()
	GetWaiter() <-chan time.Time
	InitiateNew(d time.Duration)
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

func (kickerTimer *KickerTimer) InitiateNew(d time.Duration) {
	kickerTimer.Timer = time.NewTimer(d)
}

func (kickerTicker *KickerTicker) Stop() {
	kickerTicker.Ticker.Stop()
}

func (kickerTicker *KickerTicker) GetWaiter() <-chan time.Time {
	return kickerTicker.Ticker.C
}

func (kickerTicker *KickerTicker) InitiateNew(d time.Duration) {
	kickerTicker.Ticker = time.NewTicker(d)
}
