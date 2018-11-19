package wheeltimer

import (
	"log"
	"time"
)

func Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}

func NewTimer(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}

type Timer struct {
}

type slot struct {
}

type WheelTimer struct {
	lv1slots  [10000]*slot
	lv2slots  [10000]*slot
	lv1cursor int
	curtime   time.Time
}

func NewWheelTimer() *WheelTimer {
	wt := &WheelTimer{
		curtime: time.Now(),
	}
	go wt.routine()
	return wt
}

func (wt *WheelTimer) AfterFunc(d time.Duration, f func()) *Timer {
	return nil
}

func (wt *WheelTimer) After(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	wt.AfterFunc(d, func() {
		ch <- time.Now()
	})
	return ch
}

func (wt *WheelTimer) routine() {
	for {
		now := time.Now()
		wt.fireTimers(now)
		time.Sleep(time.Millisecond)
	}
}

func (wt *WheelTimer) fireTimers(now time.Time) {
	elapsed := now.Sub(wt.curtime)
	steps := (elapsed / time.Millisecond) // steps to go forward, this should normally be 1
	log.Printf("steps %d", steps)
	for i := time.Duration(0); i < steps; i++ {
		wt.step()
	}
	wt.curtime = wt.curtime.Add(time.Millisecond * steps)
}

func (wt *WheelTimer) step() {

}
