package compressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
)

// GzipComp is a compressor which uses the gzip
type GzipComp struct {
}

// Compress compresses the data based on the gzip
func (l *GzipComp) Compress(data []byte) []byte {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, _ = writer.Write(data)
	writer.Close()
	return buffer.Bytes()
}

// Decompress decompresses the data which was compressed by gzip
func (l *GzipComp) Decompress(data []byte) []byte {
	buffer := bytes.NewBuffer(data)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		log.Fatalf("Failed to decompress the buffer: %v", err)
	}
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to read the buffer: %v", err)
	}
	reader.Close()
	return result
}
