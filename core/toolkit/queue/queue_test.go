package queue_test

import (
	"github.com/kim118000/core/toolkit/queue"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLockFreeQueue(t *testing.T) {
	const taskNum = 10000
	q := queue.NewLockFreeQueue(0)
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		for i := 0; i < taskNum; i++ {
			task := &queue.QueueElement{Elem: "a"}
			q.Enqueue(task)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < taskNum; i++ {
			task := &queue.QueueElement{Elem: "b"}
			q.Enqueue(task)
		}
		wg.Done()
	}()

	var counter int32
	go func() {
		for {
			task := q.Dequeue()
			if task != nil {
				atomic.AddInt32(&counter, 1)
			}
			if task == nil && atomic.LoadInt32(&counter) == 2*taskNum {
				break
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			task := q.Dequeue()
			if task != nil {
				atomic.AddInt32(&counter, 1)
			}
			if task == nil && atomic.LoadInt32(&counter) == 2*taskNum {
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()

	t.Logf("sent and received all %d tasks, %d", 2*taskNum, atomic.LoadInt32(&counter))
}


func TestLockFreeQueue1(t *testing.T) {
	const taskNum = 10000
	q := queue.NewLockFreeQueue(100)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < taskNum; i++ {
			task := &queue.QueueElement{Elem: "a"}
			q.Enqueue(task)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < taskNum; i++ {
			task := &queue.QueueElement{Elem: "b"}
			q.Enqueue(task)
		}
		wg.Done()
	}()

	wg.Wait()
}
