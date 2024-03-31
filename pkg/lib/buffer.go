package lib

import "sync"

type CircularBuffer[T any] struct {
	items []T
	size  uint32
	head  uint32
	tail  uint32
	mutex sync.RWMutex
}

func NewCircularBuffer[T any](size uint32) *CircularBuffer[T] {
	return &CircularBuffer[T]{
		items: make([]T, size),
		size:  size,
	}
}

func (q *CircularBuffer[T]) Enqueue(item T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if (q.tail+1)%q.size == q.head {
		q.head = (q.head + 1) % q.size
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.size
}

func (q *CircularBuffer[T]) LastN(n uint32) []T {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if n > q.size {
		n = q.size
	}
	if q.tail >= q.head {
		return q.items[q.tail-n : q.tail]
	}
	return append(q.items[q.size-n+q.tail:], q.items[:q.tail]...)
}

func (q *CircularBuffer[T]) Len() uint32 {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if q.tail >= q.head {
		return q.tail - q.head
	}
	return q.size - q.head + q.tail + 1
}

func (q *CircularBuffer[T]) Head() T {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.items[q.head]
}

func (q *CircularBuffer[T]) Tail() T {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.items[q.tail]
}
