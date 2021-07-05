package compressor_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/soyoslab/soy_log_generator/pkg/compressor"
)

func TestGzipCompressorSuccess(t *testing.T) {
	source := "Hello World"
	// var c compressor.Compressor // This also valid
	c := &compressor.GzipComp{}
	buffer, err := c.Compress([]byte(source))
	if err != nil {
		t.Errorf("compress failed %v", err)
	}
	bytes, err := c.Decompress(buffer)
	if err != nil {
		t.Errorf("decompress failed %v", err)
	}
	target := string(bytes)

	if source != target {
		t.Errorf("%s(source) != %s(target)", source, target)
	}
}

func TestZstdCompressorSuccess(t *testing.T) {
	source := "Hello World"
	c := &compressor.ZstdComp{}
	buffer, err := c.Compress([]byte(source))
	if err != nil {
		t.Errorf("compress failed %v", err)
	}
	bytes, err := c.Decompress(buffer)
	if err != nil {
		t.Errorf("decompress failed %v", err)
	}
	target := string(bytes)
	if source != target {
		t.Errorf("%s(source) != %s(target)", source, target)
	}
}

func TestGzipDecompressFailed(t *testing.T) {
	c := &compressor.GzipComp{}
	_, err := c.Decompress(nil)
	if err == nil {
		t.Errorf("error must be specified")
	}
}

func getRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "ExcbsVQs"
}

func BenchmarkGzipCompress16MB(b *testing.B) {
	c := &compressor.GzipComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096 * 4096)
		b.StartTimer()
		// execution in here
		result, _ := c.Compress([]byte(message))
		_ = result
		// fmt.Printf("\n%v to %v\n", len([]byte(message)), len(result))
	}
}

func BenchmarkGzipDecompress16MB(b *testing.B) {
	c := &compressor.GzipComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096 * 4096)
		buffer, _ := c.Compress([]byte(message))
		b.StartTimer()
		// execution in here
		result, _ := c.Decompress(buffer)
		_ = result
	}
}

func BenchmarkZstdCompress16MB(b *testing.B) {
	c := &compressor.ZstdComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096 * 4096)
		b.StartTimer()
		// execution in here
		result, _ := c.Compress([]byte(message))
		_ = result
		// fmt.Printf("\n%v to %v\n", len([]byte(message)), len(result))
	}
}

func BenchmarkZstdDecompress16MB(b *testing.B) {
	c := &compressor.ZstdComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096 * 4096)
		buffer, _ := c.Compress([]byte(message))
		b.StartTimer()
		// execution in here
		result, _ := c.Decompress(buffer)
		_ = result
	}
}

func BenchmarkGzipCompress4KB(b *testing.B) {
	c := &compressor.GzipComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096)
		b.StartTimer()
		// execution in here
		result, _ := c.Compress([]byte(message))
		_ = result
		// fmt.Printf("\n%v to %v\n", len([]byte(message)), len(result))
	}
}

func BenchmarkGzipDecompress4KB(b *testing.B) {
	c := &compressor.GzipComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096)
		buffer, _ := c.Compress([]byte(message))
		b.StartTimer()
		// execution in here
		result, _ := c.Decompress(buffer)
		_ = result
	}
}

func BenchmarkZstdCompress4KB(b *testing.B) {
	c := &compressor.ZstdComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096)
		b.StartTimer()
		// execution in here
		result, _ := c.Compress([]byte(message))
		_ = result
		// fmt.Printf("\n%v to %v\n", len([]byte(message)), len(result))
	}
}

func BenchmarkZstdDecompress4KB(b *testing.B) {
	c := &compressor.ZstdComp{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		message := getRandomString(4096)
		buffer, _ := c.Compress([]byte(message))
		b.StartTimer()
		// execution in here
		result, _ := c.Decompress(buffer)
		_ = result
	}
}
