package compressor

import (
	"github.com/klauspost/compress/zstd"
	"log"
)

type ZstdComp struct {
}

func (l *ZstdComp) Compress(data []byte) []byte {
	encoder, _ := zstd.NewWriter(nil)
	buffer := encoder.EncodeAll(data, make([]byte, 0, len(data)))
	encoder.Close()
	return buffer
}

func (l *ZstdComp) Decompress(data []byte) []byte {
	decoder, _ := zstd.NewReader(nil)
	buffer, err := decoder.DecodeAll(data, nil)
	if err != nil {
		log.Fatalf("Failed to decompress the buffer: %v", err)
	}
	decoder.Close()
	return buffer
}
