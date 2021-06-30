package compressor

import (
	"github.com/klauspost/compress/zstd"
	"log"
)

// ZstdComp is a compressor which uses the zstandard
type ZstdComp struct {
}

// Compress compresses the data based on the zstandard
func (l *ZstdComp) Compress(data []byte) []byte {
	encoder, _ := zstd.NewWriter(nil)
	buffer := encoder.EncodeAll(data, make([]byte, 0, len(data)))
	encoder.Close()
	return buffer
}

// Decompress decompresses the data which was compressed by zstandard
func (l *ZstdComp) Decompress(data []byte) []byte {
	decoder, _ := zstd.NewReader(nil)
	buffer, err := decoder.DecodeAll(data, nil)
	if err != nil {
		log.Fatalf("Failed to decompress the buffer: %v", err)
	}
	decoder.Close()
	return buffer
}
