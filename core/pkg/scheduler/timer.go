package scheduler

import (
	"container/heap"
	logger2 "github.com/kim118000/core/pkg/logger"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	infinite = -1
)

type (
	// TimerFunc represents a function which will be called periodically in main
	// logic gorontine.
	TimerFunc func()

	// TimerCondition represents a checker that returns true when cron job needs
	// to execute
	TimerCondition interface {
		Check(now time.Time) bool
	}

	TimerManager struct {
		incrementID int64      // auto increment id
		timers      TimerSlice // all timers

		muCreatedTimer sync.RWMutex
		createdTimer   []*Timer
	}

	// Timer represents a cron job
	Timer struct {
		id       int64         // timer id
		fn       TimerFunc     // function that execute
		createAt int64         // timer create time
		interval time.Duration // execution interval
		elapse   int64         // total elapse time
		closed   int32         // is timer closed
		counter  int           // counter
	}
)

type TimerSlice []*Timer

func (h TimerSlice) Len() int           { return len(h) }
func (h TimerSlice) Less(i, j int) bool { return h[i].elapse < h[j].elapse }
func (h TimerSlice) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *TimerSlice) Push(x interface{}) {
	*h = append(*h, x.(*Timer))
}
func (h *TimerSlice) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}


var timerManager TimerManager

func init() {
	var timerHeap TimerSlice = make([]*Timer, 0)
	heap.Init(&timerHeap)
	timerManager.timers = timerHeap
}

// ID returns id of current timer
func (t *Timer) ID() int64 {
	return t.id
}

// Stop turns off a timer. After Stop, fn will not be called forever
func (t *Timer) Stop() {
	if atomic.AddInt32(&t.closed, 1) != 1 {
		return
	}
	t.counter = 0
}

// execute job function with protection
func safecall(id int64, fn TimerFunc) {
	defer func() {
		if err := recover(); err != nil {
			logger2.DefaultLogger.Errorf("Handle timer panic: %+v %s", err, debug.Stack())
		}
	}()

	fn()
}

func cron() {
	if len(timerManager.createdTimer) > 0 {
		timerManager.muCreatedTimer.Lock()
		for _, t := range timerManager.createdTimer {
			heap.Push(&timerManager.timers, t)
		}
		timerManager.createdTimer = timerManager.createdTimer[:0]
		timerManager.muCreatedTimer.Unlock()
	}

	now := time.Now()
	unn := now.UnixNano()

	for timerManager.timers.Len() > 0{
		t := heap.Pop(&timerManager.timers).(*Timer)

		if t.createAt+t.elapse > unn {
			heap.Push(&timerManager.timers, t)
			return
		}

		if t.counter == infinite || t.counter > 0 {
			safecall(t.id, t.fn)
			t.elapse += int64(t.interval)
			logger2.DefaultLogger.Infof("timer execute id=%d counter=%d", t.id, t.counter)
			// update timer counter
			if t.counter != infinite && t.counter > 0 {
				t.counter--
			}
		}

		if t.counter == infinite || t.counter > 0 {
			heap.Push(&timerManager.timers, t)
		}
	}

}

func NewTimer(interval time.Duration, fn TimerFunc) *Timer {
	return NewCountTimer(interval, infinite, fn)
}

func NewCountOldTimer(timer *Timer) {
	timerManager.muCreatedTimer.Lock()
	timerManager.createdTimer = append(timerManager.createdTimer, timer)
	timerManager.muCreatedTimer.Unlock()
}

func NewCountTimer(interval time.Duration, count int, fn TimerFunc) *Timer {
	if fn == nil {
		panic("timer: nil timer function")
	}
	if interval <= 0 {
		panic("non-positive interval for NewTimer")
	}

	t := &Timer{
		id:       atomic.AddInt64(&timerManager.incrementID, 1),
		fn:       fn,
		createAt: time.Now().UnixNano(),
		interval: interval,
		elapse:   int64(interval), // first execution will be after interval
		counter:  count,
	}

	timerManager.muCreatedTimer.Lock()
	timerManager.createdTimer = append(timerManager.createdTimer, t)
	timerManager.muCreatedTimer.Unlock()
	return t
}

func NewAfterTimer(interval time.Duration, fn TimerFunc) *Timer {
	return NewCountTimer(interval, 1, fn)
}


func NewAfterTimerBySecondOnce(second int, fn TimerFunc) *Timer {
	return NewCountTimer(time.Duration(second) * time.Second, 1, fn)
}

func NewAfterTimerBySecond(second int, fn TimerFunc, count int) *Timer {
	return NewCountTimer(time.Duration(second) * time.Second, count, fn)
}

func NewAfterTimerBySecondForever(second int, fn TimerFunc) *Timer {
	return NewCountTimer(time.Duration(second) * time.Second, infinite, fn)
}