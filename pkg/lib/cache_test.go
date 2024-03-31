package lib

import (
	"testing"
)

func TestFIFOCache(t *testing.T) {
	cache := NewFIFOCache[int, int](3)
	cache.Set(1, 1)
	cache.Set(2, 2)
	cache.Set(3, 3)
	v, ok := cache.Get(1)
	if !ok {
		t.Fatal("expected key 1 to exist")
	}
	if v != 1 {
		t.Fatalf("expected key 1 to have value 1, got %d", v)
	}
	cache.Set(4, 4)
	v, ok = cache.Get(1)
	if ok {
		t.Fatal("expected key 1 to be removed")
	}
}
