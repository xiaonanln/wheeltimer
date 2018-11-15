package wheeltimer

import "time"

func After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func AfterFunc(d time.Duration, f func()) *time.Timer {
	return time.AfterFunc(d, f)
}

func Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}

func NewTimer(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}
