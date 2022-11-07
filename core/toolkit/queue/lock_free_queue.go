package queue

import (
	"github.com/kim118000/core/pkg/log"
	"sync/atomic"
	"unsafe"
)

// lockFreeQueue is a simple, fast, and practical non-blocking and concurrent queue with no lock.
type lockFreeQueue struct {
	head     unsafe.Pointer
	tail     unsafe.Pointer
	length   int32
	capacity int32
}

type node struct {
	value *QueueElement
	next  unsafe.Pointer
}

// NewLockFreeQueue instantiates and returns a lockFreeQueue.
func NewLockFreeQueue(capacity int32) AsyncQueue {
	n := unsafe.Pointer(&node{})
	return &lockFreeQueue{head: n, tail: n, capacity: capacity}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *lockFreeQueue) Enqueue(elem *QueueElement) {
	if q.capacity > 0 {
		for {
			if q.Len() < q.capacity {
				break
			}
			drop := q.Dequeue()
			log.DefaultLogger.Warnf("queue full drop first element %v", drop)
		}
	}

	n := &node{value: elem}
retry:
	tail := load(&q.tail)
	next := load(&tail.next)
	// Are tail and next consistent?
	if tail == load(&q.tail) {
		if next == nil {
			// Try to link node at the end of the linked list.
			if cas(&tail.next, next, n) { // enqueue is done.
				// Try to swing tail to the inserted node.
				cas(&q.tail, tail, n)
				atomic.AddInt32(&q.length, 1)
				return
			}
		} else { // tail was not pointing to the last node
			// Try to swing tail to the next node.
			cas(&q.tail, tail, next)
		}
	}
	goto retry
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *lockFreeQueue) Dequeue() *QueueElement {
retry:
	head := load(&q.head)
	tail := load(&q.tail)
	next := load(&head.next)
	// Are head, tail, and next consistent?
	if head == load(&q.head) {
		// Is queue empty or tail falling behind?
		if head == tail {
			// Is queue empty?
			if next == nil {
				return nil
			}
			cas(&q.tail, tail, next) // tail is falling behind, try to advance it.
		} else {
			// Read value before CAS, otherwise another dequeue might free the next node.
			task := next.value
			if cas(&q.head, head, next) { // dequeue is done, return value.
				atomic.AddInt32(&q.length, -1)
				return task
			}
		}
	}
	goto retry
}

func (q *lockFreeQueue) Len() int32 {
	return atomic.LoadInt32(&q.length)
}

// IsEmpty indicates whether this queue is empty or not.
func (q *lockFreeQueue) IsEmpty() bool {
	return atomic.LoadInt32(&q.length) == 0
}

func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}

//Clear all queue item
func (q *lockFreeQueue) Clear() {
	for !q.IsEmpty() {
		item := q.Dequeue()
		PutQueueElement(item)
	}
}