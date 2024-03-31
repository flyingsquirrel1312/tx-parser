package lib

import "sync"

type ConcurrentSet[T comparable] struct {
	capacity uint32
	mutex    sync.RWMutex
	items    map[T]struct{}
}

func NewConcurrentSet[T comparable](capacity uint32) *ConcurrentSet[T] {
	return &ConcurrentSet[T]{
		capacity: capacity,
		items:    make(map[T]struct{}),
	}
}

// Add adds a new item to the set. Returns true if the item was added, false if the set is full.
func (s *ConcurrentSet[T]) Add(item T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if uint32(len(s.items)) == s.capacity {
		return false
	}
	s.items[item] = struct{}{}
	return true
}

func (s *ConcurrentSet[T]) Has(item T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.items[item]
	return ok
}
