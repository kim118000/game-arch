package scheduler

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	var exists = struct {
		timers        int
		createdTimes  int
	}{
		timers:        len(timerManager.timers),
		createdTimes:  len(timerManager.createdTimer),
	}

	const tc = 10
	var counter int64
	for i := 0; i < tc; i++ {
		NewTimer(1*time.Millisecond, func() {
			atomic.AddInt64(&counter, 1)
		})
	}

	<-time.After(5 * time.Millisecond)
	cron()
	cron()
	if counter != tc*2 {
		t.Fatalf("expect: %d, got: %d", tc*2, counter)
	}

	if len(timerManager.timers) != exists.timers+tc {
		t.Fatalf("timers: %d", len(timerManager.timers))
	}

	if len(timerManager.createdTimer) != exists.createdTimes {
		t.Fatalf("createdTimer: %d", len(timerManager.createdTimer))
	}
}

func TestNewAfterTimer(t *testing.T) {
	var exists = struct {
		timers        int
		createdTimes  int
	}{
		timers:        len(timerManager.timers),
		createdTimes:  len(timerManager.createdTimer),
	}

	const tc = 10
	var counter int64
	for i := 0; i < tc; i++ {
		NewAfterTimer(1*time.Millisecond, func() {
			atomic.AddInt64(&counter, 1)
		})
	}

	<-time.After(5 * time.Millisecond)
	cron()
	if counter != tc {
		t.Fatalf("expect: %d, got: %d", tc, counter)
	}

	if len(timerManager.timers) != exists.timers {
		t.Fatalf("timers: %d", len(timerManager.timers))
	}

	if len(timerManager.createdTimer) != exists.createdTimes {
		t.Fatalf("createdTimer: %d", len(timerManager.createdTimer))
	}
}

func TestTimer_Stop(t *testing.T) {
	var counter int64
	tm := NewCountTimer(10*time.Second, 10, func() {
		atomic.AddInt64(&counter, 1)
	})

	go Sched()

	<- time.After(15 * time.Second)
	tm.Stop()

}