package jobkicker

import "time"

// ITimer is an interface to abstract Timer and Ticker to use thin in Job struct
type ITimer interface {
	Stop()
	GetWaiter() <-chan time.Time
}

// KickerTimer is a wrapper to time.Timer
type KickerTimer struct {
	Timer *time.Timer
}

// KickerTimer is a wrapper to time.Ticker
type KickerTicker struct {
	Ticker *time.Ticker
}

// Stop stops the timer
func (kickerTimer *KickerTimer) Stop() {
	kickerTimer.Timer.Stop()
}

// GetWaiter returns the chan of the Timer
func (kickerTimer *KickerTimer) GetWaiter() <-chan time.Time {
	return kickerTimer.Timer.C
}

// InitiateNewKickerTimer returns new KickerTimer set to a duration
func InitiateNewKickerTimer(d time.Duration) *KickerTimer {
	return &KickerTimer{Timer: time.NewTimer(d)}
}

// Stop stops the ticker
func (kickerTicker *KickerTicker) Stop() {
	kickerTicker.Ticker.Stop()
}

// GetWaiter returns the chan of the Ticker
func (kickerTicker *KickerTicker) GetWaiter() <-chan time.Time {
	return kickerTicker.Ticker.C
}

// InitiateNewKickerTicker returns new KickerTicker set to a duration
func InitiateNewKickerTicker(d time.Duration) *KickerTicker {
	return &KickerTicker{Ticker: time.NewTicker(d)}
}
