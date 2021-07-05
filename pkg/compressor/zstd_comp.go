package compressor

import (
	"github.com/klauspost/compress/zstd"
)

// ZstdComp is a compressor which uses the zstandard
type ZstdComp struct {
}

// Compress compresses the data based on the zstandard
func (l *ZstdComp) Compress(data []byte) ([]byte, error) {
	encoder, err := zstd.NewWriter(nil)
	buffer := encoder.EncodeAll(data, make([]byte, 0, len(data)))
	encoder.Close()
	return buffer, err
}

// Decompress decompresses the data which was compressed by zstandard
func (l *ZstdComp) Decompress(data []byte) ([]byte, error) {
	decoder, err := zstd.NewReader(nil)
	buffer, err := decoder.DecodeAll(data, nil)
	decoder.Close()
	return buffer, err
}
