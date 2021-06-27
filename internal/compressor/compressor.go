package compressor

import (
	"bytes"
	"log"
	//"compress/gzip"
	"github.com/klauspost/compress/zstd"
	"io/ioutil"
)

func Compress(data []byte) []byte {
	var buffer bytes.Buffer
	// writer := zstd.NewWriter(&buffer)
	writer, _ := zstd.NewWriter(&buffer)
	writer.Write(data)
	writer.Close()
	return buffer.Bytes()
}

func Decompress(data []byte) []byte {
	buffer := bytes.NewBuffer(data)
	reader, err := zstd.NewReader(buffer)
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
