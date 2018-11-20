package wheeltimer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const wtDebug = true

const (
	wheelSize = 10000
)

type Timer struct {
	C          chan time.Time
	b          *bucket
	t          int64 // timeout time in milliseconds
	f          unsafe.Pointer
	prev, next *Timer
}

func (t *Timer) getf() func() {
	return *((*func())(atomic.LoadPointer(&t.f)))
}

func (t *Timer) setf(f func()) {
	atomic.StorePointer(&t.f, unsafe.Pointer(&f))
}

func (t *Timer) Stop() {
	t.setf(nil)
	t.b.removeTimer(t)
}

type WheelTimer struct {
	lv1wheel lv1wheel
	lv2wheel lv2wheel
	addMu    sync.Mutex
	addQueue timerlist
}

type lv1wheel struct {
	buckets [wheelSize]bucket
	cursor  int64
	curtime int64
}

type lv2wheel struct {
	buckets [wheelSize]bucket
	cursor  int64
	curtime int64
}

func NewWheelTimer() *WheelTimer {
	now := now()
	wt := &WheelTimer{}
	wt.lv1wheel.curtime = now
	wt.lv2wheel.curtime = now
	go wt.routine()
	return wt
}

func now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func (wt *WheelTimer) AfterFunc(d time.Duration, f func()) *Timer {
	ms := int64((d + time.Millisecond - 1) / time.Millisecond)
	timer := &Timer{
		C: nil,
		b: nil,
		t: now() + ms,
	}
	timer.setf(f)
	wt.addTimer(timer)
	return timer
}

func (wt *WheelTimer) After(d time.Duration) <-chan time.Time {
	ms := int64((d + time.Millisecond - 1) / time.Millisecond)
	timer := &Timer{
		C: make(chan time.Time, 1),
		b: nil,
		t: now() + ms,
	}
	timer.setf(func() {
		timer.C <- time.Now()
	})
	wt.addTimer(timer)
	return timer.C
}

func (wt *WheelTimer) NewTimer(d time.Duration) *Timer {
	ms := int64((d + time.Millisecond - 1) / time.Millisecond)
	timer := &Timer{
		C: make(chan time.Time, 1),
		b: nil,
		t: now() + ms,
		f: nil,
	}
	timer.setf(func() {
		timer.C <- time.Now()
	})
	wt.addTimer(timer)
	return timer
}

func (wt *WheelTimer) Tick(d time.Duration) <-chan time.Time {
	ms := int64((d + time.Millisecond - 1) / time.Millisecond)
	if ms <= 0 {
		ms = 1
	}
	timer := &Timer{
		C: make(chan time.Time, 1),
		b: nil,
		t: now() + ms,
		f: nil,
	}

	timer.setf(func() {
		select {
		case timer.C <- time.Now():
		default:
		}

		// re-add the timer
		if timer.b != nil {
			panic("timer.b != nil")
		}
		timer.t += ms
		wt.addTimer(timer)
	})

	wt.addTimer(timer)
	return timer.C
}

type Ticker struct {
	C chan time.Time
}

func (wt *WheelTimer) NewTicker(d time.Duration) *Ticker {
	ms := int64((d + time.Millisecond - 1) / time.Millisecond)
	if ms <= 0 {
		ms = 1
	}
	ticker := &Ticker{
		C: make(chan time.Time, 1),
	}
	timer := &Timer{
		C: nil,
		b: nil,
		t: now() + ms,
		f: nil,
	}

	timer.setf(func() {
		select {
		case ticker.C <- time.Now():
		default:
		}

		// re-add the timer
		if timer.b != nil {
			panic("timer.b != nil")
		}
		timer.t += ms
		wt.addTimer(timer)
	})

	wt.addTimer(timer)
	return ticker
}

func (wt *WheelTimer) routine() {
	for {
		now := now()
		wt.fireTimers(now)
		time.Sleep(time.Millisecond)
	}
}

func (wt *WheelTimer) fireTimers(now int64) {
	//log.Printf("wt fire timers: now=%d, lv1wheel.curtime=%d, lv2wheel.curtime=%d, lv1steps=%d", now, wt.lv1wheel.curtime, wt.lv2wheel.curtime, now-wt.lv1wheel.curtime)
	wt.processAddQueue()

	lv1steps := now - wt.lv1wheel.curtime // steps to go forward, this should normally be 1
	for i := int64(0); i < lv1steps; i++ {
		timers := wt.lv1wheel.step()
		for t := timers.head; t != nil; t = t.next {
			if wtDebug {
				if t.t >= wt.lv1wheel.curtime {
					panic(fmt.Errorf("curtime and timer deadline mismatch"))
				}
			}
			f := t.getf()
			if f != nil {
				go f()
			}
		}
	}

	lv2steps := (now - wt.lv2wheel.curtime) / wheelSize
	for i := int64(0); i < lv2steps; i++ {
		timers := wt.lv2wheel.step()
		t := timers.head
		for t != nil {
			// re-add the timer to the wheel timer
			next := t.next
			t.next, t.prev = nil, nil
			wt.addTimer(t)
			t = next
		}
	}
}

func (w *lv1wheel) step() (timers timerlist) {
	bucket := &w.buckets[w.cursor]
	w.cursor = (w.cursor + 1) % wheelSize
	w.curtime += 1
	return bucket.popTimers()
}

func (w *lv2wheel) step() (timers timerlist) {
	bucket := &w.buckets[w.cursor]
	w.cursor = (w.cursor + 1) % wheelSize
	w.curtime += wheelSize
	return bucket.popTimers()
}

func (wt *WheelTimer) addTimer(t *Timer) {
	wt.addMu.Lock()
	defer wt.addMu.Unlock()

	//log.Printf("addTimer: %p, cur queue: %v", t, wt.addQueue)
	wt.addQueue.add(t)
}

func (wt *WheelTimer) addTimerImpl(t *Timer) {
	if !wt.lv1wheel.addTimer(t) {
		//log.Printf("add to lv1wheel failed, adding to lv2wheel")
		wt.lv2wheel.addTimer(t)
	}
}

func (wt *WheelTimer) processAddQueue() {
	var timers timerlist
	wt.addMu.Lock()
	timers, wt.addQueue = wt.addQueue, timerlist{}
	wt.addMu.Unlock()

	t := timers.head
	for t != nil {
		next := t.next
		t.prev, t.next = nil, nil
		//log.Printf("processAddQueue: %p, %p", t, t.next)
		wt.addTimerImpl(t)
		t = next
	}
}

func (w *lv1wheel) addTimer(t *Timer) bool {
	d := t.t - w.curtime
	//log.Printf("lv1wheel addTimer: curtime=%d, t=%d, d=%d", w.curtime, t.t, d)
	if d >= wheelSize {
		return false
	} else if d < 0 {
		d = 0
	}
	w.buckets[(w.cursor+d)%wheelSize].addTimer(t)
	return true
}

func (w *lv2wheel) addTimer(t *Timer) {
	if t.t < w.curtime+wheelSize {
		panic(fmt.Errorf("t.t is %d, but should be at least %d", t.t, w.curtime+wheelSize))
	}

	d := (t.t - w.curtime - wheelSize) / wheelSize
	//log.Printf("lv2wheel addTimer: curtime=%d, t=%d, d=%d", w.curtime, t.t, d)
	if d >= wheelSize {
		d = wheelSize - 1
	}
	w.buckets[(w.cursor+d)%wheelSize].addTimer(t)
}
