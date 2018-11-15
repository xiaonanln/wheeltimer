package wheeltimer

import (
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	d := time.Millisecond * 100
	now := time.Now()
	ch := After(d)
	<-ch
	dt := time.Now().Sub(now)
	if dt < d || dt > d+time.Second {
		t.Fatalf("wrong time")
	}
}

func TestAfterFunc(t *testing.T) {
	d := time.Millisecond * 100
	now := time.Now()
	ch := make(chan struct{}, 1)
	AfterFunc(d, func() {
		ch <- struct{}{}
	})
	<-ch
	dt := time.Now().Sub(now)
	if dt < d || dt > d+time.Second {
		t.Fatalf("wrong time")
	}
}

func TestNewTimer(t *testing.T) {
	d := time.Millisecond * 100
	now := time.Now()
	timer := NewTimer(d)
	<-timer.C
	dt := time.Now().Sub(now)
	if dt < d || dt > d+time.Second {
		t.Fatalf("wrong time")
	}
}

func TestTick(t *testing.T) {
	d := time.Millisecond * 100

	ch := Tick(d)
	t0 := time.Now()
	for i := 0; i < 3; i++ {
		<-ch
		dt := time.Now().Sub(t0)
		if dt < d || dt > d+time.Second {
			t.Fatalf("wrong time")
		}
		t0 = time.Now()
	}
}

func TestNewTicker(t *testing.T) {
	d := time.Millisecond * 100

	ticker := NewTicker(d)
	t0 := time.Now()
	for i := 0; i < 3; i++ {
		<-ticker.C
		dt := time.Now().Sub(t0)
		if dt < d || dt > d+time.Second {
			t.Fatalf("wrong time")
		}
		t0 = time.Now()
	}
}
