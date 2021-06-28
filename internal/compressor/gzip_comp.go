package compressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
)

type GzipComp struct {
}

func (l *GzipComp) Compress(data []byte) []byte {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writer.Write(data)
	writer.Close()
	return buffer.Bytes()
}

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
