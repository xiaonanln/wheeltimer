package wheeltimer

import (
	"fmt"
)

type timerlist struct {
	head *Timer
	tail *Timer
}

var (
	errWrongTimerlist = fmt.Errorf("wrong timerlist")
)

func (tl *timerlist) add(t *Timer) {
	if wtDebug && t.next != nil || t.prev != nil {
		panic("timerlist add Timer failed: next or prev is not nil")
	}

	if tl.head == nil {
		if wtDebug && tl.tail != nil {
			panic(errWrongTimerlist)
		}
		tl.head = t
		tl.tail = t
	} else {
		// insert the timer after tail
		tail := tl.tail
		if wtDebug && (tail == nil || tail.next != nil) {
			panic(errWrongTimerlist)
		}

		tail.next = t.prev
		t.prev = tail
		tl.tail = t
	}

	if wtDebug {
		if tl.head != nil {
			if tl.tail == nil || tl.head.prev != nil || tl.tail.next != nil {
				panic(errWrongTimerlist)
			}
		} else {
			if tl.tail != nil {
				panic(errWrongTimerlist)
			}
		}
	}
}

func (tl *timerlist) remove(t *Timer) {

	if tl.head == t {
		if wtDebug && t.prev != nil {
			panic(errWrongTimerlist)
		}

		if tl.tail == t {
			if wtDebug && t.next != nil {
				panic(errWrongTimerlist)
			}
			// remove the single Timer
			tl.head, tl.tail = nil, nil
		} else {
			next := t.next
			t.next = nil
			if wtDebug && next == nil {
				panic(errWrongTimerlist)
			}
			// remove Timer from list head
			tl.head = next
			next.prev = nil
		}
	} else if tl.tail == t {
		if wtDebug && t.next != nil {
			panic(errWrongTimerlist)
		}

		// remove Timer from list tail
		prev := t.prev
		t.prev = nil
		if wtDebug && prev == nil {
			panic(errWrongTimerlist)
		}

		tl.tail = prev
		prev.next = nil
	} else {
		prev := t.prev
		next := t.next
		t.prev, t.next = nil, nil
		if wtDebug && (prev == nil || next == nil) {
			panic(errWrongTimerlist)
		}

		prev.next = next
		next.prev = prev
	}

	if wtDebug {
		if t.prev != nil || t.next != nil {
			panic(errWrongTimerlist)
		}

		if tl.head != nil {
			if tl.tail == nil || tl.head.prev != nil || tl.tail.next != nil {
				panic(errWrongTimerlist)
			}
		} else {
			if tl.tail != nil {
				panic(errWrongTimerlist)
			}
		}
	}
}
