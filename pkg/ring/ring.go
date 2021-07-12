package ring

import (
	"time"

	"github.com/Workiva/go-datastructures/queue"
)

// Ring contains the ring buffer information
type Ring struct {
	BufferType string
	buffer     *queue.RingBuffer
	Kick       chan bool
}

// Init initializes the ring buffer
func (r *Ring) Init(ringCapacity uint64, bufferType string) {
	r.buffer = queue.NewRingBuffer(ringCapacity)
	r.BufferType = bufferType
	r.Kick = make(chan bool, ringCapacity)
}

// Offer inserts a value to the ring buffer (non-blocking)
func (r *Ring) Offer(v interface{}) (bool, error) {
	return r.buffer.Offer(v)
}

// Push inserts a value to the ring buffer (blocking)
func (r *Ring) Push(v interface{}) error {
	return r.buffer.Put(v)
}

// Pop receives the number of values in ring until reach the threshold
// zero(0) means infinitely pop data until ring buffer is empty.
func (r *Ring) Pop(threshold uint64) []interface{} {
	buffer := []interface{}{}
	counter := threshold
	for {
		if threshold != 0 && counter <= 0 {
			break
		}
		v, err := r.buffer.Poll(time.Duration(1) * time.Millisecond)
		if err != nil {
			break
		}
		buffer = append(buffer, v)
		counter--
	}
	return buffer
}

// Poll receives the number of values in ring until it is empty
func (r *Ring) Poll() []interface{} {
	return r.Pop(0)
}
