package wheeltimer

import (
	"math/rand"
	"testing"
	"time"
)

var (
	MULTIPLY = 10000
)

func BenchmarkWheelTimer(b *testing.B) {
	wt := NewWheelTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < MULTIPLY; j++ {
			wt.After(time.Millisecond * time.Duration(rand.Intn(100000)))
		}
	}
}

func BenchmarkTimeLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < MULTIPLY; j++ {
			time.After(time.Millisecond * time.Duration(rand.Intn(100000)))
		}
	}
}
