package compressor

import (
	"bytes"
	"encoding/gob"
	"github.com/bxcodec/faker/v3"
	"testing"
)

type TestStruct struct {
	Inta  int   `faker:"boundary_start=5, boundary_end=10"`
	Int8  int8  `faker:"boundary_start=100, boundary_end=1000"`
	Int16 int16 `faker:"boundary_start=123, boundary_end=1123"`
	Int32 int32 `faker:"boundary_start=-10, boundary_end=8123"`
	Int64 int64 `faker:"boundary_start=31, boundary_end=88"`

	UInta  uint   `faker:"boundary_start=35, boundary_end=152"`
	UInt8  uint8  `faker:"boundary_start=5, boundary_end=1425"`
	UInt16 uint16 `faker:"boundary_start=245, boundary_end=2125"`
	UInt32 uint32 `faker:"boundary_start=0, boundary_end=40"`
	UInt64 uint64 `faker:"boundary_start=14, boundary_end=50"`

	ASString []string          `faker:"len=100"`
	SString  string            `faker:"len=100"`
	MSString map[string]string `faker:"len=100"`
	MIint    map[int]int       `faker:"boundary_start=5, boundary_end=10"`
}

func TestCompressor(t *testing.T) {
	source := "Hello World"
	buffer := Compress([]byte(source))
	target := string(Decompress(buffer))
	if source != target {
		t.Errorf("%s(source) != %s(target)", source, target)
	}
}

func BenchmarkCompress1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		data := TestStruct{}
		var buffer bytes.Buffer
		_ = faker.SetRandomMapAndSliceSize(1000)
		_ = faker.FakeData(&data)
		encoder := gob.NewEncoder(&buffer)
		_ = encoder.Encode(data)
		b.StartTimer()
		// execution in here
		_ = Compress(buffer.Bytes())
	}
}

func BenchmarkDecompress1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		// data generation in here
		data := TestStruct{}
		var temp bytes.Buffer
		_ = faker.SetRandomMapAndSliceSize(1000)
		_ = faker.FakeData(&data)
		encoder := gob.NewEncoder(&temp)
		_ = encoder.Encode(data)
		buffer := Compress(temp.Bytes())
		b.StartTimer()
		// execution in here
		_ = Decompress(buffer)
	}
}
