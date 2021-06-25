package compressor

import "testing"

func TestHello(t *testing.T) {
	message := Hello("sample")
	if len(message) <= 0 {
		t.Errorf("length of message (=%d) <= 0", len(message))
	}
}
