package wheeltimer

import (
	"fmt"
	"sync"
)

type bucket struct {
	mu     sync.Mutex
	timers timerlist
}

func (b *bucket) addTimer(t *Timer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if t.b != nil {
		panic(fmt.Errorf("timer bucket should be nil"))
	}

	b.timers.add(t)
}

func (b *bucket) removeTimer(t *Timer) {
	b.mu.Lock()
	defer b.mu.Unlock()

}

func (b *bucket) popTimers() (timers timerlist) {
	b.mu.Lock()
	defer b.mu.Unlock()
	timers, b.timers = b.timers, timerlist{}
	// clear bucket for all timers

	return timers
}
