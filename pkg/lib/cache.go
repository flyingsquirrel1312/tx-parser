package lib

import "sync"

type FIFOCache[K comparable, V any] struct {
	capacity uint32
	mutex    sync.RWMutex
	items    map[K]V
	queue    *CircularBuffer[K]
}

func NewFIFOCache[K comparable, V any](capacity uint32) *FIFOCache[K, V] {
	return &FIFOCache[K, V]{
		capacity: capacity,
		items:    make(map[K]V),
		queue:    NewCircularBuffer[K](capacity),
	}
}

func (c *FIFOCache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.items[key]
	return value, ok
}

// Set adds a new key-value pair to the cache. If the key already exists, it updates the value.
// If the cache is full, it removes the oldest item.
// Returns true if the key is new, false if it already exists.
func (c *FIFOCache[K, V]) Set(key K, value V) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.items[key]; exists {
		c.items[key] = value
		return false
	}
	L := c.queue.Len()
	if L == c.capacity {
		oldest := c.queue.Tail()
		delete(c.items, oldest)
	}
	c.queue.Enqueue(key)
	c.items[key] = value
	return true
}

func (c *FIFOCache[K, V]) Has(key K) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.items[key]
	return ok
}

func (c *FIFOCache[K, V]) Len() uint32 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return uint32(len(c.items))
}
