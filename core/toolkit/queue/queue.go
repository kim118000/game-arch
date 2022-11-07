package queue

import "sync"

// Task is a wrapper that contains function and its argument.
type QueueElement struct {
	Elem interface{}
}

// AsyncQueue is a queue storing asynchronous tasks.
type AsyncQueue interface {
	Enqueue(*QueueElement)
	Dequeue() *QueueElement
	IsEmpty() bool
	Len() int32
	Clear()
}

var pool = sync.Pool{New: func() interface{} { return new(QueueElement) }}

func GetQueueElement() *QueueElement {
	return pool.Get().(*QueueElement)
}

func PutQueueElement(elem *QueueElement) {
	elem.Elem = nil
	pool.Put(elem)
}