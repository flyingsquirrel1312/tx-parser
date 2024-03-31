package lib

import (
	"fmt"
	"testing"
)

func TestCircularQueue(*testing.T) {
	q := NewCircularBuffer[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)
	q.Enqueue(5)
	fmt.Println(q.LastN(3))
	q.Enqueue(6)
	fmt.Println(q.LastN(3))
}
