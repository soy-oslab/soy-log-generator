package compressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// GzipComp is a compressor which uses the gzip
type GzipComp struct {
}

// Compress compresses the data based on the gzip
func (l *GzipComp) Compress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write(data)
	writer.Close()
	return buffer.Bytes(), err
}

// Decompress decompresses the data which was compressed by gzip
func (l *GzipComp) Decompress(data []byte) ([]byte, error) {
	var ret []byte
	var err error
	var reader *gzip.Reader

	buffer := bytes.NewBuffer(data)
	reader, err = gzip.NewReader(buffer)
	if err != nil {
		goto out
	}
	ret, err = ioutil.ReadAll(reader)
	reader.Close()
out:
	return ret, err
}
