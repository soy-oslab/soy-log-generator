package ring_test

import (
	"testing"

	"github.com/soyoslab/soy_log_generator/pkg/ring"
)

func setup(size uint64) *ring.Ring {
	r := new(ring.Ring)
	r.Init(size, "test")
	return r
}

func TestRing(t *testing.T) {
	r := setup(10)
	if r == nil {
		t.Errorf("initialization failed")
	}
	for i := 0; i < 16; i++ {
		err := r.Push(i)
		if err != nil {
			t.Errorf("push failed %v", err)
		}
	}
	ok, _ := r.Offer(16)
	if ok {
		t.Errorf("invalid push detected")
	}
	v := r.Pop(4)
	if len(v) != 4 {
		t.Errorf("ring buffer pop threshold failed")
	}
	v = r.Poll()
	if len(v) != 12 {
		t.Errorf("ring buffer polling pop failed")
	}
	v = r.Pop(1)
	if len(v) != 0 {
		t.Errorf("ring buffer size must be zero")
	}
}
